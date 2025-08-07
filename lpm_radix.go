package golpm

import (
	"errors"
	"math"
	"net"
)

type RadixTable struct {
	root         *radixNode
	defaultEntry *Entry
	ipBytesLen   int
}

type radixNode struct {
	childCnt int
	entryCnt int
	children [256]*radixNode
	entries  [8]*Entry // routing entries stored by the node
}

func traverse(deep int, node *radixNode, entries map[int][]Entry) {
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

// Show Return the lpm table in format: maskLen -> entry list
func (rt *RadixTable) Show() map[int][]Entry {
	var entries map[int][]Entry
	entries = make(map[int][]Entry)
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
	if rt.ipBytesLen == net.IPv4len && prefix.IP.To4() == nil {
		return errors.New("add ipv6 entry to ipv4 table")
	} else if rt.ipBytesLen == net.IPv6len && prefix.IP.To4() != nil {
		return errors.New("add ipv4 entry to ipv6 table")
	}

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
			curNode.childCnt++
			curNode.children[curByte] = &radixNode{}
		}
		curNode = curNode.children[curByte]
	}
	// save entry in end point
	entryIdx := (maskSize + 7) % 8
	if curNode.entries[entryIdx] == nil {
		curNode.entryCnt++
	}
	curNode.entries[entryIdx] = &Entry{
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
	if rt.ipBytesLen == net.IPv4len && prefix.IP.To4() == nil {
		return errors.New("delete ipv6 entry to ipv4 table")
	} else if rt.ipBytesLen == net.IPv6len && prefix.IP.To4() != nil {
		return errors.New("delete ipv4 entry to ipv6 table")
	}

	maskSize, _ := prefix.Mask.Size()
	if maskSize == 0 {
		rt.defaultEntry = nil
		return nil
	}

	var nodePath []*radixNode
	var curNode *radixNode
	var curByte byte

	byteCount := (maskSize + 7) / 8
	ipBytes := []byte(prefix.IP)

	curNode = rt.root

	// find corresponding node byte-by-byte
	for i := 0; i < byteCount; i++ {
		nodePath = append(nodePath, curNode)
		curByte = ipBytes[i]
		if curNode.children[curByte] == nil {
			return nil
		}
		curNode = curNode.children[curByte]
	}
	// delete entry from end point
	entryIdx := (maskSize + 7) % 8
	if curNode.entries[entryIdx] != nil {
		curNode.entryCnt--
		curNode.entries[entryIdx] = nil
	}
	// free the node memory when appropriate
	if curNode.entryCnt != 0 {
		return nil
	}
	for i := byteCount - 1; i >= 0; i-- {
		curByte = ipBytes[i]
		if nodePath[i].children[curByte].childCnt == 0 && nodePath[i].children[curByte].entryCnt == 0 {
			nodePath[i].children[curByte] = nil
			nodePath[i].childCnt--
		}
	}
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
	if rt.ipBytesLen == net.IPv4len && ip.To4() == nil ||
		rt.ipBytesLen == net.IPv6len && ip.To4() != nil {
		return nil
	}

	ipBytes := []byte(ip)
	if rt.ipBytesLen == net.IPv4len {
		ipBytes = ip.To4()
	}
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
