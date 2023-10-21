package RPCServer

import (
	"fmt"
	"github.com/aerosystems/user-service/internal/services"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

type UserServer struct {
	rpcPort     int
	log         *logrus.Logger
	userService services.UserService
}

func NewUserServer(
	rpcPort int,
	log *logrus.Logger,
	userService services.UserService,
) *UserServer {
	return &UserServer{
		rpcPort:     rpcPort,
		log:         log,
		userService: userService,
	}
}

func (ss *UserServer) Listen(rpcPort int) error {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
