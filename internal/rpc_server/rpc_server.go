package RPCServer

import (
	"fmt"
	"github.com/aerosystems/customer-service/internal/services"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

type CustomerServer struct {
	rpcPort     int
	log         *logrus.Logger
	userService services.CustomerService
}

func NewUserServer(
	rpcPort int,
	log *logrus.Logger,
	customerService services.CustomerService,
) *CustomerServer {
	return &CustomerServer{
		rpcPort:     rpcPort,
		log:         log,
		userService: customerService,
	}
}

func (us *CustomerServer) Listen(rpcPort int) error {
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
