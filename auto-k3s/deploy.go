package autok3s

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var wg2 sync.WaitGroup

func commonDeploy(cmd *exec.Cmd) string {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	fmt.Println("Running: ", cmd.String())
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	CheckErr(err)

	return fmt.Sprintf("%s\n", stdout.String())
}

func (p PiConfig) DeployMaster() (string, error) {

	script := []byte("#!/bin/bash\nset -e\ncurl -s -f -L https://get.k3s.io | sh - 2> errors\n exit 0")
	err := os.WriteFile("install.sh", script, 0777)
	CheckErr(err)

	cmd := exec.Command("/bin/bash", "install.sh")
	commonDeploy(cmd)

	fmt.Println("Finished Install Script: ")
	if os.Getenv("K3S_MASTER") == "true" {

		time.Sleep(100 * time.Millisecond)

		cmd2 := exec.Command("sudo", "cat", nodeTokenPath)
		o := commonDeploy(cmd2)
		cmd3 := exec.Command("sudo", "chmod", "777", "/etc/rancher/k3s/k3s.yaml")
		oo := commonDeploy(cmd3)

		fmt.Println("Retrieved node token...")
		//output, err = cmd2.Output()
		return fmt.Sprintf("%s\n%s", o, oo), nil
	}
	return "done", nil
}

func (p PiConfig) DeployNode(token string, m string) {

	fmt.Println("Starting node deployment...")

	s := []string{p.PiEnvConfig.BaseNet, m}
	masterIp := strings.Join(s, ".")
	masterURL := fmt.Sprintf("https://%s:6443", masterIp)

	fmt.Println("Setting URL")
	err := os.Setenv("K3S_URL", masterURL)
	CheckErr(err)

	fmt.Println("Set URL. setting token")
	err = os.Setenv("K3S_TOKEN", token)
	CheckErr(err)
	fmt.Println("Set token")

	_, err = p.DeployMaster()
	CheckErr(err)

	time.Sleep(100 * time.Millisecond)

	script := []byte(fmt.Sprintf("#!/bin/bash\nset -e\nsudo nohup k3s agent --server %s --token %s &\n exit 0", os.Getenv("K3S_URL"), os.Getenv("K3S_TOKEN")))
	err = os.WriteFile("agent.sh", script, 0777)
	CheckErr(err)

	wg2.Add(1)

	go func() {
		fmt.Println("Starting the K3S agent in a separate process...")
		cmd := exec.Command("/bin/bash", "agent.sh")
		wg2.Done()
		commonDeploy(cmd)
	}()
	wg2.Wait()
	//Adding a sleep to give the agent a little time to fully start
	time.Sleep(30 * time.Second)
	fmt.Println("done with agent setup")
	//output, err = cmd2.Output()
	return
}
