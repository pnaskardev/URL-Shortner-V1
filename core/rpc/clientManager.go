package rpc_service

type RPCClients struct {
	Auth *AuthRPCClient
	// TODO ADD MORE CLIENTS IN HERE
}

var rpcClientInstance *RPCClients

func LoadNewClients() error {

	// TODO add Auth Microservice URL
	authClient, err := NewAuthRPCClientConnection("")
	if err != nil {
		return err
	}

	rpcClientInstance = &RPCClients{
		Auth: authClient,
	}

	return nil
}

func GetRPCClient() *RPCClients {
	if rpcClientInstance == nil {
		err := LoadNewClients()
		if err != nil {
			panic(err)
		}
	}
	return rpcClientInstance
}
