package golpm

import (
	"math"
	"net"
)

type RadixTable struct {
	root         *radixNode
	defaultEntry *Entry
}

type radixNode struct {
	children [256]*radixNode
	entries  [8]*Entry // routing entries stored by the node
}

func traverse(deep int, node *radixNode, entries [][]Entry) {
	if node == nil {
		return
	}
	for maskTail, entry := range node.entries {
		if entry == nil {
			continue
		}
		maskSize := deep*8 + maskTail + 1
		entries[maskSize] = append(entries[maskSize], *entry)
	}
	for _, child := range node.children {
		traverse(deep+1, child, entries)
	}
}

func (rt *RadixTable) Show() [][]Entry {
	var entries [][]Entry
	entries = make([][]Entry, 33)
	for _, node := range rt.root.children {
		traverse(0, node, entries)
	}
	if rt.defaultEntry != nil {
		entries[0] = append(entries[0], *rt.defaultEntry)
	}
	return entries
}

func (rt *RadixTable) Add(prefix string, entry interface{}) error {
	_, ipNet, err := net.ParseCIDR(prefix)
	if err != nil {
		return err
	}
	return rt.AddIPNet(ipNet, entry)
}

func (rt *RadixTable) AddIPNet(prefix *net.IPNet, entry interface{}) error {
	maskSize, _ := prefix.Mask.Size()
	if maskSize == 0 {
		rt.defaultEntry = &Entry{
			Prefix: prefix,
			Entry:  entry,
		}
		return nil
	}

	var curNode *radixNode
	var curByte byte

	byteCount := (maskSize + 7) / 8
	ipBytes := []byte(prefix.IP)

	curNode = rt.root

	// process add byte-by-byte
	for i := 0; i < byteCount; i++ {
		curByte = ipBytes[i]
		if curNode.children[curByte] == nil {
			curNode.children[curByte] = &radixNode{}
		}
		curNode = curNode.children[curByte]
	}
	// save entry in end point
	curNode.entries[(maskSize+7)%8] = &Entry{
		Prefix: prefix,
		Entry:  entry,
	}
	return nil
}

func (rt *RadixTable) Delete(prefix string) error {
	_, ipNet, err := net.ParseCIDR(prefix)
	if err != nil {
		return err
	}
	return rt.DeleteIPNet(ipNet)
}

func (rt *RadixTable) DeleteIPNet(prefix *net.IPNet) error {
	maskSize, _ := prefix.Mask.Size()
	if maskSize == 0 {
		rt.defaultEntry = nil
		return nil
	}

	var curNode *radixNode
	var curByte byte

	byteCount := (maskSize + 7) / 8
	ipBytes := []byte(prefix.IP)

	curNode = rt.root

	// find corresponding node byte-by-byte
	for i := 0; i < byteCount; i++ {
		curByte = ipBytes[i]
		if curNode.children[curByte] == nil {
			return nil
		}
		curNode = curNode.children[curByte]
	}
	// delete entry from end point
	curNode.entries[(maskSize+7)%8] = nil
	return nil
}

func (rt *radixNode) lookupOneNode(val byte) *Entry {
	preVal := val + 1
	mask := byte(math.MaxUint8)

	for i := 0; i < 8; i++ {
		val &= mask
		if val == preVal || rt.children[val] == nil {
			preVal = val
			mask <<= 1
			continue
		}

		for prefixMaskSize := 7 - i; prefixMaskSize >= 0; prefixMaskSize-- {
			if rt.children[val].entries[prefixMaskSize] != nil {
				return rt.children[val].entries[prefixMaskSize]
			}
		}
		preVal = val
		mask <<= 1
	}
	return nil
}

func (rt *RadixTable) Lookup(ip string) *Entry {
	ipp := net.ParseIP(ip)
	if ipp == nil {
		return nil
	}
	return rt.LookupIP(ipp)
}

func (rt *RadixTable) LookupIP(ip net.IP) *Entry {
	// 转换为IPv4地址
	ip4 := ip.To4()
	if ip4 == nil {
		return nil // 仅处理IPv4
	}

	ipBytes := []byte(ip4)
	var nodes []*radixNode // nodes that may hit route
	var curNode *radixNode

	curNode = rt.root
	// Try to find the deepest node
	for _, curByte := range ipBytes {
		if curNode == nil {
			break
		}
		nodes = append(nodes, curNode)
		curNode = curNode.children[curByte]
	}

	// Try to traverse the query backwards starting from the deepest node
	for i := len(nodes) - 1; i >= 0; i-- {
		node := nodes[i]
		entry := node.lookupOneNode(ipBytes[i])
		if entry != nil {
			return entry
		}
	}

	return rt.defaultEntry
}
