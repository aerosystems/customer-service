package RPCServer

import "github.com/google/uuid"

type CustomerRPCPayload struct {
	Uuid uuid.UUID
}

func (us *CustomerServer) CreateCustomer(_ string, payload *CustomerRPCPayload) error {
	user, err := us.customerService.CreateUser()
	if err != nil {
		return err
	}
	payload.Uuid = user.Uuid
	return nil
}
