package mqtt

type MQTTBrokerInfo struct {
	Host      string `yaml:"host"`
	Qos       int    `yaml:"qos"`
	KeepAlive int    `yaml:"keepAlive"`
	ClientId  string `yaml:"clientId"`

	ConnEstablishingRetry int `yaml:"connEstablishingRetry"`

	AuthMode      string `yaml:"authMode"`
	CredentialDir string `yaml:"credentialDir"`

	SubTopic string `yaml:"subTopic"`
	PubTopic string `yaml:"pubTopic"`
}
