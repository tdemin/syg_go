package main

/*
This program is largely based on cmd/genkeys from yggdrasil-go, located at
https://github.com/yggdrasil-network/yggdrasil-go

See cmd/genkeys/main.go@78b5f88e4bb734d0dd6a138ff08d34ca39dcaea3
*/

import (
	"encoding/hex"
	"flag"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"

	"github.com/yggdrasil-network/yggdrasil-go/src/address"
	"github.com/yggdrasil-network/yggdrasil-go/src/crypto"
)

var version = "v0.1.3"

func main() {
	rxflag := flag.String("regex", "::", "regex to match addresses against")
	threads := flag.Int("threads", runtime.GOMAXPROCS(0), "how many threads to use for mining")
	iterationsPerOutput := flag.Uint("iter", 100000, "per how many iterations to output status")
	displayVersion := flag.Bool("version", false, "display version")
	origCode := flag.Bool("original", false, "use original Yggdrasil code")
	highAddressMode := flag.Bool("highaddr", false, "high address mining mode, excludes regex")
	flag.Parse()
	if *displayVersion {
		println("syg_go", version)
		return
	}

	if *origCode {
		log.Println("using unmodified Yggdrasil code")
		addrForNodeID = address.AddrForNodeID
	} else {
		log.Println("using syg_go vendored code")
		addrForNodeID = AddrForNodeID
	}

	regex, err := regexp.Compile(*rxflag)
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	newKeys := make(chan keySet, *threads)
	var currentBest []byte

	if !*highAddressMode {
		log.Printf("starting mining for %v with %v threads\n", regex, *threads)
	} else {
		log.Printf("starting mining higher addresses with %v threads\n", *threads)
	}
	for i := 0; i < *threads; i++ {
		go doBoxKeys(newKeys)
	}

	counter := uint64(0)
	i := uint64(*iterationsPerOutput)
	if !*highAddressMode {
		for {
			newKey := <-newKeys
			if regex.MatchString(newKey.ip) {
				newKey.print()
			}
			counter++
			if counter%i == 0 {
				log.Printf("reached %v iterations\n", counter)
			}
		}
	} else {
		for {
			newKey := <-newKeys
			if isBetter(currentBest[:], newKey.id) || len(currentBest) == 0 {
				currentBest = newKey.id
				newKey.print()
			}
			counter++
			if counter%i == 0 {
				log.Printf("reached %v iterations\n", counter)
			}
		}
	}
}

type keySet struct {
	priv []byte
	pub  []byte
	id   []byte
	ip   string
}

func (k *keySet) print() {
	log.Printf("priv: %s | pub: %s | nodeid: %s | ip: %s\n",
		hex.EncodeToString(k.priv[:]),
		hex.EncodeToString(k.pub[:]),
		hex.EncodeToString(k.id[:]),
		k.ip)
}

func doBoxKeys(out chan<- keySet) {
	for {
		pub, priv := crypto.NewBoxKeys()
		id := crypto.GetNodeID(pub)
		ip := net.IP(addrForNodeID(id)[:]).String()
		out <- keySet{priv[:], pub[:], id[:], ip}
	}
}

func isBetter(oldID, newID []byte) bool {
	for i := range oldID {
		if newID[i] > oldID[i] {
			return true
		}
		if newID[i] < oldID[i] {
			return false
		}
	}
	return false
}
