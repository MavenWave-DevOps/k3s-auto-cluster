package autok3s

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

var intLocalFourthOctet int
var intRemoteFourthOctet int
var NodeIPs []string
var wg sync.WaitGroup

//ex:
//my ip: 21
//nodeIPs: [22, 20]
// CheckMaster loops through node ips and checks whether the nodes IP is the lowest out of all the nodes

func (p PiConfig) CheckMaster(NodeIPs []string) bool {
	for _, StrRemoteIp := range NodeIPs {
		intLocalFourthOctet, err := strconv.Atoi(p.PiIpConfig.LocalFourthOctet)
		CheckErr(err)
		intRemoteFourthOctet, err = strconv.Atoi(StrRemoteIp)
		CheckErr(err)
		if intLocalFourthOctet > intRemoteFourthOctet {
			return false
		}
	}
	return true
}

func Run(d DeploymentMatrix, pc net.PacketConn, c chan string, c2 chan string, ipconfig IpConfig, envConfig EnvConfig) {

	//Set up piconfig
	PiConfig := PiConfig{
		PiEnvConfig:      envConfig,
		PiIpConfig:       ipconfig,
		DeploymentMatrix: d,
		Pc:               pc,
	}

	for {
		wg.Add(1)

		switch PiConfig.DeploymentMatrix.DeployNodeReady {
		case false:
			fmt.Println("UDP Server - Receiving an IP")
			go func() {
				PiConfig.Receive(c, &wg)
			}()
		case true:
			fmt.Println("UDP Server - Receiving a token")
			go func() {
				PiConfig.ReceiveToken(c2, &wg)
			}()
		}

		fmt.Println("Sending a packet...")
		PiConfig.Send(PiConfig.PiIpConfig.LocalFourthOctet)

		fmt.Println("Waiting for receiving to finish...")
		select {
		case channelReceive := <-c:
			fmt.Printf("Received an IP: %s", channelReceive)
			master := false

			if len(NodeIPs) > 0 {
				var Nodeappend = false

				for _, StrRemoteIp := range NodeIPs {
					if channelReceive == StrRemoteIp {
						Nodeappend = false
						break
					}
					Nodeappend = true
				}
				if Nodeappend == true {
					NodeIPs = append(NodeIPs, channelReceive)
				}
			} else {
				NodeIPs = append(NodeIPs, channelReceive)
			}

			if len(NodeIPs) == PiConfig.PiEnvConfig.NodeQuantity-1 {
				fmt.Printf("Node IPs and node quantity match\n Node IPs are: ", NodeIPs)
				// Make updates here to wait for all IPs to come in
				master = PiConfig.CheckMaster(NodeIPs)

				if master == true {
					err := os.Setenv("K3S_MASTER", "true")
					CheckErr(err)

					fmt.Println("Set K3s Master to TRUE")
					if PiConfig.DeploymentMatrix.DeployMasterReady == true {
						PiConfig.Send(PiConfig.PiIpConfig.LocalFourthOctet)
						nodeToken, _ := PiConfig.DeployMaster()
						fmt.Println(nodeToken)
						PiConfig.Send(string(nodeToken[:]))
						PiConfig.DeploymentMatrix.DeployMasterReady = false
					}
				} else {
					os.Setenv("K3S_MASTER", "false")
					fmt.Println("Set K3s Master to FALSE")
					if PiConfig.DeploymentMatrix.DeploymentComplete == false {
						PiConfig.DeploymentMatrix.DeployNodeReady = true
					}
				}
			}

		case channelReceive := <-c2:
			fmt.Println(intRemoteFourthOctet)
			stringRemoteFourthOctet := strconv.Itoa(intRemoteFourthOctet)
			fmt.Printf("%s:%T", stringRemoteFourthOctet, stringRemoteFourthOctet)
			PiConfig.DeployNode(channelReceive, stringRemoteFourthOctet)
			PiConfig.DeploymentMatrix.DeployNodeReady = false
			PiConfig.DeploymentMatrix.DeploymentComplete = true
		}

		wg.Wait()

		fmt.Println("K3S_MASTER: ", os.Getenv("K3S_MASTER"))
	}
}
