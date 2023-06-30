package networks

import (
	"time"
)

type NetworkList struct {
	Name    string    `json:"name"`
	ID      string    `json:"id"`
	Driver  string    `json:"driver"`
	Scope   string    `json:"scope"`
	Created time.Time `json:"created"`
}

type NetworkCreateStruct struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	EnableIPAM bool              `json:"enableIPAM"`
	IPAM       IPAM              `json:"ipam"`
	Options    map[string]string `json:"options"`
	Labels     map[string]string `json:"labels"`
}

// IPAM represents IP Address Management
type IPAM struct {
	Driver  string            `json:"driver"`
	Options map[string]string `json:"options"` //Per network IPAM driver options
	Config  []IPAMConfig      `json:"config"`
}

// IPAMConfig represents IPAM configurations
type IPAMConfig struct {
	Subnet     string            `json:"subnet,omitempty"`
	IPRange    string            `json:"ipRange,omitempty"`
	Gateway    string            `json:"gateway,omitempty"`
	AuxAddress map[string]string `json:"AuxAddress,AuxiliaryAddresses,omitempty"`
}

type NetworkDeleteStruct struct {
	Name []string `json:"name"`
}
