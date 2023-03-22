package controller

import (
	"context"
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	edgev1alpha1 "github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	deviceclientset "github.com/skeeey/device-addon/pkg/client/clientset/versioned"
	devicesv1alpha1informer "github.com/skeeey/device-addon/pkg/client/informers/externalversions/apis/v1alpha1"
	devicesv1alpha1lister "github.com/skeeey/device-addon/pkg/client/listers/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/utils"
)

type Subscriber struct {
	deviceClient deviceclientset.Interface
	deviceLister devicesv1alpha1lister.DeviceLister
	mqttClient   mqtt.Client
	clusterName  string
	subTopic     string
}

func NewSubscriber(
	deviceClient deviceclientset.Interface,
	deviceInformer devicesv1alpha1informer.DeviceInformer,
	mqttClient mqtt.Client,
	clusterName string,
	subTopic string,
) *Subscriber {
	return &Subscriber{
		deviceClient: deviceClient,
		deviceLister: deviceInformer.Lister(),
		mqttClient:   mqttClient,
		clusterName:  clusterName,
		subTopic:     subTopic,
	}
}

func (s *Subscriber) Run(ctx context.Context) {
	t := s.mqttClient.Subscribe(s.subTopic+"/+", 0, func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		klog.Infof("msg from MQQT topic=%s payload=%s", topic, string(payload))

		name := strings.ReplaceAll(msg.Topic(), s.subTopic+"/", "")

		attrs := map[string]string{}
		if err := json.Unmarshal(payload, &attrs); err != nil {
			klog.Errorf("failed to unmarshal payload %s, %v", string(payload), err)
		}

		// TODO think about how to handle errors
		if err := s.updateDevice(ctx, s.clusterName, name, attrs); err != nil {
			klog.Errorf("failed to update device status, %v", err)
		}
	})
	<-t.Done()
}

func (s *Subscriber) updateDevice(ctx context.Context, clusterName, name string, newAttrs map[string]string) error {
	device, err := s.deviceLister.Devices(s.clusterName).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	modelName := device.Spec.DeviceDataModelRef.Name
	if _, err := s.deviceClient.EdgeV1alpha1().DeviceDataModels().Get(ctx, modelName, metav1.GetOptions{}); err != nil {
		return err
	}

	_, _, err = utils.UpdateDeviceStatus(
		ctx,
		s.deviceClient,
		clusterName,
		name,
		func(oldStatus *edgev1alpha1.DeviceStatus) error {
			oldAttrs := map[string]string{}
			for _, attr := range oldStatus.ReportedAttrs {
				oldAttrs[attr.Name] = attr.Value
			}

			if equality.Semantic.DeepEqual(oldAttrs, newAttrs) {
				return nil
			}

			reportedAttrs := []edgev1alpha1.ReportedAttr{}

			//TODO only update existed attr
			klog.Infof("new attrs: %+v", newAttrs)
			for k, v := range newAttrs {
				reportedAttrs = append(reportedAttrs, edgev1alpha1.ReportedAttr{
					LastUpdatedTime: metav1.Now(),
					Name:            k,
					Value:           v,
				})
			}

			oldStatus.ReportedAttrs = reportedAttrs
			return nil
		},
	)
	return err
}
