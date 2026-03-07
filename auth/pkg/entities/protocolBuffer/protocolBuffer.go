package protocolbuffer

import authPB "github.com/pnaskardev/URL-Shortner-V1/url-shortner-rpc/auth"

type ProtoAuthServer struct {
	authPB.UnimplementedAuthServer
}
