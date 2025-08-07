package golpm

import (
	"net"
)

const (
	ArchRadix = "radix"
)

type LPMTable interface {
	Show() map[int][]Entry
	Add(prefix string, entry interface{}) error
	AddIPNet(prefix *net.IPNet, entry interface{}) error
	Delete(prefix string) error
	DeleteIPNet(prefix *net.IPNet) error
	Lookup(ip string) *Entry
	LookupIP(ip net.IP) *Entry
}

// Entry An entry in entries table
type Entry struct {
	Prefix *net.IPNet
	Entry  interface{}
}

// NewLPMTable Create a lpm table based on specify arch.
func NewLPMTable(arch string, isIPv6 bool) LPMTable {
	ipBytesLen := net.IPv4len
	if isIPv6 {
		ipBytesLen = net.IPv6len
	}
	switch arch {
	case ArchRadix:
		return &RadixTable{
			ipBytesLen: ipBytesLen,
			root:       &radixNode{},
		}
	default:
		return &RadixTable{
			ipBytesLen: ipBytesLen,
			root:       &radixNode{},
		}
	}
}

// NewRadixLPMTable Create a lpm table based on radix arch.
func NewRadixLPMTable(isIPv6 bool) LPMTable {
	return NewLPMTable(ArchRadix, isIPv6)
}
