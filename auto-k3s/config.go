package autok3s

import (
	"strings"
)

//example BASE_NET=192.168.86
//example NODE_QUANTITY=4

const (
	nodeTokenPath = "/var/lib/rancher/k3s/server/node-token"
)

func ParseBaseNet(baseNet string) string {
	splitNet := strings.Split(baseNet, ".")
	splitNet = splitNet[:len(splitNet)-1]
	return strings.Join(splitNet, ".")
}

func ParseNodeQuantity(nodeQuantity int) int {
	//TODO - add some checks here
	return nodeQuantity
}
