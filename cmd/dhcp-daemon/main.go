package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/ipamd2"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/version"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var (
		mode      string
		stateFile string

		rangeStart string
		rangeEnd   string
	)

	flag.StringVar(&mode, "mode", "", "")
	flag.StringVar(&stateFile, "file", "/var/lib/run/dhcp-cni/lease.json", "")
	flag.StringVar(&rangeStart, "start", "192.168.0.200", "IPAM Range Start")
	flag.StringVar(&rangeEnd, "end", "192.168.0.220", "IPAM Range End")
	flag.Parse()

	ipam := ipamd2.NewIpamServer("eth0")
	fmt.Println("Start RPC Handler")
	// Start the RPC listener
	err := ipam.RunRPCHandler(version.Version, net.ParseIP(rangeStart), net.ParseIP(rangeEnd))
	if err != nil {
		fmt.Printf("Failed to set up gRPC handler: %v\n", err)
		return 1
	}
	return 0
}
