package addr

import (
	"net"
	"strconv"
	"strings"
)

type Addr interface {
	net.Addr
	Host(args ...string) string
	Port(args ...string) string
	MyIP() string
	PortInt() int
	SetNetwork(string)
}

type addr struct {
	host string
	port string
	name string
}

var _ Addr = &addr{}

func NewAddr(args ...string) Addr {
	ad := &addr{}
	switch len(args) {
	case 1:
		ad.Host(args[0])
		if i := strings.IndexByte(args[0], ':'); i >= 0 {
			return NewAddr(args[0][:i], args[0][i+1:])
		}
	case 2:
		ad.Host(args[0])
		ad.Port(args[1])
	case 3:
		ad.Host(args[0])
		ad.Port(args[1])
		ad.SetNetwork(args[2])
	}
	return ad
}

func (ad *addr) Network() string {
	return ad.name
}

func (ad *addr) SetNetwork(name string) {
	ad.name = name
}

func (ad *addr) String() string {
	return net.JoinHostPort(ad.Host(), ad.Port())
}

func (ad *addr) Host(args ...string) string {

	if len(args) >= 1 {
		ad.host = strings.Replace(args[0], " ", "", -1)
	}
	return ad.host
}

func (ad *addr) Port(args ...string) string {
	if len(args) >= 1 {
		ad.port = strings.Replace(args[0], " ", "", -1)
	}
	return ad.port
}

func (ad *addr) PortInt() int {
	p, ok := strconv.Atoi(ad.port)
	if ok == nil {
		return p
	}
	return 0
}

// TODO:
func (ad *addr) MyIP() string {
	return ""
}
