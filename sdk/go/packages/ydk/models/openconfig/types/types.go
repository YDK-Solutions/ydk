package types

//////////////////////////////////////////////////////////////////////////
// Ipv6UnicastIdentity
//////////////////////////////////////////////////////////////////////////
type Ipv6UnicastIdentity struct {
}

func (id Ipv6UnicastIdentity) String() string {
	return "openconfig-bgp-types:IPV6_UNICAST"
}

type Ipv4UnicastIdentity struct {
}

func (id Ipv4UnicastIdentity) String() string {
	return "openconfig-bgp-types:IPV4_UNICAST"
}
