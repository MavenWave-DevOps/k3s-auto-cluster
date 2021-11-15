package autok3s

import "net"

type DeploymentMatrix struct {
	AlreadyDeployed   bool
	DeployNodeReady   bool
	DeployMasterReady bool
}

type IpConfig struct {
	MyIps            []string
	LocalFourthOctet string
}

type AutoClusterConfig struct {
}

type Packet struct {
	Conn    net.PacketConn
	Payload string
}
