package main

import (
	"autok3s/auto-k3s"
	"fmt"
	"net"
)

func main() {

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

	autok3s.Run(deploymentMatrix, pc, c, c2, ipconfig)
	//defer pc.Close()
}
