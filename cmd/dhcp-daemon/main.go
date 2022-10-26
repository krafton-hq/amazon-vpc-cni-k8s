package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/ipamd2"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/version"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var mode string
	var stateFile string
	flag.StringVar(&mode, "mode", "", "")
	flag.StringVar(&stateFile, "file", "/var/lib/run/dhcp-cni/lease.json", "")

	ipam := ipamd2.NewIpamServer("eth0")
	fmt.Println("Start RPC Handler")
	// Start the RPC listener
	err := ipam.RunRPCHandler(version.Version)
	if err != nil {
		fmt.Printf("Failed to set up gRPC handler: %v\n", err)
		return 1
	}
	return 0
}
