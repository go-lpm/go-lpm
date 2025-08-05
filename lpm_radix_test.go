package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
)

func TestRadixTable_Show(t *testing.T) {
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

		entries := table.Show()

		for maskLen := 0; maskLen < len(entries); maskLen++ {
			if maskLen == 0 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "0.0.0.0/0")
				So(entries[maskLen][0].Entry, ShouldEqual, 0)
			} else if maskLen == 8 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.0.0.0/8")
				So(entries[maskLen][0].Entry, ShouldEqual, 192)
			} else if maskLen == 16 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.0.0/16")
				So(entries[maskLen][0].Entry, ShouldEqual, 168)
			} else if maskLen == 24 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.1.0/24")
				So(entries[maskLen][0].Entry, ShouldEqual, 1)
			} else if maskLen == 32 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.1.255/32")
				So(entries[maskLen][0].Entry, ShouldEqual, 255)
			} else {
				So(len(entries[maskLen]), ShouldEqual, 0)
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

		entries := table.Show()

		for maskLen := 0; maskLen < len(entries); maskLen++ {
			if maskLen == 32 {
				So(len(entries[maskLen]), ShouldEqual, 4)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.1.0/32")
				So(entries[maskLen][0].Entry, ShouldEqual, 0)
				So(entries[maskLen][1].Prefix.String(), ShouldEqual, "192.168.1.7/32")
				So(entries[maskLen][1].Entry, ShouldEqual, 7)
				So(entries[maskLen][2].Prefix.String(), ShouldEqual, "192.168.1.250/32")
				So(entries[maskLen][2].Entry, ShouldEqual, 250)
				So(entries[maskLen][3].Prefix.String(), ShouldEqual, "192.168.1.255/32")
				So(entries[maskLen][3].Entry, ShouldEqual, 255)
			} else {
				So(len(entries[maskLen]), ShouldEqual, 0)
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

		entries := table.Show()

		for maskLen := 0; maskLen < len(entries); maskLen++ {
			if maskLen == 25 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.1.128/25")
				So(entries[maskLen][0].Entry, ShouldEqual, 128)
			} else if maskLen == 27 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.1.224/27")
				So(entries[maskLen][0].Entry, ShouldEqual, 224)
			} else if maskLen == 31 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.1.254/31")
				So(entries[maskLen][0].Entry, ShouldEqual, 254)
			} else if maskLen == 32 {
				So(len(entries[maskLen]), ShouldEqual, 1)
				So(entries[maskLen][0].Prefix.String(), ShouldEqual, "192.168.1.255/32")
				So(entries[maskLen][0].Entry, ShouldEqual, 255)
			} else {
				So(len(entries[maskLen]), ShouldEqual, 0)
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

	Convey("Add an overlay entry", t, func() {
		table := RadixTable{
			root: &radixNode{},
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
}

func TestRadixTable_Delete(t *testing.T) {
	var err error

	Convey("Delete invalid entry", t, func() {
		table := RadixTable{
			root: &radixNode{},
		}

		err = table.Delete("3.3.3.3/33")
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

		for maskLen := 0; maskLen <= 32; maskLen++ {
			_, cidr, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip0, maskLen))
			err = table.Delete(cidr.String())
			So(err, ShouldBeNil)
			_, cidr, _ = net.ParseCIDR(fmt.Sprintf("%s/%d", ip255, maskLen))
			err = table.Delete(cidr.String())
			So(err, ShouldBeNil)
		}

		entries = table.Show()

		for maskLen := 0; maskLen < len(entries); maskLen++ {
			So(len(entries[maskLen]), ShouldEqual, 0)
		}
	})
}

func TestRadixTable_Lookup(t *testing.T) {
	route := RadixTable{
		root: &radixNode{},
	}
	route.Add("192.168.0.0/24", "192.168.0.0/24")
	route.Add("192.168.0.1/32", "192.168.0.1/32")
	route.Add("192.168.0.2/32", "192.168.0.2/32")

	route.Add("172.16.0.0/12", "172.16.0.0/12")
	route.Add("172.16.0.4/30", "172.16.0.4/30")
	route.Add("172.16.0.12/30", "172.16.0.12/30")

	route.Add("10.0.0.0/8", "10.0.0.0/8")
	route.Add("10.10.0.0/17", "10.10.0.0/17")
	route.Add("10.10.128.0/17", "10.10.128.0/17")

	Convey("Lookup 192.168 private cidr", t, func() {
		entry := route.Lookup("192.168.0.1")
		So(entry.Entry, ShouldEqual, "192.168.0.1/32")
		entry = route.Lookup("192.168.0.2")
		So(entry.Entry, ShouldEqual, "192.168.0.2/32")
		entry = route.Lookup("192.168.0.3")
		So(entry.Entry, ShouldEqual, "192.168.0.0/24")
		entry = route.Lookup("192.168.1.3")
		So(entry, ShouldBeNil)
	})

	Convey("Lookup 172.16 private cidr", t, func() {
		entry := route.Lookup("172.16.0.4")
		So(entry.Entry, ShouldEqual, "172.16.0.4/30")
		entry = route.Lookup("172.16.0.5")
		So(entry.Entry, ShouldEqual, "172.16.0.4/30")
		entry = route.Lookup("172.16.0.6")
		So(entry.Entry, ShouldEqual, "172.16.0.4/30")
		entry = route.Lookup("172.16.0.8")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = route.Lookup("172.16.0.9")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = route.Lookup("172.16.0.10")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = route.Lookup("172.16.0.12")
		So(entry.Entry, ShouldEqual, "172.16.0.12/30")
		entry = route.Lookup("172.16.0.13")
		So(entry.Entry, ShouldEqual, "172.16.0.12/30")
		entry = route.Lookup("172.16.0.14")
		So(entry.Entry, ShouldEqual, "172.16.0.12/30")
		entry = route.Lookup("172.16.0.16")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = route.Lookup("172.16.0.17")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = route.Lookup("172.16.0.18")
		So(entry.Entry, ShouldEqual, "172.16.0.0/12")
		entry = route.Lookup("172.48.0.18")
		So(entry, ShouldBeNil)
	})

	Convey("Lookup 10 private cidr", t, func() {
		entry := route.Lookup("10.10.0.1")
		So(entry.Entry, ShouldEqual, "10.10.0.0/17")
		entry = route.Lookup("10.10.127.1")
		So(entry.Entry, ShouldEqual, "10.10.0.0/17")
		entry = route.Lookup("10.10.128.1")
		So(entry.Entry, ShouldEqual, "10.10.128.0/17")
		entry = route.Lookup("10.10.255.1")
		So(entry.Entry, ShouldEqual, "10.10.128.0/17")
		entry = route.Lookup("10.11.0.1")
		So(entry.Entry, ShouldEqual, "10.0.0.0/8")
		entry = route.Lookup("10.11.128.1")
		So(entry.Entry, ShouldEqual, "10.0.0.0/8")
		entry = route.Lookup("11.0.0.0")
		So(entry, ShouldBeNil)
	})

	Convey("Lookup default entries", t, func() {
		route.Add("0.0.0.0/0", "0.0.0.0/0")
		entry := route.Lookup("192.168.1.3")
		So(entry.Entry, ShouldEqual, "0.0.0.0/0")
		entry = route.Lookup("172.48.0.18")
		So(entry.Entry, ShouldEqual, "0.0.0.0/0")
		entry = route.Lookup("11.0.0.0")
		So(entry.Entry, ShouldEqual, "0.0.0.0/0")
	})
}
