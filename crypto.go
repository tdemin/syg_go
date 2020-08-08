package main

import (
	"github.com/yggdrasil-network/yggdrasil-go/src/address"
	"github.com/yggdrasil-network/yggdrasil-go/src/crypto"
)

// AddrForNodeID is a variant of Yggdrasil's src/address.AddrForNodeID that
// might be slightly optimized for performance.
//
// This function is a modded variant of address.AddrForNodeID from Yggdrasil.
// See src/address/address.go@78b5f88e4bb734d0dd6a138ff08d34ca39dcaea3
func AddrForNodeID(nid *crypto.NodeID) *address.Address {
	// 128 bit address, begins with GetPrefix(), with last bit set to 0
	// (indicates an address). Next 7 bits, interpreted as a uint, are the count
	// of leading 1s in the NodeID. Leading 1s and first leading 0 of the NodeID
	// are truncated off. The rest is appended to the IPv6 address (truncated to
	// 128 bits total).
	var addr address.Address
	temp := make([]byte, 0, len(nid))
	done := false
	ones := byte(0)
	bits := byte(0)
	nBits := 0
	for idx := 0; idx < 8*len(nid); idx++ {
		bit := (nid[idx/8] & (0x80 >> byte(idx%8))) >> byte(7-(idx%8))
		if !done && bit != 0 {
			ones++
			continue
		}
		if !done && bit == 0 {
			done = true
			continue
		}
		bits = (bits << 1) | bit
		nBits++
		if nBits == 8 {
			nBits = 0
			temp = append(temp, bits)
		}
	}
	prefix := address.GetPrefix()
	copy(addr[:], prefix[:])
	addr[len(prefix)] = ones
	copy(addr[len(prefix)+1:], temp)
	return &addr
}
