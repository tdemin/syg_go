// +build original

package main

import (
	"log"

	"github.com/yggdrasil-network/yggdrasil-go/src/address"
)

func init() {
	log.Println("using unmodified Yggdrasil code")
	addrForNodeID = address.AddrForNodeID
}
