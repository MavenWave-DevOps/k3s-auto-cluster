package autok3s

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

func (p PiConfig) Send(payload string) {
	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:8830")
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)

	for i := 0; i < 1; i++ {
		_, err = p.Pc.WriteTo([]byte(payload), addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Done sending packet %d - Packet sent successfully\n", i)
		time.Sleep(5 * time.Second)
	}
}

func (p PiConfig) ReceiveToken(c2 chan string, wg *sync.WaitGroup) {
	for {
		buf := make([]byte, 1024)
		n, addr, err := p.Pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		if len(buf[:n]) < 25 {
			continue
		}

		loopOn := false
		for _, ip := range p.PiIpConfig.MyIps {
			if strings.Split(addr.String(), ":")[0] == strings.Split(ip, "/")[0] {
				fmt.Printf("Lol, this packet is from me - payload: %s\n", buf[:n])
				loopOn = true
				break
			}
		}
		if loopOn == false {
			fmt.Printf("%s sent this: %s\n", addr, buf[:n])
			c2 <- string(buf[:n])
			wg.Done()
			return
		}
		continue
	}
}

func (p PiConfig) Receive(c chan string, wg *sync.WaitGroup) {
	for {
		buf := make([]byte, 1024)
		n, addr, err := p.Pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		loopOn := false
		fmt.Println("length of addr ", len(buf[:n]))
		for _, ip := range p.PiIpConfig.MyIps {
			if strings.Split(addr.String(), ":")[0] == strings.Split(ip, "/")[0] {
				fmt.Printf("Lol, this packet is from me - payload: %s\n", buf[:n])
				loopOn = true
				break
			}
		}

		if loopOn == false {
			fmt.Printf("%s sent this: %s\n", addr, buf[:n])
			c <- string(buf[:n])
			wg.Done()
			return
		}
		continue
	}
}
