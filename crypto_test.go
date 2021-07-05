package main

import (
	"crypto/ed25519"
	"net"
	"os"
	"regexp"
	"testing"

	"github.com/yggdrasil-network/yggdrasil-go/src/address"
)

var (
	testAddr  *address.Address
	testPub   ed25519.PublicKey
	testRegex *regexp.Regexp = regexp.MustCompile("::")
)

func TestAddrForKey(t *testing.T) {
	for i := 100000; i > 0; i-- {
		pub, _ := GenerateKeyEd25519()
		origIP := net.IP(address.AddrForKey(pub)[:])
		modIP := net.IP(AddrForKey(pub)[:])
		if !origIP.Equal(modIP) {
			t.Errorf("got %s, expected %s", modIP, origIP)
		}
	}
}

func TestMain(m *testing.M) {
	testPub, _ = GenerateKeyEd25519()
	os.Exit(m.Run())
}

func BenchmarkOrigAddrKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testAddr = address.AddrForKey(testPub)
	}
}

func BenchmarkModdedAddrForKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testAddr = AddrForKey(testPub)
	}
}

// measures overall performance of code from cmd/genkeys
func BenchmarkOrigLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pub, _ := GenerateKeyEd25519()
		ip := net.IP(address.AddrForKey(pub)[:]).String()
		testRegex.MatchString(ip)
	}
}

// measures overall performance of functions we vendor
func BenchmarkModdedLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pub, _ := GenerateKey()
		ip := net.IP(AddrForKey(pub)[:]).String()
		testRegex.MatchString(ip)
	}
}
