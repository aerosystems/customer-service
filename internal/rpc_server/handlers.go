package RPCServer

type CreateUserRPCPayload struct {
	Email    string
	Password string
}

type UserRPCPayload struct {
	UserId   uint
	IsActive bool
}

func (ss *UserServer) GetUserByEmail(email string) (user UserRPCPayload, err error) {
	return
}
