package ipamd2

import (
	"context"
	"fmt"

	"github.com/aws/amazon-vpc-cni-k8s/rpc"
)

type rpcHandler struct {
	version   string
	innerIpam Ipam
}

func (s *rpcHandler) AddNetwork(ctx context.Context, request *rpc.AddNetworkRequest) (*rpc.AddNetworkReply, error) {
	fmt.Printf("AddNetwork Requested: %v\n", request)

	ipv4, err := s.innerIpam.LeaseIp(&IpamKey{
		ContainerId: request.ContainerID,
		IfName:      request.IfName,
	})
	if err != nil {
		fmt.Printf("LeaseIp Failed, err = %s\n", err.Error())
		return &rpc.AddNetworkReply{
			Success: false,
		}, nil
	}

	return &rpc.AddNetworkReply{
		Success:         true,
		IPv4Addr:        ipv4.String(),
		DeviceNumber:    0,
		UseExternalSNAT: true,
		VPCv4CIDRs:      []string{"100.66.0.0/19"},
		VPCv6CIDRs:      nil,
	}, nil
}

func (s *rpcHandler) DelNetwork(ctx context.Context, request *rpc.DelNetworkRequest) (*rpc.DelNetworkReply, error) {
	fmt.Printf("DelNetwork Requested: %v\n", request)

	ipv4, err := s.innerIpam.LeaseIp(&IpamKey{
		ContainerId: request.ContainerID,
		IfName:      request.IfName,
	})
	if err != nil {
		fmt.Printf("LeaseIp Failed, err = %s\n", err.Error())
		return &rpc.DelNetworkReply{
			Success: false,
		}, nil
	}

	return &rpc.DelNetworkReply{
		Success:      true,
		IPv4Addr:     ipv4.String(),
		DeviceNumber: 0,
	}, nil
}
