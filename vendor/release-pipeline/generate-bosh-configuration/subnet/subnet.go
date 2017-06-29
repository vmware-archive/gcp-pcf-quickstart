package subnet

import (
	"fmt"
	"net"
)

type Subnet struct {
	net.IP
	*net.IPNet
}

func ParseSubnet(s string) (Subnet, error) {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return Subnet{}, err
	}

	return Subnet{ip, ipNet}, nil
}

func (s Subnet) IPAddress(octet int) (string, error) {
	if octet < 0 {
		return "", fmt.Errorf("IP address containing \"%d\" is invalid, cannot be negative", octet)
	}

	if octet > 255 {
		return "", fmt.Errorf("IP address containing \"%d\" is invalid, cannot exceed 255", octet)
	}

	ip := s.IP.To4()
	ip = ip.Mask(s.IPNet.Mask)

	for i := 0; i < octet; i++ {
		ip[3]++
	}

	return ip.String(), nil
}

func (s Subnet) Range(start, end int) (string, error) {
	if start < 0 {
		return "", fmt.Errorf("subnet range start \"%d\" cannot be negative", start)
	}

	if end > 255 {
		return "", fmt.Errorf("subnet range end \"%d\" cannot exceed 255", end)
	}

	if start > end {
		return "", fmt.Errorf("subnet range start \"%d\" cannot exceed subnet range end \"%d\"", start, end)
	}

	startIP := s.IP.To4()
	for i := 0; i < start; i++ {
		startIP[3]++
	}

	startRange := startIP.String()

	endIP := s.IP.To4()
	for i := start; i < end; i++ {
		endIP[3]++
	}

	return fmt.Sprintf("%s-%s", startRange, endIP), nil
}
