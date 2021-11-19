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

func CheckMaster(NodeIPs []string, ipconfig IpConfig) bool {
	for _, StrRemoteIp := range NodeIPs {
		intLocalFourthOctet, err := strconv.Atoi(ipconfig.LocalFourthOctet)
		CheckErr(err)
		intRemoteFourthOctet, err = strconv.Atoi(StrRemoteIp)
		CheckErr(err)
		if intLocalFourthOctet > intRemoteFourthOctet {
			return false
		}
	}
	return true
}

func Run(d DeploymentMatrix, pc net.PacketConn, c chan string, c2 chan string, ipconfig IpConfig) {
	for {

		wg.Add(1)

		switch d.DeployNodeReady {
		case false:
			fmt.Println("UDP Server - Receiving an IP")
			go func() {
				Receive(pc, ipconfig.MyIps, c, c2, &wg)
			}()
		case true:
			fmt.Println("UDP Server - Receiving a token")
			go func() {
				ReceiveToken(pc, ipconfig.MyIps, c2, &wg)
			}()
		}

		fmt.Println("Sending a packet...")
		Send(pc, ipconfig.LocalFourthOctet)

		fmt.Println("Waiting for receiving to finish...")
		select {
		case channelReceive := <-c:
			fmt.Printf("Received an IP: %s", channelReceive)
			master := false

			if len(NodeIPs) > 0 {
				for _, StrRemoteIp := range NodeIPs {
					if channelReceive != StrRemoteIp {
						NodeIPs = append(NodeIPs, channelReceive)
					}
				}
			} else {
				NodeIPs = append(NodeIPs, channelReceive)
			}

			if len(NodeIPs) == nodeQuantity-1 {
				fmt.Printf("Node IPs and node quantity match\n Node IPs are: ", NodeIPs)
				// Make updates here to wait for all IPs to come in
				master = CheckMaster(NodeIPs, ipconfig)

				if master == true {
					err := os.Setenv("K3S_MASTER", "true")
					CheckErr(err)

					fmt.Println("Set K3s Master to TRUE")
					if d.DeployMasterReady == true {
						Send(pc, ipconfig.LocalFourthOctet)
						nodeToken, _ := DeployMaster()
						fmt.Println(nodeToken)
						Send(pc, string(nodeToken[:]))
						d.DeployMasterReady = false
					}
				} else {
					os.Setenv("K3S_MASTER", "false")
					fmt.Println("Set K3s Master to FALSE")
					if d.AlreadyDeployed == false {
						d.DeployNodeReady = true
					}
				}
			}

		case channelReceive := <-c2:
			fmt.Println(intRemoteFourthOctet)
			stringRemoteFourthOctet := strconv.Itoa(intRemoteFourthOctet)
			fmt.Printf("%s:%T", stringRemoteFourthOctet, stringRemoteFourthOctet)
			DeployNode(channelReceive, stringRemoteFourthOctet)
			d.DeployNodeReady = false
			d.AlreadyDeployed = true
		}

		wg.Wait()

		fmt.Println("K3S_MASTER: ", os.Getenv("K3S_MASTER"))
	}
}
