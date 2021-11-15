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
var wg sync.WaitGroup

func Run(d DeploymentMatrix, pc net.PacketConn, c chan string, c2 chan string, ipconfig IpConfig) {
	for {
		wg.Add(1)
		fmt.Println("Launching Goroutine for udp server...")
		switch d.DeployNodeReady {
		case false:
			fmt.Println("Receiving an IP")
			go func() {
				Receive(pc, ipconfig.MyIps, c, &wg)
			}()
		case true:
			fmt.Println("Receiving a token")
			go func() {
				ReceiveToken(pc, ipconfig.MyIps, c2, &wg)
			}()
		}

		fmt.Println("Sending a packet")
		Send(pc, ipconfig.LocalFourthOctet)

		fmt.Println("Waiting for receiving to finish...")
		select {
		case channelReceive := <-c:
			fmt.Println("Received IP", channelReceive)

			intLocalFourthOctet, err := strconv.Atoi(ipconfig.LocalFourthOctet)
			CheckErr(err)
			intRemoteFourthOctet, err = strconv.Atoi(channelReceive)
			CheckErr(err)

			if intLocalFourthOctet < intRemoteFourthOctet {
				err = os.Setenv("K3S_MASTER", "true")
				CheckErr(err)

				fmt.Println("Set K3s Master to TRUE")
				if d.DeployMasterReady == true {
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
