package main

import (
	"autok3s/auto-k3s"
	"flag"
	"fmt"
	"net"
)

var baseNet string
var nodeQuantity int

func main() {

	flag.StringVar(&baseNet, "baseNet", "none", "The base network for the pi")
	flag.IntVar(&nodeQuantity, "nodeQuantity", 1, "The number of nodes in the cluster")
	flag.Parse()

	baseNet = autok3s.ParseBaseNet(baseNet)
	fmt.Println("Base network is: ", baseNet)
	nodeQuantity = autok3s.ParseNodeQuantity(nodeQuantity)
	fmt.Println("node Quantity is: ", nodeQuantity)

	c := make(chan string)
	c2 := make(chan string)

	deploymentMatrix := autok3s.DeploymentMatrix{
		AlreadyDeployed:   false,
		DeployNodeReady:   false,
		DeployMasterReady: true,
	}

	fmt.Println("Looking up my IPs by interface...")

	myIps, err := autok3s.GetMyIp()
	autok3s.CheckErr(err)
	fmt.Println("Found them.\n Determining 4th octet.")
	localFourthOctet, err := autok3s.GetFourthOctet(myIps)
	autok3s.CheckErr(err)
	fmt.Println("Done - setting up connection", localFourthOctet)

	ipconfig := autok3s.IpConfig{
		MyIps:            myIps,
		LocalFourthOctet: localFourthOctet,
	}

	pc, err := net.ListenPacket("udp4", ":8830")

	autok3s.CheckErr(err)

	autok3s.Run(deploymentMatrix, pc, c, c2, ipconfig, baseNet, nodeQuantity)
	//defer pc.Close()
}
