package main

/*
Parts of this program are taken from yggdrasil-go, located at
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

var version = "v0.1.1"

func main() {
	rxflag := flag.String("regex", "::", "regex to match addresses against")
	threads := flag.Int("threads", runtime.GOMAXPROCS(0), "how many threads to use for mining")
	iterationsPerOutput := flag.Uint("iter", 100000, "per how many iterations to output status")
	displayVersion := flag.Bool("version", false, "display version")
	origCode := flag.Bool("original", false, "use original Yggdrasil code")
	highAddressMode := flag.Bool("highaddr", false, "high address mining mode (2xx::), excludes regex")
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

	var (
		regex *regexp.Regexp
		err   error
	)
	if !*highAddressMode {
		regex, err = regexp.Compile(*rxflag)
		if err != nil {
			log.Printf("%v\n", err)
			os.Exit(1)
		}
	} else {
		regex = regexp.MustCompile("^2..::$")
	}

	newKeys := make(chan keySet, *threads)
	log.Printf("starting mining for %v with %v threads\n", regex, *threads)
	for i := 0; i < *threads; i++ {
		go doBoxKeys(newKeys)
	}

	counter := uint64(0)
	i := uint64(*iterationsPerOutput)
	for {
		newKey := <-newKeys
		if regex.MatchString(newKey.ip) {
			log.Printf("priv: %s | pub: %s | nodeid: %s | ip: %s\n",
				hex.EncodeToString(newKey.priv[:]),
				hex.EncodeToString(newKey.pub[:]),
				hex.EncodeToString(newKey.id[:]),
				newKey.ip)
		}
		counter++
		if counter%i == 0 {
			log.Printf("reached %v iterations\n", counter)
		}
	}
}

type keySet struct {
	priv []byte
	pub  []byte
	id   []byte
	ip   string
}

func doBoxKeys(out chan<- keySet) {
	for {
		pub, priv := crypto.NewBoxKeys()
		id := crypto.GetNodeID(pub)
		ip := net.IP(addrForNodeID(id)[:]).String()
		out <- keySet{priv[:], pub[:], id[:], ip}
	}
}
