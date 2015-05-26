// go generate gen.go
// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package ipv6

// Internet Control Message Protocol version 6 (ICMPv6) Parameters, Updated: 2014-09-22
const (
    ICMPTypeDestinationUnreachable                ICMPType = 1   // Destination Unreachable
    ICMPTypePacketTooBig                          ICMPType = 2   // Packet Too Big
    ICMPTypeTimeExceeded                          ICMPType = 3   // Time Exceeded
    ICMPTypeParameterProblem                      ICMPType = 4   // Parameter Problem
    ICMPTypeEchoRequest                           ICMPType = 128 // Echo Request
    ICMPTypeEchoReply                             ICMPType = 129 // Echo Reply
    ICMPTypeMulticastListenerQuery                ICMPType = 130 // Multicast Listener Query
    ICMPTypeMulticastListenerReport               ICMPType = 131 // Multicast Listener Report
    ICMPTypeMulticastListenerDone                 ICMPType = 132 // Multicast Listener Done
    ICMPTypeRouterSolicitation                    ICMPType = 133 // Router Solicitation
    ICMPTypeRouterAdvertisement                   ICMPType = 134 // Router Advertisement
    ICMPTypeNeighborSolicitation                  ICMPType = 135 // Neighbor Solicitation
    ICMPTypeNeighborAdvertisement                 ICMPType = 136 // Neighbor Advertisement
    ICMPTypeRedirect                              ICMPType = 137 // Redirect Message
    ICMPTypeRouterRenumbering                     ICMPType = 138 // Router Renumbering
    ICMPTypeNodeInformationQuery                  ICMPType = 139 // ICMP Node Information Query
    ICMPTypeNodeInformationResponse               ICMPType = 140 // ICMP Node Information Response
    ICMPTypeInverseNeighborDiscoverySolicitation  ICMPType = 141 // Inverse Neighbor Discovery Solicitation Message
    ICMPTypeInverseNeighborDiscoveryAdvertisement ICMPType = 142 // Inverse Neighbor Discovery Advertisement Message
    ICMPTypeVersion2MulticastListenerReport       ICMPType = 143 // Version 2 Multicast Listener Report
    ICMPTypeHomeAgentAddressDiscoveryRequest      ICMPType = 144 // Home Agent Address Discovery Request Message
    ICMPTypeHomeAgentAddressDiscoveryReply        ICMPType = 145 // Home Agent Address Discovery Reply Message
    ICMPTypeMobilePrefixSolicitation              ICMPType = 146 // Mobile Prefix Solicitation
    ICMPTypeMobilePrefixAdvertisement             ICMPType = 147 // Mobile Prefix Advertisement
    ICMPTypeCertificationPathSolicitation         ICMPType = 148 // Certification Path Solicitation Message
    ICMPTypeCertificationPathAdvertisement        ICMPType = 149 // Certification Path Advertisement Message
    ICMPTypeMulticastRouterAdvertisement          ICMPType = 151 // Multicast Router Advertisement
    ICMPTypeMulticastRouterSolicitation           ICMPType = 152 // Multicast Router Solicitation
    ICMPTypeMulticastRouterTermination            ICMPType = 153 // Multicast Router Termination
    ICMPTypeFMIPv6                                ICMPType = 154 // FMIPv6 Messages
    ICMPTypeRPLControl                            ICMPType = 155 // RPL Control Message
    ICMPTypeILNPv6LocatorUpdate                   ICMPType = 156 // ILNPv6 Locator Update Message
    ICMPTypeDuplicateAddressRequest               ICMPType = 157 // Duplicate Address Request
    ICMPTypeDuplicateAddressConfirmation          ICMPType = 158 // Duplicate Address Confirmation
)

// Internet Control Message Protocol version 6 (ICMPv6) Parameters, Updated: 2014-09-22
var icmpTypes = map[ICMPType]string{
    1:   "destination unreachable",
    2:   "packet too big",
    3:   "time exceeded",
    4:   "parameter problem",
    128: "echo request",
    129: "echo reply",
    130: "multicast listener query",
    131: "multicast listener report",
    132: "multicast listener done",
    133: "router solicitation",
    134: "router advertisement",
    135: "neighbor solicitation",
    136: "neighbor advertisement",
    137: "redirect message",
    138: "router renumbering",
    139: "icmp node information query",
    140: "icmp node information response",
    141: "inverse neighbor discovery solicitation message",
    142: "inverse neighbor discovery advertisement message",
    143: "version 2 multicast listener report",
    144: "home agent address discovery request message",
    145: "home agent address discovery reply message",
    146: "mobile prefix solicitation",
    147: "mobile prefix advertisement",
    148: "certification path solicitation message",
    149: "certification path advertisement message",
    151: "multicast router advertisement",
    152: "multicast router solicitation",
    153: "multicast router termination",
    154: "fmipv6 messages",
    155: "rpl control message",
    156: "ilnpv6 locator update message",
    157: "duplicate address request",
    158: "duplicate address confirmation",
}
