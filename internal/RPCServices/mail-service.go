package RPCServices

import "net/rpc"

type MailService interface {
	SendEmail(to, subject, body string) error
}

type MailRPC struct {
	rpcClient *rpc.Client
}

func NewMailRPC(rpcClient *rpc.Client) *MailRPC {
	return &MailRPC{
		rpcClient: rpcClient,
	}
}

type MailRPCPayload struct {
	To      string
	Subject string
	Body    string
}

func (ms *MailRPC) SendEmail(to, subject, body string) error {
	var result string
	if err := ms.rpcClient.Call(
		"MailServer.SendEmail",
		MailRPCPayload{
			To:      to,
			Subject: subject,
			Body:    body,
		},
		&result); err != nil {
		return err
	}

	return nil
}
