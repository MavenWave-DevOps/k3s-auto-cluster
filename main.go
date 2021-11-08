package main

import (
	"autok3s/udp"
	"fmt"
	"net"
	"sync"
)

var wg sync.WaitGroup

type packet struct{}

func main() {
	fmt.Println("Starting main function")

	fmt.Println("Looking up my IPs by interface...")
	myIps, err := udp.GetMyIp()
	if err != nil {
		panic(err)
	}
	fmt.Println("Found them.\n Setting up connection.")
	pc, err := net.ListenPacket("udp4", ":8830")
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	fmt.Println("Launching Goroutine for udp server...")

	go func(){
		defer wg.Done()
		udp.Receive(pc, myIps)
	}()

	fmt.Println("Sending a packet")
	udp.Send(pc)
	fmt.Println("Waiting for receiving to finish...")
	wg.Wait()
	defer pc.Close()
}
