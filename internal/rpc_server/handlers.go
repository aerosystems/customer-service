package RPCServer

type UserRPCPayload struct {
	UserId       uint
	IsActive     bool
	Role         string
	Email        string
	PasswordHash string
}

func (us *CustomerServer) CreateCustomer(payload UserRPCPayload, userId *uint) error {
	user, err := us.userService.CreateUser(payload.Email, payload.PasswordHash)
	if err != nil {
		return err
	}
	*userId = user.Id
	return nil
}
