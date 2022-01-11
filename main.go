package main

import (
	"autok3s/auto-k3s"
	"flag"
	"fmt"
	"net"
)

var baseNet string
var nodeQuantity int

func getIpConfig() ([]string, string) {
	fmt.Println("Looking up my IPs by interface...")
	myIps, err := autok3s.GetMyIp()
	autok3s.CheckErr(err)

	fmt.Println("Determining 4th octet.")
	localFourthOctet, err := autok3s.GetFourthOctet(myIps)
	autok3s.CheckErr(err)

	return myIps, localFourthOctet
}

func parseFlags() (e autok3s.EnvConfig) {
	flag.StringVar(&baseNet, "baseNet", "none", "The base network for the pi")
	flag.IntVar(&nodeQuantity, "nodeQuantity", 1, "The number of nodes in the cluster")
	flag.Parse()

	baseNet = autok3s.ParseBaseNet(baseNet)
	fmt.Println("Base network is: ", baseNet)

	nodeQuantity = autok3s.ParseNodeQuantity(nodeQuantity)
	fmt.Println("node Quantity is: ", nodeQuantity)

	envConfig := autok3s.EnvConfig{
		NodeQuantity: nodeQuantity,
		BaseNet:      baseNet,
	}

	return envConfig
}

func main() {
	// instantiate environment configuration based on flag values
	envConfig := parseFlags()

	//set up deployment matrix
	deploymentMatrix := autok3s.DeploymentMatrix{
		DeploymentComplete: false,
		DeployNodeReady:    false,
		DeployMasterReady:  true,
	}

	channel1 := make(chan string)
	channel2 := make(chan string)

	//Set up IP Config
	myIps, localFourthOctet := getIpConfig()
	ipconfig := autok3s.IpConfig{
		MyIps:            myIps,
		LocalFourthOctet: localFourthOctet,
	}

	fmt.Println("Setting UDP server", localFourthOctet)
	pc, err := net.ListenPacket("udp4", ":8830")
	autok3s.CheckErr(err)

	autok3s.Run(deploymentMatrix, pc, channel1, channel2, ipconfig, envConfig)
	//defer pc.Close()
}
