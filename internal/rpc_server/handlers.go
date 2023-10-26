package RPCServer

type UserRPCPayload struct {
	UserId       uint
	IsActive     bool
	Role         string
	Email        string
	PasswordHash string
}

func (us *UserServer) CreateUser(payload UserRPCPayload, userId *uint) error {
	user, err := us.userService.CreateUser(payload.Email, payload.PasswordHash)
	if err != nil {
		return err
	}
	*userId = user.Id
	return nil
}

func (us *UserServer) GetUserByEmail(email string, payload *UserRPCPayload) error {
	user, err := us.userService.GetUserByEmail(email)
	if err != nil {
		return err
	}
	payload.UserId = user.Id
	payload.Email = user.Email
	payload.IsActive = user.IsActive
	payload.Role = user.Role
	payload.PasswordHash = user.PasswordHash
	return nil
}

func (us *UserServer) ResetPassword(payload UserRPCPayload, result *string) error {
	user, err := us.userService.GetUserById(payload.UserId)
	if err != nil {
		return err
	}
	return us.userService.ResetPassword(user.Id, payload.PasswordHash)
}

func (us *UserServer) ActivateUser(userId uint, result *string) error {
	return us.userService.ActivateUser(userId)
}

func (us *UserServer) MatchPassword(payload UserRPCPayload, result *bool) error {
	if res, err := us.userService.MatchPassword(payload.Email, payload.PasswordHash); err != nil {
		return err
	} else {
		*result = res
	}
	return nil
}
