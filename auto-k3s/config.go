package autok3s

import (
	"os"
)

const (
	nodeTokenPath    = "/var/lib/rancher/k3s/server/node-token"
)

//example BASE_NET=192.168.86
var base_net string     = os.Getenv("BASE_NET")

//example NODE_QUANTITY=4
var nodeQuantity = os.Getenv("NODE_QUANTITY")
