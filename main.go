package main

import (
	"autok3s/udp"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

type packet struct{}

func main() {

	c := make(chan string)

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

	go func() {
		udp.Receive(pc, myIps, c, &wg)
	}()

	fmt.Println("Sending a packet")
	udp.Send(pc, octet)
	fmt.Println("Waiting for receiving to finish...")
	select {
	case otherIP := <-c:
		fmt.Println("Received", otherIP)
		myOctet, err := strconv.Atoi(octet)
		if err != nil {
			log.Fatal(err)
		}

		otherOctet, err := strconv.Atoi(otherIP)
		if err != nil {
			log.Fatal(err)
		}

		if myOctet > otherOctet {
			os.Setenv("K3S_MASTER", "true")
		} else {
			os.Setenv("K3S_MASTER", "false")
		}
	}
	wg.Wait()

	fmt.Println("K3S_MASTER: ", os.Getenv("K3S_MASTER"))

	defer pc.Close()
}
