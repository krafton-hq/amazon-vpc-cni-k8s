package ipamd2

import (
	"context"
	"fmt"

	"github.com/aws/amazon-vpc-cni-k8s/rpc"
)

type rpcHandler struct {
	version string
}

func (s *rpcHandler) AddNetwork(ctx context.Context, request *rpc.AddNetworkRequest) (*rpc.AddNetworkReply, error) {
	fmt.Printf("AddNetwork Requested: %v\n", request)

	//lastIp := rand.IntnRange(200, 240)
	return &rpc.AddNetworkReply{
		Success:         true,
		IPv4Addr:        fmt.Sprintf("100.66.5.241"),
		DeviceNumber:    0,
		UseExternalSNAT: true,
		VPCv4CIDRs:      []string{"100.66.0.0/19"},
		VPCv6CIDRs:      nil,
	}, nil
}

func (s *rpcHandler) DelNetwork(ctx context.Context, request *rpc.DelNetworkRequest) (*rpc.DelNetworkReply, error) {
	fmt.Printf("DelNetwork Requested: %v\n", request)

	//lastIp := rand.IntnRange(200, 240)
	return &rpc.DelNetworkReply{
		Success:      true,
		IPv4Addr:     fmt.Sprintf("100.66.5.241/19"),
		DeviceNumber: 0,
	}, nil
}
