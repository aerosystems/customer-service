package RPCServer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
)

const rpcPort = 5001

type Server struct {
	log             *logrus.Logger
	customerUsecase CustomerUsecase
}

func NewServer(
	log *logrus.Logger,
	customerUsecase CustomerUsecase,
) *Server {
	return &Server{
		log:             log,
		customerUsecase: customerUsecase,
	}
}

func (s Server) Run() error {
	if err := rpc.Register(s); err != nil {
		return err
	}
	s.log.Infof("starting customer-service RPC server on port %d\n", rpcPort)
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
