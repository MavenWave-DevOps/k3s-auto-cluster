package autok3s

type DeploymentMatrix struct {
	AlreadyDeployed bool
	DeployNodeReady bool
	DeployMasterReady bool
}

type IpConfig struct {
	MyIps []string
	LocalFourthOctet string
}

type AutoClusterConfig struct {
	
}
