package main

import "github.com/hashicorp/consul/api"
import "fmt"

func main() {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	// PUT a new KV pair
	p := &api.KVPair{Key: "aamesmsmsmss", Value: []byte("hlhlhlh")}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	}

	// Lookup the pair
	pair, _, err := kv.Get("aamesmsmsmss", nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)
}
