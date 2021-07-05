package main

/*
This program is largely based on cmd/genkeys from yggdrasil-go, located at
https://github.com/yggdrasil-network/yggdrasil-go

See cmd/genkeys/main.go@78b5f88e4bb734d0dd6a138ff08d34ca39dcaea3
*/

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"log"
	"net"
	"os"
	"regexp"
	"runtime"

	"github.com/yggdrasil-network/yggdrasil-go/src/address"
)

var (
	programName = "syg_go"
	version     = "v0.1.4"
	copyright   = "Copyright (c) 2020 Timur Demin"
)

var (
	stdout = log.New(os.Stdout, "", log.Flags())
	stderr = log.New(os.Stderr, "", log.Flags())
)

func main() {
	rxflag := flag.String("regex", "::", "regex to match addresses against")
	threads := flag.Int("threads", runtime.GOMAXPROCS(0), "how many threads to use for mining")
	iterationsPerOutput := flag.Uint("iter", 100000, "per how many iterations to output status")
	displayVersion := flag.Bool("version", false, "display version")
	origCode := flag.Bool("original", false, "use original Yggdrasil code")
	highAddressMode := flag.Bool("highaddr", false, "high address mining mode, excludes regex")
	flag.Parse()
	if *displayVersion {
		println(programName, version)
		println(copyright)
		return
	}

	if *origCode {
		stdout.Println("using unmodified Yggdrasil code")
		addrForKey = address.AddrForKey
		generateKey = GenerateKeyEd25519
	} else {
		stdout.Println("using syg_go vendored code")
		addrForKey = AddrForKey
		generateKey = GenerateKey
	}

	regex, err := regexp.Compile(*rxflag)
	if err != nil {
		stderr.Printf("%v\n", err)
		os.Exit(1)
	}

	newKeys := make(chan keySet, *threads)
	var currentBest = make(ed25519.PublicKey, ed25519.PublicKeySize)
	for i := range currentBest {
		currentBest[i] = 0xff
	}

	for i := 0; i < *threads; i++ {
		go doBoxKeys(newKeys)
	}

	counter := uint64(0)
	i := uint64(*iterationsPerOutput)
	if !*highAddressMode {
		stdout.Printf("starting mining for %v with %v threads\n", regex, *threads)
		for {
			newKey := <-newKeys
			if regex.MatchString(newKey.ip) {
				newKey.print()
			}
			counter++
			if counter%i == 0 {
				stderr.Printf("reached %v iterations\n", counter)
			}
		}
	} else {
		stdout.Printf("starting mining higher addresses with %v threads\n", *threads)
		for {
			newKey := <-newKeys
			if isBetter(currentBest, newKey.pub) {
				currentBest = newKey.pub
				newKey.print()
			}
			counter++
			if counter%i == 0 {
				stderr.Printf("reached %v iterations\n", counter)
			}
		}
	}
}

type keySet struct {
	priv []byte
	pub  []byte
	ip   string
}

func (k *keySet) print() {
	stdout.Printf("priv: %s | pub: %s | ip: %s\n",
		hex.EncodeToString(k.priv[:]),
		hex.EncodeToString(k.pub[:]),
		k.ip)
}

func doBoxKeys(out chan<- keySet) {
	for {
		pub, priv := generateKey()
		ip := net.IP(addrForKey(pub)[:]).String()
		out <- keySet{priv[:], pub[:], ip}
	}
}

func isBetter(oldPub, newPub ed25519.PublicKey) bool {
	for i := range oldPub {
		if newPub[i] < oldPub[i] {
			return true
		}
		if newPub[i] > oldPub[i] {
			break
		}
	}
	return false
}
