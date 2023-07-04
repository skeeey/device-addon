package opcua

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"

	"k8s.io/klog/v2"

	"github.com/skeeey/device-addon/pkg/apis/v1alpha1"
	"github.com/skeeey/device-addon/pkg/device/messagebuses"
	"github.com/skeeey/device-addon/pkg/device/util"
)

type request struct {
	nodeId *ua.NodeID
	handle uint32
	res    v1alpha1.DeviceResource
}

type OPCUADriver struct {
	sync.Mutex
	serverInfo       *OPCUAServerInfo
	msgBuses         []messagebuses.MessageBus
	devices          map[string]v1alpha1.DeviceConfig
	deviceRequests   map[string][]request
	deviceSubCancels map[string]context.CancelFunc
}

func NewOPCUADriver() *OPCUADriver {
	return &OPCUADriver{
		devices:          make(map[string]v1alpha1.DeviceConfig),
		deviceRequests:   make(map[string][]request),
		deviceSubCancels: make(map[string]context.CancelFunc),
	}
}

func (d *OPCUADriver) Initialize(driverConfig util.ConfigProperties, msgBuses []messagebuses.MessageBus) error {
	var serverInfo = &OPCUAServerInfo{}
	if err := util.ToConfigObj(driverConfig, serverInfo); err != nil {
		return err
	}

	d.msgBuses = msgBuses
	d.serverInfo = serverInfo
	return nil
}

func (d *OPCUADriver) GetType() string {
	return "opcua"
}

func (d *OPCUADriver) Start() error {
	//do nothing
	return nil
}

func (d *OPCUADriver) Stop() error {
	return nil
}

func (d *OPCUADriver) AddDevice(device v1alpha1.DeviceConfig) error {
	d.Lock()
	defer d.Unlock()

	_, ok := d.devices[device.Name]
	if !ok {
		go func() {
			if err := d.startSubscription(device); err != nil {
				klog.Errorf("failed to sub device %s, %v", device.Name, err)
			}
		}()

		d.devices[device.Name] = device
	}

	return nil
}

func (d *OPCUADriver) UpdateDevice(device v1alpha1.DeviceConfig) error {
	//TODO
	return nil
}

func (d *OPCUADriver) RemoveDevice(deviceName string) error {
	//TODO
	return nil
}

func (d *OPCUADriver) HandleCommands(deviceName string, command util.Command) error {
	//TODO
	return nil
}

func (d *OPCUADriver) startSubscription(device v1alpha1.DeviceConfig) error {
	_, ok := d.deviceSubCancels[device.Name]
	if ok {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.deviceSubCancels[device.Name] = cancel

	endpoint, err := d.findEndpoint(device)
	if err != nil {
		return err
	}

	endpoints, err := opcua.GetEndpoints(ctx, endpoint)
	if err != nil {
		return err
	}
	ep := opcua.SelectEndpoint(
		endpoints,
		d.serverInfo.SecurityPolicy,
		ua.MessageSecurityModeFromString(d.serverInfo.SecurityMode),
	)
	if ep == nil {
		return fmt.Errorf("failed to find suitable endpoint")
	}

	opts := []opcua.Option{
		opcua.SecurityPolicy(d.serverInfo.SecurityPolicy),
		opcua.SecurityModeString(d.serverInfo.SecurityMode),
		opcua.CertificateFile(d.serverInfo.CertFile),
		opcua.PrivateKeyFile(d.serverInfo.KeyFile),
		opcua.AuthAnonymous(),
		opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous),
	}

	client := opcua.NewClient(ep.EndpointURL, opts...)
	if err := client.Connect(ctx); err != nil {
		return err
	}
	defer client.CloseWithContext(ctx)

	klog.Infof("Connected to opcua server %s", endpoint)

	notifyCh := make(chan *opcua.PublishNotificationData)

	sub, err := client.SubscribeWithContext(ctx, &opcua.SubscriptionParameters{
		Interval: time.Duration(500) * time.Millisecond,
	}, notifyCh)
	if err != nil {
		return err
	}
	defer sub.Cancel(ctx)

	klog.Infof("Created subscription with id %v", sub.SubscriptionID)

	for index, deviceResource := range device.Profile.DeviceResources {
		req, err := d.toRequest(device.Name, index, deviceResource)
		if err != nil {
			return err
		}

		resp, err := sub.Monitor(ua.TimestampsToReturnBoth, valueRequest(req))
		if err != nil || resp.Results[0].StatusCode != ua.StatusOK {
			return err
		}
	}

	// read from subscription's notification channel until ctx is cancelled
	for {
		select {
		case <-ctx.Done():
			return nil
		case res := <-notifyCh:
			if res.Error != nil {
				klog.Errorf("%v", res.Error)
				continue
			}

			switch x := res.Value.(type) {
			case *ua.DataChangeNotification:
				for _, item := range x.MonitoredItems {
					data := item.Value.Value.Value()
					klog.Infof("MonitoredItem with client handle %v = %v", item.ClientHandle, data)

					req := d.findRequest(device.Name, item.ClientHandle)
					if req == nil {
						continue
					}

					result, err := util.NewResult(req.res, data)
					if err != nil {
						klog.Errorf("The device %s attribute %s  is unsupported, %v", device.Name, req.res.Name, err)
						continue
					}

					for _, msgBus := range d.msgBuses {
						msgBus.Publish(device.Name, *result)
					}
				}

			case *ua.EventNotificationList:
				// do nothing
			default:
				klog.Infof("unknown publish result: %T", res.Value)
			}
		}
	}
}

func (d *OPCUADriver) findEndpoint(device v1alpha1.DeviceConfig) (string, error) {
	protocols := device.Protocols
	properties, ok := protocols[Protocol]
	if !ok {
		return "", fmt.Errorf("opcua protocol properties is not defined")
	}

	endpoint, ok := properties.Data[Endpoint]
	if !ok {
		return "", fmt.Errorf("endpoint not found in the opcua protocol properties")
	}
	return fmt.Sprintf("%v", endpoint), nil
}

func (d *OPCUADriver) toRequest(deviceName string, index int, res v1alpha1.DeviceResource) (*request, error) {
	nodeId, err := getNodeID(res.Attributes, NODE)
	if err != nil {
		return nil, err
	}

	id, err := ua.ParseNodeID(nodeId)
	if err != nil {
		return nil, err
	}

	req := request{
		nodeId: id,
		handle: uint32(index + 42),
		res:    res,
	}

	requests, ok := d.deviceRequests[deviceName]
	if !ok {
		d.deviceRequests[deviceName] = []request{req}
		return &req, nil
	}

	requests = append(requests, req)
	d.deviceRequests[deviceName] = requests
	return &req, nil
}

func (d *OPCUADriver) findRequest(deviceName string, handle uint32) *request {
	requests, ok := d.deviceRequests[deviceName]
	if !ok {
		return nil
	}

	for _, req := range requests {
		if req.handle == handle {
			return &req
		}
	}

	return nil
}

func valueRequest(req *request) *ua.MonitoredItemCreateRequest {
	return opcua.NewMonitoredItemCreateRequestWithDefaults(req.nodeId, ua.AttributeIDValue, req.handle)
}

func getNodeID(attrs v1alpha1.Values, id string) (string, error) {
	identifier, ok := attrs.Data[id]
	if !ok {
		return "", fmt.Errorf("attribute %s does not exist", id)
	}

	return identifier.(string), nil
}
