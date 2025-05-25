package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeaeder = "grpcgateway-user-agent"
	userAgentHeader             = "user-agent"
	xForwarderForHeader         = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClintIP   string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	meta := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeaeder); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}

		if clientIps := md.Get(xForwarderForHeader); len(clientIps) > 0 {
			meta.ClintIP = clientIps[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		meta.ClintIP = p.Addr.String()
	}

	return meta
}
