package conceal

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

// Addr conceals addr string in "0.0.0.0:0" ip:port format
func (c *Conceal) Addr(addr string) (eaddr string, err error) {
	addrs := strings.Split(addr, ":")
	if len(addrs) != 2 {
		return "", ErrInvalidAddress
	}
	eport, err := c.Port(addrs[1])
	if err != nil {
		return "", err
	}
	eip, err := c.IP(addrs[0])
	if err != nil {
		return "", err
	}
	return eip + ":" + eport, nil
}

// IP conceals IPv4 address in standard string format
func (c *Conceal) IP(ip string) (eip string, err error) {
	nip := net.ParseIP(ip)
	if nip == nil {
		return "", ErrInvalidIP
	}
	ips := strings.Split(ip, ".")
	if len(ips) != 4 {
		return "", ErrInvalidIP
	}
	nipbytes := []byte(nip.To4())
	for i := range nipbytes {
		x := int(nipbytes[i])
		xi := x + i
		if xi > 255 {
			nipbytes[i] = byte(xi - 256)
			continue
		}
		nipbytes[i] = byte(xi)
	}
	nip = c.Bytes(nipbytes)
	return nip.String(), nil
}

// Port conceals standard port in numeric format
func (c *Conceal) Port(port string) (eport string, err error) {
	porti, err := strconv.Atoi(port)
	if err != nil || porti < 0 || porti > 65535 {
		return "", ErrInvalidPort
	}
	portb := []byte{0x00, 0x00}
	binary.LittleEndian.PutUint16(portb, uint16(porti))
	portb = c.Bytes(portb)
	porti = int(binary.LittleEndian.Uint16(portb))
	return strconv.Itoa(porti), nil
}
