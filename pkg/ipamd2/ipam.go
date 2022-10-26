package ipamd2

import (
	"encoding/binary"
	"fmt"
	"net"
)

type IpamKey struct {
	ContainerId string `json:"containerId"`
	IfName      string `json:"ifName"`
}

type Ipam interface {
	LeaseIp(key *IpamKey) (*net.IP, error)
	ReleaseIp(key *IpamKey) error
}

// Only support IPv4
type inMemoryIpam struct {
	rangeStart uint32
	rangeEnd   uint32
	cursor     uint32
	data       map[IpamKey]*net.IP
}

func newInMemoryIpam(rangeStart net.IP, rangeEnd net.IP) *inMemoryIpam {
	startIp := ip2int(rangeStart)
	endIp := ip2int(rangeEnd)

	return &inMemoryIpam{
		rangeStart: startIp,
		rangeEnd:   endIp,
		cursor:     startIp,
		data:       map[IpamKey]*net.IP{},
	}
}

func (i *inMemoryIpam) LeaseIp(key *IpamKey) (*net.IP, error) {
	fmt.Printf("LeaseIp: Stats pre-start: %d, end: %d, cursor: %d", i.rangeStart, i.rangeEnd, i.cursor)
	if val, found := i.data[*key]; found {
		return val, nil
	}

	i.cursor++
	if i.cursor > i.rangeEnd {
		return nil, fmt.Errorf("IP exhausted, No available ips in ip-range. start: %d, end: %d", i.rangeStart, i.rangeEnd)
	}

	ip := int2ip(i.cursor)
	i.data[*key] = &ip

	fmt.Printf("LeaseIp: Stats lease one IP")

	return &ip, nil
}

func (i *inMemoryIpam) ReleaseIp(key *IpamKey) error {
	if _, found := i.data[*key]; found {
		delete(i.data, *key)
		fmt.Printf("ReleaseIp: Stats: %d, end: %d, cursor: %d", i.rangeStart, i.rangeEnd, i.cursor)
		return nil
	} else {
		return fmt.Errorf("ReleaseIp: Failed to find IPAM entry, key: %v", key)
	}
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		panic("no sane way to convert ipv6 into uint32")
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
