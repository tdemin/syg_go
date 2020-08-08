package main

import (
	"net"
	"os"
	"regexp"
	"testing"

	"github.com/yggdrasil-network/yggdrasil-go/src/address"
	"github.com/yggdrasil-network/yggdrasil-go/src/crypto"
)

var (
	testAddr   *address.Address
	testPub    *crypto.BoxPubKey
	testNodeID *crypto.NodeID
	testRegex  *regexp.Regexp
)

func TestAddrForNodeID(t *testing.T) {
	for i := 20; i > 0; i-- {
		pub, _ := crypto.NewBoxKeys()
		id := crypto.GetNodeID(pub)
		origIP := net.IP(address.AddrForNodeID(id)[:])
		modIP := net.IP(AddrForNodeID(id)[:])
		if !origIP.Equal(modIP) {
			t.Errorf("got %s, expected %s", modIP, origIP)
		}
	}
}

func TestMain(m *testing.M) {
	testPub, _ = crypto.NewBoxKeys()
	testNodeID = crypto.GetNodeID(testPub)
	testRegex = regexp.MustCompile("::")
	os.Exit(m.Run())
}

func BenchmarkOrigAddrForNodeID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testAddr = address.AddrForNodeID(testNodeID)
	}
}

func BenchmarkModdedAddrForNodeID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testAddr = AddrForNodeID(testNodeID)
	}
}

// measures overall performance of code from cmd/genkeys
func BenchmarkOrigLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pub, _ := crypto.NewBoxKeys()
		id := crypto.GetNodeID(pub)
		ip := net.IP(address.AddrForNodeID(id)[:]).String()
		testRegex.MatchString(ip)
	}
}

// measures overall performance of functions we vendor
func BenchmarkModdedLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pub, _ := crypto.NewBoxKeys()
		id := crypto.GetNodeID(pub)
		ip := net.IP(AddrForNodeID(id)[:]).String()
		testRegex.MatchString(ip)
	}
}
