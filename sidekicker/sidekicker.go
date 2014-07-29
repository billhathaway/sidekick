package main

import (
	"flag"
	"fmt"
	"github.com/billhathaway/sidekick"
	"os"
)

func main() {
	etcdServers := flag.String("servers", "http://localhost:4001", "etcdServers in CSV list")
	key := flag.String("key", "", "key in etcd")
	value := flag.String("val", "", "value for key")
	flag.Parse()
	if *key == "" {
		fmt.Println("Error: key must be set")
		flag.Usage()
		os.Exit(1)
	}

	_, err := sidekick.New(*etcdServers, *key, *value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	select {}
}
