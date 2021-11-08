package udp

import (
	"fmt"
	"net"
	"strings"
	"time"
)
var myifaces []string

func getFourthOctet() (string, error) {
	return "test", nil
}

func GetMyIp() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces{
		if iface.Name == "eth0"{
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				myifaces = append(myifaces, addr.String())
			}
		}
	}
	return myifaces, nil
}

func Send(pc net.PacketConn) {
	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:8830")
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)

	for i := 0; i < 10; i++ {
		_, err = pc.WriteTo([]byte("this is a test"), addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Done sending packet %d - Packet sent successfully\n", i)
		time.Sleep(5*time.Second)
	}
}
func Receive(pc net.PacketConn, myIps []string) {
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
			break
		}

		continue
	}
}