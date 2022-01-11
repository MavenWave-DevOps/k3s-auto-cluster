package autok3s

import "net"

type DeploymentMatrix struct {
	DeploymentComplete bool
	DeployNodeReady    bool
	DeployMasterReady  bool
}

type EnvConfig struct {
	NodeQuantity int
	BaseNet      string
}

type IpConfig struct {
	MyIps            []string
	LocalFourthOctet string
}

type PiConfig struct {
	PiEnvConfig      EnvConfig
	PiIpConfig       IpConfig
	DeploymentMatrix DeploymentMatrix
	Pc               net.PacketConn
}

type Packet struct {
	Conn    net.PacketConn
	Payload string
}
