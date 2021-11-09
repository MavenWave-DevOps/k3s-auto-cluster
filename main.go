package main

import (
	"autok3s/udp"
	"fmt"
	"log"
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
	fmt.Println("Found them.\n Determining 4th octet.")
	octet, err := udp.GetFourthOctet(myIps)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done - setting up connection", octet)
	pc, err := net.ListenPacket("udp4", ":8830")
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	fmt.Println("Launching Goroutine for udp server...")

	go func() string {
		defer wg.Done()
		r := udp.Receive(pc, myIps)
		return r
	}()

	fmt.Println("Sending a packet")
	udp.Send(pc, octet)
	fmt.Println("Waiting for receiving to finish...")
	wg.Wait()
	defer pc.Close()
}
