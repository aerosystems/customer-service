package RpcServer

import "github.com/google/uuid"

type CustomerRPCPayload struct {
	Uuid uuid.UUID
}

func (s Server) CreateCustomer(_ string, payload *CustomerRPCPayload) error {
	user, err := s.customerUsecase.CreateCustomer()
	if err != nil {
		return err
	}
	payload.Uuid = user.Uuid
	return nil
}
