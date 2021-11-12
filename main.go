package main

import (
	k3s_deploy "autok3s/k3s-deploy"
	"autok3s/udp"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var myOctet int
var otherOctet int

type packet struct{}

func main() {

	c := make(chan string)
	c2 := make(chan string)

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

	alreadyDeployed := false
	deployNodeReady := false
	deployMasterReady := true
	for {
		wg.Add(1)
		fmt.Println("Launching Goroutine for udp server...")

		switch deployNodeReady {
		case false:
			go func() {
				udp.Receive(pc, myIps, c, &wg)
			}()
		case true:
			go func() {
				udp.ReceiveToken(pc, myIps, c2, &wg)
			}()
		}

		fmt.Println("Sending a packet")
		udp.Send(pc, octet)
		fmt.Println("Waiting for receiving to finish...")
		select {
		case channelReceive := <-c:
			fmt.Println("Received IP", channelReceive)

			myOctet, err = strconv.Atoi(octet)
			if err != nil {
				log.Fatal(err)
			}

			otherOctet, err = strconv.Atoi(channelReceive)
			if err != nil {
				log.Fatal(err)
			}

			if myOctet < otherOctet {
				os.Setenv("K3S_MASTER", "true")
				fmt.Println("Set K3s Master to TRUE")
				if deployMasterReady == true {
					nodeToken, _ := k3s_deploy.DeployMaster()
					fmt.Println(nodeToken)
					udp.Send(pc, string(nodeToken[:]))
					deployMasterReady = false
				}
			} else {
				os.Setenv("K3S_MASTER", "false")
				fmt.Println("Set K3s Master to FALSE")
				if alreadyDeployed == false {
					deployNodeReady = true
				}
			}

		case channelReceive := <-c2:
			fmt.Println(otherOctet)
			o := strconv.Itoa(otherOctet)
			fmt.Printf("%s:%T", o, o)
			k3s_deploy.DeployNode(channelReceive, o)
			deployNodeReady = false
			alreadyDeployed = true
		}

		wg.Wait()

		fmt.Println("K3S_MASTER: ", os.Getenv("K3S_MASTER"))
	}
	//defer pc.Close()
}
