package mqtt

import "github.com/skeeey/device-addon/pkg/device/client"

type Config struct {
	client.MQTTBrokerInfo `json:"inline"`

	SubTopic string `json:"subTopic"`
	PubTopic string `json:"pubTopic"`
}
