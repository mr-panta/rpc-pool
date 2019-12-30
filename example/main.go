package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/mr-panta/rpc-pool"
)

func main() {
	startServer()
	startClient()
}

type Echo struct{}

func (e *Echo) Do(input *int, output *int) error {
	*output = *input
	return nil
}

func startServer() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	e := &Echo{}
	s := rpc.NewServer()
	err = s.Register(e)
	if err != nil {
		log.Fatal(err)
	}
	go s.Accept(l)
}

func startClient() {
	client, err := rpc_pool.NewRPCPool(":8888", 0, 1, time.Second, time.Second, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	n := 100
	wg := &sync.WaitGroup{}
	for n > 0 {
		wg.Add(1)
		go call(client, n, wg)
		n--
	}
	wg.Wait()
}

func call(client rpc_pool.RPCPool, n int, wg *sync.WaitGroup) {
	input := n
	output := 0
	err := client.Call("Echo.Do", &input, &output)
	if err != nil {
		log.Println(err)
	}
	log.Println(output)
	wg.Done()
}
