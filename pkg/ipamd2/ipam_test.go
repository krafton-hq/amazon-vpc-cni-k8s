package ipamd2

import (
	"fmt"
	"net"
	"testing"
)

func Test_inMemoryIpam_LeaseIp(t *testing.T) {
	ipam := newInMemoryIpam(net.IPv4(192, 168, 0, 200), net.IPv4(192, 168, 1, 100))
	for i := 0; i < 155; i++ {
		_, err := ipam.LeaseIp(&IpamKey{
			ContainerId: fmt.Sprintf("cont-%d", i),
			IfName:      "eth0",
		})
		if err != nil {
			t.Errorf("LeaseIp() error = %v\n", err)
			return
		}
	}
	r1, r2, l := ipam.Status()
	fmt.Printf("Status() range = %d, remain = %d\n, leases = %d\n", r1, r2, l)
}

func Test_inMemoryIpam_Re_LeaseIp(t *testing.T) {
	ipam := newInMemoryIpam(net.IPv4(192, 168, 0, 200), net.IPv4(192, 168, 1, 100))
	key1 := &IpamKey{ContainerId: "cont-1", IfName: "eth0"}
	key2 := &IpamKey{ContainerId: "cont-2", IfName: "eth0"}
	key3 := &IpamKey{ContainerId: "cont-3", IfName: "eth0"}

	if _, err := ipam.LeaseIp(key1); err != nil {
		t.Errorf("LeaseIp() error = %v\n", err)
		return
	}
	if _, err := ipam.LeaseIp(key2); err != nil {
		t.Errorf("LeaseIp() error = %v\n", err)
		return
	}
	if _, err := ipam.LeaseIp(key3); err != nil {
		t.Errorf("LeaseIp() error = %v\n", err)
		return
	}

	r1, r2, l := ipam.Status()
	fmt.Printf("Before ReLease Status() range = %d, remain = %d, leases = %d\n", r1, r2, l)

	if _, err := ipam.LeaseIp(key1); err != nil {
		t.Errorf("LeaseIp() error = %v\n", err)
		return
	}
	if _, err := ipam.LeaseIp(key2); err != nil {
		t.Errorf("LeaseIp() error = %v\n", err)
		return
	}
	if _, err := ipam.LeaseIp(key3); err != nil {
		t.Errorf("LeaseIp() error = %v\n", err)
		return
	}

	r1, r2, l = ipam.Status()
	fmt.Printf("After ReLease Status() range = %d, remain = %d, leases = %d\n", r1, r2, l)
}

func Test_inMemoryIpam_ReleaseIp(t *testing.T) {
	ipam := newInMemoryIpam(net.IPv4(192, 168, 0, 200), net.IPv4(192, 168, 1, 100))
	key1 := &IpamKey{ContainerId: "cont-1", IfName: "eth0"}
	if _, err := ipam.LeaseIp(key1); err != nil {
		t.Errorf("LeaseIp() error = %v\n", err)
		return
	}

	r1, r2, l := ipam.Status()
	fmt.Printf("Before Release Status() range = %d, remain = %d, leases = %d\n", r1, r2, l)

	if err := ipam.ReleaseIp(key1); err != nil {
		t.Errorf("ReleaseIp() error = %v\n", err)
		return
	}

	r1, r2, l = ipam.Status()
	fmt.Printf("After Release Status() range = %d, remain = %d, leases = %d\n", r1, r2, l)
}
