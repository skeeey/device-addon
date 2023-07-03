package opcua

const (
	Protocol = "opcua"
	Endpoint = "endpoint"
)

const NODE = "nodeId"

type OPCUAServerInfo struct {
	// Security policy: None, Basic128Rsa15, Basic256, Basic256Sha256
	SecurityPolicy string `yaml:"securityPolicy"`
	// Security mode: None, Sign, SignAndEncrypt
	SecurityMode string `yaml:"securityMode"`
	CertFile     string `yaml:"certFile"`
	KeyFile      string `yaml:"keyFile"`
}
