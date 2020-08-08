// +build !original

package main

import "log"

func init() {
	log.Println("using vendored syg_go code")
	addrForNodeID = AddrForNodeID
}
