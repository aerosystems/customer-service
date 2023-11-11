package RPCServer

type CustomerRPCPayload struct {
	UserId int
}

func (us *CustomerServer) CreateCustomer(arg, payload *CustomerRPCPayload) error {
	user, err := us.userService.CreateUser()
	if err != nil {
		return err
	}
	payload.UserId = user.Id
	return nil
}
