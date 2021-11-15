package autok3s

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	base_net = "192.168.80"
	nodeTokenPath = "/var/lib/rancher/k3s/server/node-token"
)

var wg2 sync.WaitGroup

func DeployMaster() (string, error) {

	command := []byte("#!/bin/bash\nset -e\ncurl -s -f -L https://get.k3s.io | sh - 2> errors\n exit 0")
	err := os.WriteFile("install.sh", command, 0777)
	if err != nil {
		log.Fatal(err)
	}
	var stderr bytes.Buffer
	fmt.Println("Set up download command")
	cmd := exec.Command("/usr/bin/bash", "install.sh")
	fmt.Println(cmd.String())
	cmd.Stderr = &stderr
	fmt.Println("executed download command")
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
		log.Fatal(err)
	}
	fmt.Println("done with command 1")
		if os.Getenv("K3S_MASTER") == "true" {
			//fmt.Println(output)
			time.Sleep(100 * time.Millisecond)
			cmd2 := exec.Command("sudo", "cat", nodeTokenPath)
			var out bytes.Buffer
			cmd2.Stdout = &out
			cmd2.Stderr = &stderr
			err = cmd2.Run()
			o := fmt.Sprintf("%s\n", out.String())
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
				log.Fatal(err)
			}
			fmt.Println("done with command 2")
			//output, err = cmd2.Output()
			return o, nil
		}
	return "done", nil
}

func DeployNode(token string, m string) {
	fmt.Println("Starting node deployment")
	fmt.Println("m is: ", m)
	s := []string{base_net, m}
	masterIp := strings.Join(s, ".")
	masterURL := fmt.Sprintf("https://%s:6443", masterIp)
	fmt.Println(masterURL)
	fmt.Println("Setting URL")
	err := os.Setenv("K3S_URL", masterURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Set URL")
	fmt.Println("Setting Token")
	err = os.Setenv("K3S_TOKEN", token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Set token")
	o, err := DeployMaster()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(o)


	time.Sleep(100 * time.Millisecond)



	script := []byte(fmt.Sprintf("#!/bin/bash\nset -e\nsudo k3s agent --server %s --token %s>output.txt 2>&1\n exit 0", os.Getenv("K3S_URL"), os.Getenv("K3S_TOKEN")))
	err = os.WriteFile("agent.sh", script, 0777)
	if err != nil {
		log.Fatal(err)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	wg2.Add(1)

	go func(){
		fmt.Println("Starting the K3S agent in a separate process...")
		cmd := exec.Command("/usr/bin/bash", "agent.sh")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		wg2.Done()
		err = cmd.Run()
		outya := fmt.Sprintf("%s\n", out.String())
		fmt.Println(outya)
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ":" + stderr.String())
			log.Fatal(err)
		}
	}()
	wg2.Wait()
	//Adding a sleep to give the agent a little time to fully start
	time.Sleep(30*time.Second)
	fmt.Println("done with agent setup")
	//output, err = cmd2.Output()
	return
}