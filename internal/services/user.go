package services

import (
	"errors"
	"fmt"
	"github.com/aerosystems/user-service/internal/RPCServices"
	"github.com/aerosystems/user-service/internal/helpers"
	"github.com/aerosystems/user-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

type UserService interface {
	GetUserById(userId uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(email, passwordHash string) (*models.User, error)
	ResetPassword(userId uint, passwordHash string) error
	ActivateUser(userId uint) error
	MatchPassword(email, passwordHash string) (bool, error)
}

type UserServiceImpl struct {
	userRepo        models.UserRepository
	codeRepo        models.CodeRepository
	checkmailRPC    *RPCServices.CheckmailRPC
	mailRPC         *RPCServices.MailRPC
	projectRPC      *RPCServices.ProjectRPC
	subscriptionRPC *RPCServices.SubscriptionRPC
}

func NewUserServiceImpl(userRepo models.UserRepository, codeRepo models.CodeRepository, checkmailRPC *RPCServices.CheckmailRPC, mailRPC *RPCServices.MailRPC, projectRPC *RPCServices.ProjectRPC, subscriptionRPC *RPCServices.SubscriptionRPC) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:        userRepo,
		codeRepo:        codeRepo,
		checkmailRPC:    checkmailRPC,
		mailRPC:         mailRPC,
		projectRPC:      projectRPC,
		subscriptionRPC: subscriptionRPC,
	}
}

func (us *UserServiceImpl) Register(email, password, clientIp string) error {
	// hashing password
	passwordHash, _ := us.hashPassword(password)
	// checking email in blacklist via RPC
	if isTrust, err := us.checkmailRPC.IsTrustEmail(email, clientIp); err != nil {
		if !isTrust {
			return errors.New("email address contains in blacklist")
		}
	}
	// getting user by email via RPC
	user, _ := us.userRepo.FindByEmail(email)
	// if user with this email already exists
	if user != nil {
		if user.IsActive {
			return errors.New("user with this email already exists")
		} else {
			// updating password for inactive user
			user.PasswordHash = passwordHash
			if err := us.userRepo.Update(user); err != nil {
				return errors.New("could not update password for inactive user")
			}
			code, _ := us.codeRepo.GetLastIsActiveCode(user.Id, "registration")
			if code == nil {
				// generating confirmation code
				codeObj := NewCode(user.Id, "registration", "")
				if err := us.codeRepo.Create(codeObj); err != nil {
					return errors.New("could not create new code")
				}
			} else {
				// extend expiration code and return previous active code
				if err := us.codeRepo.ExtendExpiration(code); err != nil {
					return fmt.Errorf("could not extend expiration code: %s", err.Error())
				}
			}
			// sending confirmation code via RPC
			if err := us.mailRPC.SendEmail(email, "Confirm your email🗯", fmt.Sprintf("Your confirmation code is %s", code.Code)); err != nil {
				return fmt.Errorf("could not send email: %s", err.Error())
			}
			return nil
		}
	}
	// creating new user via RPC
	userId, err := us.userRPC.CreateUser(email, passwordHash)
	if err != nil {
		return fmt.Errorf("could not create new user: %s", err.Error())
	}
	// generating confirmation code
	code, err := us.codeRepo.NewCode(userId, "registration", "")
	if err != nil {
		return fmt.Errorf("could not gen new code: %s", err.Error())
	}
	// sending confirmation code via RPC
	if err := us.mailRPC.SendEmail(email, "Confirm your email🗯", fmt.Sprintf("Your confirmation code is %s", code.Code)); err != nil {
		return fmt.Errorf("could not send email: %s", err.Error())
	}
	return nil
}

func (us *UserServiceImpl) Confirm(code *models.Code) error {
	switch code.Action {
	case "registration":
		if err := us.ActivateUser(code.UserId); err != nil {
			return fmt.Errorf("could not activate user: %s", err.Error())
		}
		code.IsUsed = true
		if err := us.codeRepo.Update(code); err != nil {
			return errors.New("could not confirm registration")
		}
		// create default project via RPC
		if err := us.projectRPC.CreateDefaultProject(code.UserId); err != nil {
			return fmt.Errorf("could not create default project: %s", err.Error())
		}
		// create default subscription via RPC
		if err := us.subscriptionRPC.CreateFreeTrial(code.UserId); err != nil {
			return fmt.Errorf("could not create default subscription: %s", err.Error())
		}
	case "reset_password":
		if err := us.ActivateUser(code.UserId); err != nil {
			return fmt.Errorf("could not activate user: %s", err.Error())
		}
		code.User.PasswordHash = code.Data
		if err := us.userRepo.Update(code.User); err != nil {
			return errors.New("could not update password")
		}
		code.IsUsed = true
		if err := us.codeRepo.Update(code); err != nil {
			return fmt.Errorf("could not confirm reset password: %s", err.Error())
		}
	}
	return nil
}

func (us *UserServiceImpl) GetUserById(userId uint) (*models.User, error) {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this id does not exist")
	}
	return user, nil
}

func (us *UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("could not get user id")
	}
	if user == nil {
		return nil, errors.New("user with this email does not exist")
	}
	return user, nil
}

func (us *UserServiceImpl) CreateUser(email, passwordHash string) (*models.User, error) {
	user := models.NewUserEntity(email, passwordHash, "user")
	if err := us.userRepo.Create(user); err != nil {
		return nil, errors.New("could not create new user")
	}
	return user, nil
}

func (us *UserServiceImpl) ResetPasswordDeprecated(userId uint, passwordHash string) error {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return errors.New("could not get user id")
	}
	if user == nil {
		return errors.New("user with this email does not exist")
	}
	user.PasswordHash = passwordHash
	if err := us.userRepo.Update(user); err != nil {
		return errors.New("could not update password")
	}
	return nil
}

func (us *UserServiceImpl) ResetPassword(email, password string) error {
	// hashing password
	passwordHash, _ := us.hashPassword(password)
	// get user by email via RPC
	user, err := us.userRPC.GetUserByEmail(email)
	if err != nil {
		return errors.New("could not get user")
	}
	code, err := us.codeRepo.GetLastIsActiveCode(user.UserId, "reset_password")
	if err != nil {
		return errors.New("could not get last active code")
	}
	if code == nil || code.IsUsed {
		_, err := us.codeRepo.NewCode(user.UserId, "reset_password", passwordHash)
		if err != nil {
			return errors.New("could not gen new code")
		}
	}
	// extend expiration code and return previous active code
	code.Data = passwordHash
	if err := us.codeRepo.ExtendExpiration(code); err != nil {
		return errors.New("could not extend expiration code")
	}
	// sending confirmation code via RPC
	if err := us.mailRPC.SendEmail(email, "Reset your password🗯", fmt.Sprintf("Your confirmation code is %s", code.Code)); err != nil {
		return errors.New("could not send email")
	}
	return nil
}

func (us *UserServiceImpl) ActivateUser(userId uint) error {
	user, err := us.userRepo.FindById(userId)
	if err != nil {
		return errors.New("could not get user id")
	}
	if user == nil {
		return errors.New("user with this email does not exist")
	}
	user.IsActive = true
	if err := us.userRepo.Update(user); err != nil {
		return errors.New("could not activate user")
	}
	return nil
}

func (us *UserServiceImpl) MatchPasswordDeprecated(email, passwordHash string) (bool, error) {
	user, err := us.userRepo.FindByEmail(email)
	if err != nil {
		return false, errors.New("could not get user id")
	}
	if user == nil {
		return false, errors.New("user with this email does not exist")
	}
	if user.PasswordHash != passwordHash {
		return false, errors.New("password is incorrect")
	}
	return true, nil
}

func (us *UserServiceImpl) MatchPassword(email, password string) (*RPCServices.UserRPCPayload, error) {
	// get user by email via RPC
	user, err := us.userRPC.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("could not get user")
	}
	if user.IsActive == false {
		return nil, errors.New("user is not active")
	}
	// match password via RPC
	if err := us.userRPC.MatchPassword(email, password); err != nil {
		return nil, errors.New("password does not match")
	}
	return user, nil
}

func (us *UserServiceImpl) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", errors.New("could not hash password")
	}
	return string(hash), nil
}

func normalizeEmail(data string) string {
	addr := strings.ToLower(data)

	arrAddr := strings.Split(addr, "@")
	username := arrAddr[0]
	domain := arrAddr[1]

	googleDomains := strings.Split(os.Getenv("GOOGLEMAIL_DOMAINS"), ",")

	//checking Google mail aliases
	if helpers.Contains(googleDomains, domain) {
		//removing all dots from username mail
		username = strings.ReplaceAll(username, ".", "")
		//removing all characters after +
		if strings.Contains(username, "+") {
			res := strings.Split(username, "+")
			username = res[0]
		}
		addr = username + "@gmail.com"
	}

	return addr
}
