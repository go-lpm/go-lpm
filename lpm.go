package main

import (
	"net"
)

const (
	ArchRadix = "radix"
)

type LPMTable interface {
	Show() [][]Entry
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
func NewLPMTable(arch string) LPMTable {
	switch arch {
	case ArchRadix:
		return &RadixTable{}
	default:
		return &RadixTable{}
	}
}

// NewRadixLPMTable Create a lpm table based on radix arch.
func NewRadixLPMTable() LPMTable {
	return NewLPMTable(ArchRadix)
}
