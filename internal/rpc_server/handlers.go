package RPCServer

import (
	"github.com/google/uuid"
)

type CustomerRPCPayload struct {
	Uuid uuid.UUID
}

func (us *CustomerServer) CreateCustomer(arg, payload *CustomerRPCPayload) error {
	user, err := us.userService.CreateUser()
	if err != nil {
		return err
	}
	payload.Uuid = user.Uuid
	return nil
}
