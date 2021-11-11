package k3s_deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	base_net = "192.168.80"
)

func DeployMaster() {
	command := strings.Split("curl -sfL https://get.k3s.io | sh -", " ")
	cmd := exec.Command("/bin/bash", "-c", command[0], command[1], command[2], command[3], command[4], command[5])
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	return
}

func DeployNode(m string) {
	s := []string{base_net, m}
	masterIp := strings.Join(s, ".")
	os.Setenv("K3S_URL", fmt.Sprintf("https://%s:6443", masterIp))
	os.Setenv("K3S_TOKEN", "")
	DeployMaster()
	return
}