package udp

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var myifaces []string



func GetFourthOctet(ips []string) (string, error) {
	fmt.Println(ips)
	fmt.Println(len(ips))
	if len(ips) == 1 {
		return strings.Split(strings.Split(ips[0], "/")[0], ".")[3], nil
	} else if len(ips) == 2 {
		for _, ip := range ips {
			if strings.Contains(ip, ":") {
				continue
			} else {
				return strings.Split(strings.Split(ips[0], "/")[0], ".")[3], nil
			}
		}
	} else {
		return "", errors.New("Length is wrong")
	}

	return "", errors.New("Something weird happened...")
}

func GetMyIp() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Name == "eth0" {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				myifaces = append(myifaces, addr.String())
			}
		}
	}
	return myifaces, nil
}

func Send(pc net.PacketConn, payload string) {
	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:8830")
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)

	for i := 0; i < 1; i++ {
		_, err = pc.WriteTo([]byte(payload), addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Done sending packet %d - Packet sent successfully\n", i)
		time.Sleep(5 * time.Second)
	}
}
func Receive(pc net.PacketConn, myIps []string, c chan string, wg *sync.WaitGroup) {
	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		loop_on := false

		for _, ip := range myIps {
			if strings.Split(addr.String(), ":")[0] == strings.Split(ip, "/")[0] {
				fmt.Printf("Lol, this packet is from me - payload: %s\n", buf[:n])
				loop_on = true
				break
			}
		}
		if loop_on == false {
			fmt.Printf("%s sent this: %s\n", addr, buf[:n])
			c <- string(buf[:n])
			wg.Done()
			return
		}
		continue
	}
}
