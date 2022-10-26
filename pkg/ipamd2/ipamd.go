package ipamd2

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/amazon-vpc-cni-k8s/rpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	ipamdgRPCaddress      = "127.0.0.1:50051"
	grpcHealthServiceName = "grpc.health.v1.aws-node"
)

type IpamServer struct {
	masterLink string
	grpcServer *grpc.Server
}

func NewIpamServer(masterLink string) *IpamServer {
	return &IpamServer{masterLink: masterLink}
}

func (ipam *IpamServer) RunRPCHandler(version string) error {
	listener, err := net.Listen("tcp", ipamdgRPCaddress)
	if err != nil {
		fmt.Printf("Failed to listen gRPC port: %v\n", err)
		return errors.Wrap(err, "ipamd: failed to listen to gRPC port")
	}

	grpcServer := grpc.NewServer()
	rpc.RegisterCNIBackendServer(grpcServer, &rpcHandler{version: version})
	healthServer := health.NewServer()
	// If ipamd can talk to the API server and to the EC2 API, the pod is healthy.
	// No need to ever change this to HealthCheckResponse_NOT_SERVING since it's a local service only
	healthServer.SetServingStatus(grpcHealthServiceName, healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	ipam.grpcServer = grpcServer

	go ipam.shutdownListener()
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("Failed to start server on gRPC port: %v\n", err)
		return errors.Wrap(err, "ipamd: failed to start server on gPRC port")
	}
	return nil
}

// shutdownListener - Listen to signals and set ipamd to be in status "terminating"
func (ipam *IpamServer) shutdownListener() {
	fmt.Println("Setting up shutdown hook.")
	sig := make(chan os.Signal, 1)

	// Interrupt signal sent from terminal
	signal.Notify(sig, syscall.SIGINT)
	// Terminate signal sent from Kubernetes
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	fmt.Println("Received shutdown signal, shutdown grpc server")
	ipam.grpcServer.GracefulStop()
	fmt.Println("Shutdown finished")
}
