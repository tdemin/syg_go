package main

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/yggdrasil-network/yggdrasil-go/src/address"
)

// custom function selectors, see crypto_xxx.go
var (
	addrForKey  func(ed25519.PublicKey) *address.Address
	generateKey func() (ed25519.PublicKey, ed25519.PrivateKey)
)

// copied from Yggdrasil in order to save time on GetPrefix() call
var prefix = [...]byte{0x02}

// AddrForKey is a variant of Yggdrasil's src/address.AddrForKey that
// might be slightly optimized for performance.
//
// This function is a modded variant of address.AddrForKey from Yggdrasil.
// See src/address/address.go@0cff56fcc17d1acaf5297a7024477a9ca1bd3590
func AddrForKey(key ed25519.PublicKey) *address.Address {
	var buf [ed25519.PublicKeySize]byte
	copy(buf[:], key)
	for idx := range buf {
		buf[idx] = ^buf[idx]
	}
	var addr address.Address
	var temp = make([]byte, 0, 32)
	done := false
	ones := byte(0)
	bits := byte(0)
	nBits := 0
	for idx := 0; idx < 8*len(buf); idx++ {
		bit := (buf[idx/8] & (0x80 >> byte(idx%8))) >> byte(7-(idx%8))
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
	copy(addr[:], prefix[:])
	addr[len(prefix)] = ones
	copy(addr[len(prefix)+1:], temp)
	return &addr
}

// GenerateKey generates a ed25519 public/private key pair using entropy from
// crypto/rand. It assumes crypto/rand always blocks in order to bypass some
// checks to spend less time generating keys.
//
// This function is a modded variant of GenerateKey from crypto/ed25519 built-in
// Go package from Go 1.16.5.
func GenerateKey() (ed25519.PublicKey, ed25519.PrivateKey) {
	seed := make([]byte, ed25519.SeedSize)
	rand.Reader.Read(seed)

	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privateKey[32:])

	return publicKey, privateKey
}

// GenerateKeyEd25519 wraps crypto/ed25519.GenerateKey(). It panics if errors
// are encountered.
func GenerateKeyEd25519() (ed25519.PublicKey, ed25519.PrivateKey) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	return pub, priv
}
