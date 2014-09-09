// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package vmx

import "testing"

func TestParsingTag(t *testing.T) {
	tests := []struct {
		tag       string
		name      string
		omitempty bool
		err       string
	}{
		{"vmx:displayname", "", false, "Tag name has to be enclosed in double quotes: vmx:displayname"},
		{"vmx:", "", false, "Invalid tag: vmx:"},
		{`vmx:""`, "", false, `Tag name is missing: vmx:""`},
		{"vm", "", false, "Invalid tag: vm"},
		{`vmx:"displayname,omitempty`, "displayname", true, ""},
		{`vmx:"displayname,blah"`, "displayname", false, ""},
		{`vmx:"-"`, "-", false, ""},
	}

	for _, tt := range tests {
		name, omitempty, err := parseTag(tt.tag)
		equals(t, tt.name, name)
		equals(t, tt.omitempty, omitempty)
		if err != nil {
			equals(t, tt.err, err.Error())
		} else {
			equals(t, tt.err, "")
		}
	}
}

func TestMarshal(t *testing.T) {
	type VM struct {
		Encoding     string `vmx:".encoding"`
		Annotation   string `vmx:"annotation"`
		Hwversion    uint8  `vmx:"virtualHW.version"`
		HwProdCompat string `vmx:"virtualHW.productCompatibility"`
		Memsize      uint   `vmx:"memsize"`
		Numvcpus     uint   `vmx:"numvcpus"`
		MemHotAdd    bool   `vmx:"mem.hotadd"`
		DisplayName  string `vmx:"displayName"`
		GuestOS      string `vmx:"guestOS"`
		Autoanswer   bool   `vmx:"msg.autoAnswer"`
	}

	vm := new(VM)
	vm.Encoding = "utf-8"
	vm.Annotation = "Test VM"
	vm.Hwversion = 10
	vm.HwProdCompat = "hosted"
	vm.Memsize = 1024
	vm.Numvcpus = 2
	vm.MemHotAdd = false
	vm.DisplayName = "test"
	vm.GuestOS = "other3xlinux-64"
	vm.Autoanswer = true

	data, err := Marshal(vm)
	ok(t, err)
	expected := `.encoding = "utf-8"
annotation = "Test VM"
virtualHW.version = "10"
virtualHW.productCompatibility = "hosted"
memsize = "1024"
numvcpus = "2"
mem.hotadd = "false"
displayName = "test"
guestOS = "other3xlinux-64"
msg.autoAnswer = "true"
`
	equals(t, expected, string(data))
}

func TestMarshalEmbedded(t *testing.T) {
	type Vhardware struct {
		Version string `vmx:"version"`
		Compat  string `vmx:"productCompatibility"`
	}

	type VM struct {
		Encoding    string    `vmx:".encoding"`
		Annotation  string    `vmx:"annotation"`
		Vhardware   Vhardware `vmx:"virtualHW"`
		Memsize     uint      `vmx:"memsize"`
		Numvcpus    uint      `vmx:"numvcpus"`
		MemHotAdd   bool      `vmx:"mem.hotadd"`
		DisplayName string    `vmx:"displayName"`
		GuestOS     string    `vmx:"guestOS"`
		Autoanswer  bool      `vmx:"msg.autoAnswer"`
	}

	vm := new(VM)
	vm.Encoding = "utf-8"
	vm.Annotation = "Test VM"
	vm.Vhardware = Vhardware{
		Version: "10",
		Compat:  "hosted",
	}
	vm.Memsize = 1024
	vm.Numvcpus = 2
	vm.MemHotAdd = false
	vm.DisplayName = "test"
	vm.GuestOS = "other3xlinux-64"
	vm.Autoanswer = true

	data, err := Marshal(vm)
	ok(t, err)
	expected := `.encoding = "utf-8"
annotation = "Test VM"
virtualHW.version = "10"
virtualHW.productCompatibility = "hosted"
memsize = "1024"
numvcpus = "2"
mem.hotadd = "false"
displayName = "test"
guestOS = "other3xlinux-64"
msg.autoAnswer = "true"
`
	equals(t, expected, string(data))
}

func TestMarshalArray(t *testing.T) {
	type Vhardware struct {
		Version string `vmx:"version"`
		Compat  string `vmx:"productCompatibility"`
	}

	type Ethernet struct {
		StartConnected       bool   `vmx:"startConnected"`
		Present              bool   `vmx:"present"`
		ConnectionType       string `vmx:"connectionType"`
		VirtualDev           string `vmx:"virtualDev"`
		WakeOnPcktRcv        bool   `vmx:"wakeOnPcktRcv"`
		AddressType          string `vmx:"addressType"`
		LinkStatePropagation bool   `vmx:"linkStatePropagation.enable,omitempty"`
	}

	type VM struct {
		Encoding    string     `vmx:".encoding"`
		Annotation  string     `vmx:"annotation"`
		Vhardware   Vhardware  `vmx:"virtualHW"`
		Memsize     uint       `vmx:"memsize"`
		Numvcpus    uint       `vmx:"numvcpus"`
		MemHotAdd   bool       `vmx:"mem.hotadd"`
		DisplayName string     `vmx:"displayName"`
		GuestOS     string     `vmx:"guestOS"`
		Autoanswer  bool       `vmx:"msg.autoAnswer"`
		Ethernet    []Ethernet `vmx:"ethernet"`
	}

	vm := new(VM)
	vm.Encoding = "utf-8"
	vm.Annotation = "Test VM"
	vm.Vhardware = Vhardware{
		Version: "9",
		Compat:  "hosted",
	}
	vm.Ethernet = []Ethernet{
		{
			StartConnected:       true,
			Present:              true,
			ConnectionType:       "bridged",
			VirtualDev:           "e1000",
			WakeOnPcktRcv:        false,
			AddressType:          "generated",
			LinkStatePropagation: true,
		},
		{
			StartConnected: true,
			Present:        true,
			ConnectionType: "nat",
			VirtualDev:     "e1000",
			WakeOnPcktRcv:  false,
			AddressType:    "generated",
		},
	}
	vm.Memsize = 1024
	vm.Numvcpus = 2
	vm.MemHotAdd = false
	vm.DisplayName = "test"
	vm.GuestOS = "other3xlinux-64"
	vm.Autoanswer = true

	data, err := Marshal(vm)
	ok(t, err)
	expected := `.encoding = "utf-8"
annotation = "Test VM"
virtualHW.version = "9"
virtualHW.productCompatibility = "hosted"
memsize = "1024"
numvcpus = "2"
mem.hotadd = "false"
displayName = "test"
guestOS = "other3xlinux-64"
msg.autoAnswer = "true"
ethernet0.startConnected = "true"
ethernet0.present = "true"
ethernet0.connectionType = "bridged"
ethernet0.virtualDev = "e1000"
ethernet0.wakeOnPcktRcv = "false"
ethernet0.addressType = "generated"
ethernet0.linkStatePropagation.enable = "true"
ethernet1.startConnected = "true"
ethernet1.present = "true"
ethernet1.connectionType = "nat"
ethernet1.virtualDev = "e1000"
ethernet1.wakeOnPcktRcv = "false"
ethernet1.addressType = "generated"
`
	equals(t, expected, string(data))
}
