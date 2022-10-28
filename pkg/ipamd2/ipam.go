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
	ReleaseIp(key *IpamKey) (*net.IP, error)
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

	fmt.Printf("New Inmemory IPAM Initialize: startIp(int32) = %d, endIp(int32) = %d, startIp = %s, endIp = %s\n", startIp, endIp, rangeStart.String(), rangeEnd.String())

	return &inMemoryIpam{
		rangeStart: startIp,
		rangeEnd:   endIp,
		cursor:     startIp,
		data:       map[IpamKey]*net.IP{},
	}
}

func (i *inMemoryIpam) LeaseIp(key *IpamKey) (*net.IP, error) {
	r1, r2, l := i.Status()
	fmt.Printf("LeaseIp: pre-start, Status: range = %d, remain = %d, leases = %d\n", r1, r2, l)
	if val, found := i.data[*key]; found {
		return val, nil
	}

	i.cursor++
	if i.cursor > i.rangeEnd {
		return nil, fmt.Errorf("IP exhausted, No available ips in ip-range. start: %d, end: %d", i.rangeStart, i.rangeEnd)
	}

	ip := int2ip(i.cursor)
	i.data[*key] = &ip

	r1, r2, l = i.Status()
	fmt.Printf("LeaseIp: finish, Status: range = %d, remain = %d, leases = %d\n", r1, r2, l)

	return &ip, nil
}

func (i *inMemoryIpam) ReleaseIp(key *IpamKey) (*net.IP, error) {
	if val, found := i.data[*key]; found {
		delete(i.data, *key)
		r1, r2, l := i.Status()
		fmt.Printf("ReleaseIp: Status: range = %d, remain = %d, leases = %d\n", r1, r2, l)
		return val, nil
	} else {
		return nil, fmt.Errorf("ReleaseIp: Failed to find IPAM entry, key: %v", key)
	}
}

func (i *inMemoryIpam) Status() (uint32, uint32, int) {
	range2 := i.rangeEnd - i.rangeStart
	remain := i.rangeEnd - i.cursor
	leases := len(i.data)
	return range2, remain, leases
}

func ip2int(ip net.IP) uint32 {
	if ip.To4() == nil {
		panic("no sane way to convert ipv6 into uint32")
	}
	return binary.BigEndian.Uint32(ip.To4())
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
