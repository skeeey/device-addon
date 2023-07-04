package mqtt

type MQTTBrokerInfo struct {
	Host      string `json:"host"`
	ClientId  string `json:"clientId"`
	Qos       int    `json:"qos"`
	KeepAlive int    `json:"keepAlive"`

	ConnEstablishingRetry int `json:"connEstablishingRetry"`

	AuthMode      string `json:"authMode"`
	CredentialDir string `json:"credentialDir"`

	SubTopic string `json:"subTopic"`
	PubTopic string `json:"pubTopic"`
}
