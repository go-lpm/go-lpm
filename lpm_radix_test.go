package golpm

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
)

func TestRadixTable_Show(t *testing.T) {
	var err error

	Convey("Show nodes across multiple layers", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}
		_, cidr, _ := net.ParseCIDR("0.0.0.0/0")
		table.defaultEntry = &Entry{
			Prefix: cidr,
			Entry:  0,
		}
		_, cidr, _ = net.ParseCIDR("192.0.0.0/8")
		layer1 := table.root
		layer1.children[192] = &radixNode{}
		layer1.children[192].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  192,
		}
		_, cidr, _ = net.ParseCIDR("192.168.0.0/16")
		layer2 := layer1.children[192]
		layer2.children[168] = &radixNode{}
		layer2.children[168].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  168,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.0/24")
		layer3 := layer2.children[168]
		layer3.children[1] = &radixNode{}
		layer3.children[1].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  1,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.255/32")
		layer4 := layer3.children[1]
		layer4.children[255] = &radixNode{}
		layer4.children[255].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  255,
		}

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 0 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "0.0.0.0/0")
				So(entries[0].Entry, ShouldEqual, 0)
			} else if maskLen == 8 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.0.0.0/8")
				So(entries[0].Entry, ShouldEqual, 192)
			} else if maskLen == 16 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.0.0/16")
				So(entries[0].Entry, ShouldEqual, 168)
			} else if maskLen == 24 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.1.0/24")
				So(entries[0].Entry, ShouldEqual, 1)
			} else if maskLen == 32 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.1.255/32")
				So(entries[0].Entry, ShouldEqual, 255)
			} else {
				So(len(entries), ShouldEqual, 0)
			}
		}
	})

	Convey("Show nodes in same layers", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}
		layer1 := table.root
		layer1.children[192] = &radixNode{}
		layer2 := layer1.children[192]
		layer2.children[168] = &radixNode{}
		layer3 := layer2.children[168]
		layer3.children[1] = &radixNode{}
		layer4 := layer3.children[1]

		_, cidr, _ := net.ParseCIDR("192.168.1.0/32")
		layer4.children[0] = &radixNode{}
		layer4.children[0].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  0,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.7/32")
		layer4.children[7] = &radixNode{}
		layer4.children[7].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  7,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.250/32")
		layer4.children[250] = &radixNode{}
		layer4.children[250].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  250,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.255/32")
		layer4.children[255] = &radixNode{}
		layer4.children[255].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  255,
		}

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 32 {
				So(len(entries), ShouldEqual, 4)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.1.0/32")
				So(entries[0].Entry, ShouldEqual, 0)
				So(entries[1].Prefix.String(), ShouldEqual, "192.168.1.7/32")
				So(entries[1].Entry, ShouldEqual, 7)
				So(entries[2].Prefix.String(), ShouldEqual, "192.168.1.250/32")
				So(entries[2].Entry, ShouldEqual, 250)
				So(entries[3].Prefix.String(), ShouldEqual, "192.168.1.255/32")
				So(entries[3].Entry, ShouldEqual, 255)
			} else {
				So(len(entries), ShouldEqual, 0)
			}
		}
	})

	Convey("Show multiple entries in same node", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}
		layer1 := table.root
		layer1.children[192] = &radixNode{}
		layer2 := layer1.children[192]
		layer2.children[168] = &radixNode{}
		layer3 := layer2.children[168]
		layer3.children[1] = &radixNode{}
		layer4 := layer3.children[1]

		_, cidr, _ := net.ParseCIDR("192.168.1.128/25")
		layer4.children[1] = &radixNode{}
		layer4.children[1].entries[0] = &Entry{
			Prefix: cidr,
			Entry:  128,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.224/27")
		layer4.children[7] = &radixNode{}
		layer4.children[7].entries[2] = &Entry{
			Prefix: cidr,
			Entry:  224,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.254/31")
		layer4.children[254] = &radixNode{}
		layer4.children[254].entries[6] = &Entry{
			Prefix: cidr,
			Entry:  254,
		}
		_, cidr, _ = net.ParseCIDR("192.168.1.255/32")
		layer4.children[255] = &radixNode{}
		layer4.children[255].entries[7] = &Entry{
			Prefix: cidr,
			Entry:  255,
		}

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 25 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.1.128/25")
				So(entries[0].Entry, ShouldEqual, 128)
			} else if maskLen == 27 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.1.224/27")
				So(entries[0].Entry, ShouldEqual, 224)
			} else if maskLen == 31 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.1.254/31")
				So(entries[0].Entry, ShouldEqual, 254)
			} else if maskLen == 32 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "192.168.1.255/32")
				So(entries[0].Entry, ShouldEqual, 255)
			} else {
				So(len(entries), ShouldEqual, 0)
			}
		}
	})

	Convey("Show ipv6 nodes across multiple layers", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		ip := "2406:d440:202:f01::ffff:ffff"

		for maskLen := 0; maskLen <= 128; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip, maskLen))
			err = table.Add(cidr.String(), maskLen)
			So(err, ShouldBeNil)
		}

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip, maskLen))
			So(len(entries), ShouldEqual, 1)
			So(entries[0].Prefix.String(), ShouldEqual, cidr.String())
			So(entries[0].Entry, ShouldEqual, maskLen)
		}
	})

	Convey("Show ipv6 nodes in same layers", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		err = table.Add("2406:d440:202:f01::fff/128", "fff")
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:f01::ffe/128", "ffe")
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:f01::ff7/128", "ff7")
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:f01::f7f/128", "f7f")
		So(err, ShouldBeNil)

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 128 {
				So(len(entries), ShouldEqual, 4)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:f01::f7f/128")
				So(entries[0].Entry, ShouldEqual, "f7f")
				So(entries[1].Prefix.String(), ShouldEqual, "2406:d440:202:f01::ff7/128")
				So(entries[1].Entry, ShouldEqual, "ff7")
				So(entries[2].Prefix.String(), ShouldEqual, "2406:d440:202:f01::ffe/128")
				So(entries[2].Entry, ShouldEqual, "ffe")
				So(entries[3].Prefix.String(), ShouldEqual, "2406:d440:202:f01::fff/128")
				So(entries[3].Entry, ShouldEqual, "fff")
			} else {
				So(len(entries), ShouldEqual, 0)
			}
		}
	})

	Convey("Show multiple ipv6 entries in same node", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		err = table.Add("2406:d440:202:fff::/64", 64)
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:ffe::/63", 63)
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:ff8::/61", 61)
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:f80::/57", 57)
		So(err, ShouldBeNil)

		err = table.Add("2406:d440:202:f01::fff/128", 128)
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:f01::ffe/127", 127)
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:f01::ff8/125", 125)
		So(err, ShouldBeNil)
		err = table.Add("2406:d440:202:f01::f80/121", 121)
		So(err, ShouldBeNil)

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 57 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:f80::/57")
				So(entries[0].Entry, ShouldEqual, 57)
			} else if maskLen == 61 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:ff8::/61")
				So(entries[0].Entry, ShouldEqual, 61)
			} else if maskLen == 63 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:ffe::/63")
				So(entries[0].Entry, ShouldEqual, 63)
			} else if maskLen == 64 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:fff::/64")
				So(entries[0].Entry, ShouldEqual, 64)
			} else if maskLen == 121 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:f01::f80/121")
				So(entries[0].Entry, ShouldEqual, 121)
			} else if maskLen == 125 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:f01::ff8/125")
				So(entries[0].Entry, ShouldEqual, 125)
			} else if maskLen == 127 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:f01::ffe/127")
				So(entries[0].Entry, ShouldEqual, 127)
			} else if maskLen == 128 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "2406:d440:202:f01::fff/128")
				So(entries[0].Entry, ShouldEqual, 128)
			} else {
				So(len(entries), ShouldEqual, 0)
			}
		}
	})
}

func TestRadixTable_Add(t *testing.T) {
	var err error

	Convey("Add invalid prefix", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}
		err = table.Add("3.3.3.3/33", 3)
		So(err, ShouldBeError)
		err = table.Add("1234::1::1/128", 128)
		So(err, ShouldBeError)
		err = table.Add("1234::1/129", 129)
		So(err, ShouldBeError)
	})

	Convey("Add ipv6 entry to ipv4 table", t, func() {
		table := RadixTable{
			ipBytesLen: net.IPv4len,
			root:       &radixNode{},
		}
		err = table.Add("1234::1/128", 128)
		So(err, ShouldBeError)
	})

	Convey("Add ipv4 entry to ipv6 table", t, func() {
		table := RadixTable{
			ipBytesLen: net.IPv6len,
			root:       &radixNode{},
		}
		err = table.Add("192.168.0.1/32", 32)
		So(err, ShouldBeError)
	})

	Convey("Add table entries of all mask len", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}
		ip0 := "0.0.0.0"
		ip255 := "255.255.255.255"
		for maskLen := 0; maskLen <= 32; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
			err = table.Add(cidr.String(), cidr.String())
			So(err, ShouldBeNil)
			_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ip255, maskLen))
			err = table.Add(cidr.String(), cidr.String())
			So(err, ShouldBeNil)
		}

		entries := table.Show()

		for maskLen := 0; maskLen < len(entries); maskLen++ {
			if maskLen == 0 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "0.0.0.0/0")
				So(entries[maskLen][0].Entry, ShouldEqual, "0.0.0.0/0")
			} else {
				So(len(entries[maskLen]), ShouldEqual, 2)
				_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[maskLen][0].Entry, ShouldEqual, cidr.String())
				_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ip255, maskLen))
				So(entries[maskLen][1].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[maskLen][1].Entry, ShouldEqual, cidr.String())
			}
		}
	})

	Convey("Add ipv6 table entries of all mask len", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		ip1 := "2406:d440:202:f01::ffff:ffff"
		ipf := "ffff:ffff:fff:fff::ffff:ffff"

		for maskLen := 0; maskLen <= 128; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip1, maskLen))
			err = table.Add(cidr.String(), maskLen)
			So(err, ShouldBeNil)
			_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ipf, maskLen))
			err = table.Add(cidr.String(), maskLen)
			So(err, ShouldBeNil)
		}

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 0 {
				_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip1, maskLen))
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[0].Entry, ShouldEqual, maskLen)
			} else {
				So(len(entries), ShouldEqual, 2)
				_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip1, maskLen))
				So(entries[0].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[0].Entry, ShouldEqual, maskLen)
				_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ipf, maskLen))
				So(entries[1].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[1].Entry, ShouldEqual, maskLen)
			}
		}
	})

	Convey("Add an overlay entry", t, func() {
		table := RadixTable{
			ipBytesLen: net.IPv4len,
			root:       &radixNode{},
		}
		err = table.Add("192.168.0.1/32", 1)
		So(err, ShouldBeNil)
		entries := table.Show()
		So(entries[32][0].Entry, ShouldEqual, 1)
		err = table.Add("192.168.0.1/32", 2)
		So(err, ShouldBeNil)
		entries = table.Show()
		So(entries[32][0].Entry, ShouldEqual, 2)
	})

	Convey("Add an overlay ipv6 entry", t, func() {
		table := RadixTable{
			ipBytesLen: net.IPv6len,
			root:       &radixNode{},
		}
		err = table.Add("1234::1/128", 1)
		So(err, ShouldBeNil)
		entries := table.Show()
		So(entries[128][0].Entry, ShouldEqual, 1)
		err = table.Add("1234::1/128", 2)
		So(err, ShouldBeNil)
		entries = table.Show()
		So(entries[128][0].Entry, ShouldEqual, 2)
	})
}

func TestRadixTable_Delete(t *testing.T) {
	var err error

	Convey("Delete invalid entry", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		err = table.Delete("3.3.3.3/33")
		So(err, ShouldBeError)
		err = table.Delete("1234::1::1/128")
		So(err, ShouldBeError)
		err = table.Delete("1234::1/129")
		So(err, ShouldBeError)
	})

	Convey("Delete non-exist entry", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		err = table.Delete("1.1.1.1/32")
		So(err, ShouldBeNil)
		entries := table.Show()
		for maskLen := 0; maskLen < len(entries); maskLen++ {
			So(len(entries[maskLen]), ShouldEqual, 0)
		}
	})

	Convey("Delete ipv6 entry to ipv4 table", t, func() {
		table := RadixTable{
			ipBytesLen: net.IPv4len,
			root:       &radixNode{},
		}
		err = table.Delete("1234::1/128")
		So(err, ShouldBeError)
	})

	Convey("Delete ipv4 entry to ipv6 table", t, func() {
		table := RadixTable{
			ipBytesLen: net.IPv6len,
			root:       &radixNode{},
		}
		err = table.Delete("192.168.0.1/32")
		So(err, ShouldBeError)
	})

	Convey("Delete entry", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		ip0 := "0.0.0.0"
		ip255 := "255.255.255.255"
		for maskLen := 0; maskLen <= 32; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
			err = table.Add(cidr.String(), cidr.String())
			So(err, ShouldBeNil)
			_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ip255, maskLen))
			err = table.Add(cidr.String(), cidr.String())
			So(err, ShouldBeNil)
		}

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 0 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "0.0.0.0/0")
				So(entries[0].Entry, ShouldEqual, "0.0.0.0/0")
			} else {
				So(len(entries), ShouldEqual, 2)
				_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
				So(entries[0].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[0].Entry, ShouldEqual, cidr.String())
				_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ip255, maskLen))
				So(entries[1].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[1].Entry, ShouldEqual, cidr.String())
			}
		}

		for maskLen := 0; maskLen <= 32; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
			err = table.Delete(cidr.String())
			So(err, ShouldBeNil)
			_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ip255, maskLen))
			err = table.Delete(cidr.String())
			So(err, ShouldBeNil)
		}

		maskLen2entries = table.Show()

		for _, entries := range maskLen2entries {
			So(len(entries), ShouldEqual, 0)
		}
	})

	Convey("Delete ipv6 entry", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		ip0 := "::"
		ipf := "ffff:ffff:ffff::ffff:ffff:ffff"
		for maskLen := 0; maskLen <= 128; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
			err = table.Add(cidr.String(), cidr.String())
			So(err, ShouldBeNil)
			_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ipf, maskLen))
			err = table.Add(cidr.String(), cidr.String())
			So(err, ShouldBeNil)
		}

		maskLen2entries := table.Show()

		for maskLen, entries := range maskLen2entries {
			if maskLen == 0 {
				So(len(entries), ShouldEqual, 1)
				So(entries[0].Prefix.String(), ShouldEqual, "::/0")
				So(entries[0].Entry, ShouldEqual, "::/0")
			} else {
				So(len(entries), ShouldEqual, 2)
				_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
				So(entries[0].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[0].Entry, ShouldEqual, cidr.String())
				_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ipf, maskLen))
				So(entries[1].Prefix.String(), ShouldEqual, cidr.String())
				So(entries[1].Entry, ShouldEqual, cidr.String())
			}
		}

		for maskLen := 0; maskLen <= 128; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
			err = table.Delete(cidr.String())
			So(err, ShouldBeNil)
			_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ipf, maskLen))
			err = table.Delete(cidr.String())
			So(err, ShouldBeNil)
		}

		maskLen2entries = table.Show()

		for _, entries := range maskLen2entries {
			So(len(entries), ShouldEqual, 0)
		}
	})
}

func TestRadixTable_Lookup(t *testing.T) {
	table := RadixTable{
		ipBytesLen: net.IPv4len,
		root:       &radixNode{},
	}
	table.Add("192.168.0.0/24", "192.168.0.0/24")
	table.Add("192.168.0.1/32", "192.168.0.1/32")
	table.Add("192.168.0.2/32", "192.168.0.2/32")

	table.Add("172.16.0.0/12", "172.16.0.0/12")
	table.Add("172.16.0.4/30", "172.16.0.4/30")
	table.Add("172.16.0.12/30", "172.16.0.12/30")

	table.Add("10.0.0.0/8", "10.0.0.0/8")
	table.Add("10.10.0.0/17", "10.10.0.0/17")
	table.Add("10.10.128.0/17", "10.10.128.0/17")

	Convey("Lookup 192.168 private cidr", t, func() {
		entry := table.Lookup("192.168.0.1")
		So(entry.Entry, ShouldEqual, "192.168.0.1/32")
		entry = table.Lookup("192.168.0.2")
		So(entry.Entry, ShouldEqual, "192.168.0.2/32")
		entry = table.Lookup("192.168.0.3")
		So(entry.Entry, ShouldEqual, "192.168.0.0/24")
		entry = table.Lookup("192.168.1.3")
		So(entry, ShouldBeNil)
	})

	Convey("Lookup 172.16 private cidr", t, func() {
		entry := table.Lookup("172.16.0.4")
		So(entry.Entry, ShouldEqual, "172.16.0.4/30")
		entry = table.Lookup("172.16.0.5")
		So(entry.Entry, ShouldEqual, "172.16.0.4/30")
		entry = table.Lookup("172.16.0.6")
		So(entry.Entry, ShouldEqual, "172.16.0.4/30")
		entry = table.Lookup("172.16.0.8")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = table.Lookup("172.16.0.9")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = table.Lookup("172.16.0.10")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = table.Lookup("172.16.0.12")
		So(entry.Entry, ShouldEqual, "172.16.0.12/30")
		entry = table.Lookup("172.16.0.13")
		So(entry.Entry, ShouldEqual, "172.16.0.12/30")
		entry = table.Lookup("172.16.0.14")
		So(entry.Entry, ShouldEqual, "172.16.0.12/30")
		entry = table.Lookup("172.16.0.16")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = table.Lookup("172.16.0.17")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = table.Lookup("172.16.0.18")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = table.Lookup("172.48.0.18")
		So(entry, ShouldBeNil)
	})

	Convey("Lookup 10 private cidr", t, func() {
		entry := table.Lookup("10.10.0.1")
		So(entry.Entry, ShouldEqual, "10.10.0.0/17")
		entry = table.Lookup("10.10.127.1")
		So(entry.Entry, ShouldEqual, "10.10.0.0/17")
		entry = table.Lookup("10.10.128.1")
		So(entry.Entry, ShouldEqual, "10.10.128.0/17")
		entry = table.Lookup("10.10.255.1")
		So(entry.Entry, ShouldEqual, "10.10.128.0/17")
		entry = table.Lookup("10.11.0.1")
		So(entry.Entry, ShouldEqual, "10.0.0.0/8")
		entry = table.Lookup("10.11.128.1")
		So(entry.Entry, ShouldEqual, "10.0.0.0/8")
		entry = table.Lookup("11.0.0.0")
		So(entry, ShouldBeNil)
	})

	Convey("Lookup default entries", t, func() {
		table.Add("0.0.0.0/0", "0.0.0.0/0")
		entry := table.Lookup("192.168.1.3")
		So(entry.Entry, ShouldEqual, "0.0.0.0/0")
		entry = table.Lookup("172.48.0.18")
		So(entry.Entry, ShouldEqual, "0.0.0.0/0")
		entry = table.Lookup("11.0.0.0")
		So(entry.Entry, ShouldEqual, "0.0.0.0/0")
		table.Delete("0.0.0.0/0")
	})
}

func TestRadixTable_Lookup_ipv6(t *testing.T) {
	table := RadixTable{
		ipBytesLen: net.IPv6len,
		root:       &radixNode{},
	}
	table.Add("2406:d440:202:f01::ffff:ff00/120", "2406:d440:202:f01::ffff:ff00/120")
	table.Add("2406:d440:202:f01::ffff:ffff/128", "2406:d440:202:f01::ffff:ffff/128")
	table.Add("2406:d440:202:f01::ffff:fffe/128", "2406:d440:202:f01::ffff:fffe/128")

	table.Add("2406:d440:202:f01::/100", "2406:d440:202:f01::/100")
	table.Add("2406:d440:202:f01::f00:0/105", "2406:d440:202:f01::f00:0/105")
	table.Add("2406:d440:202:f01::c00:0/105", "2406:d440:202:f01::c00:0/105")

	table.Add("2406:d440:202:8000::/49", "2406:d440:202:8000::/49")
	table.Add("2406:d440:200::/49", "2406:d440:200::/49")
	table.Add("2406:d440:200::/40", "2406:d440:200::/40")

	Convey("Lookup 120 cidr", t, func() {
		entry := table.Lookup("2406:d440:202:f01::ffff:ffff")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::ffff:ffff/128")
		entry = table.Lookup("2406:d440:202:f01::ffff:fffe")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::ffff:fffe/128")
		entry = table.Lookup("2406:d440:202:f01::ffff:fffd")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::ffff:ff00/120")
	})

	Convey("Lookup 100 cidr", t, func() {
		entry := table.Lookup("2406:d440:202:f01::f00:0")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::f00:0/105")
		entry = table.Lookup("2406:d440:202:f01::f00:f")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::f00:0/105")
		entry = table.Lookup("2406:d440:202:f01::f00:f00")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::f00:0/105")

		entry = table.Lookup("2406:d440:202:f01::c00:0")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::c00:0/105")
		entry = table.Lookup("2406:d440:202:f01::c00:f0")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::c00:0/105")
		entry = table.Lookup("2406:d440:202:f01::c00:f000")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::c00:0/105")

		entry = table.Lookup("2406:d440:202:f01::a00:0")
		So(entry.Entry, ShouldEqual, "2406:d440:202:f01::/100")
	})

	Convey("Lookup 40 cidr", t, func() {
		entry := table.Lookup("2406:d440:202:8000::")
		So(entry.Entry, ShouldEqual, "2406:d440:202:8000::/49")
		entry = table.Lookup("2406:d440:202:8000::f0f0")
		So(entry.Entry, ShouldEqual, "2406:d440:202:8000::/49")
		entry = table.Lookup("2406:d440:202:8f0f::")
		So(entry.Entry, ShouldEqual, "2406:d440:202:8000::/49")

		entry = table.Lookup("2406:d440:200::")
		So(entry.Entry, ShouldEqual, "2406:d440:200::/49")
		entry = table.Lookup("2406:d440:200:f0f::")
		So(entry.Entry, ShouldEqual, "2406:d440:200::/49")
		entry = table.Lookup("2406:d440:200:7777::")
		So(entry.Entry, ShouldEqual, "2406:d440:200::/49")
		entry = table.Lookup("2406:d440:200:8000::")
		So(entry.Entry, ShouldEqual, "2406:d440:200::/40")
	})

	Convey("Lookup default entries", t, func() {
		table.Add("::/0", "::/0")
		entry := table.Lookup("1234::")
		So(entry.Entry, ShouldEqual, "::/0")
		table.Delete("::/0")
		entry = table.Lookup("1234::")
		So(entry, ShouldBeNil)
	})
}
