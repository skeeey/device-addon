package opcua

import (
	"context"
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/config/device"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/models"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/msgbus"
	"github.com/skeeey/device-addon/pkg/addon/spoke/agent/util"
	"k8s.io/klog/v2"
)

type request struct {
	nodeId *ua.NodeID
	handle uint32
	res    device.DeviceResource
}

type OPCUADriver struct {
	sync.Mutex
	serverInfo       *OPCUAServerInfo
	msgBuses         []msgbus.MessageBus
	devices          map[string]models.Device
	deviceRequests   map[string][]request
	deviceSubCancels map[string]context.CancelFunc
}

func NewOPCUADriver() *OPCUADriver {
	return &OPCUADriver{
		devices:          make(map[string]models.Device),
		deviceRequests:   make(map[string][]request),
		deviceSubCancels: make(map[string]context.CancelFunc),
	}
}

func (d *OPCUADriver) Initialize(driverInfo config.DriverInfo, msgBuses []msgbus.MessageBus) error {
	var serverInfo = &OPCUAServerInfo{}
	if err := util.LoadConfig(path.Join(driverInfo.ConfigDir, config.DriverConfigFileName), serverInfo); err != nil {
		return err
	}

	d.msgBuses = msgBuses
	d.serverInfo = serverInfo
	return nil
}

func (d *OPCUADriver) Start() error {
	//do nothing
	return nil
}

func (d *OPCUADriver) Stop() error {
	return nil
}

func (d *OPCUADriver) AddDevice(device models.Device) error {
	d.Lock()
	defer d.Unlock()

	_, ok := d.devices[device.DeviceName]
	if !ok {
		go func() {
			if err := d.startSubscription(device); err != nil {
				klog.Errorf("failed to sub device %s, %v", device.DeviceName, err)
			}
		}()

		d.devices[device.DeviceName] = device
	}

	return nil
}

func (d *OPCUADriver) UpdateDevice(device models.Device) error {
	//TODO
	return nil
}

func (d *OPCUADriver) RemoveDevice(deviceName string) error {
	//TODO
	return nil
}

func (d *OPCUADriver) HandleCommands(deviceName string, command models.Command) error {
	//TODO
	return nil
}

func (d *OPCUADriver) startSubscription(device models.Device) error {
	_, ok := d.deviceSubCancels[device.DeviceName]
	if ok {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d.deviceSubCancels[device.DeviceName] = cancel

	endpoint, err := d.findEndpoint(device.Device)
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

	for index, deviceResource := range device.DeviceProfile.DeviceResources {
		req, err := d.toRequest(device.DeviceName, index, deviceResource)
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

					req := d.findRequest(device.DeviceName, item.ClientHandle)
					if req == nil {
						continue
					}

					result, err := util.NewResult(req.res, data)
					if err != nil {
						klog.Errorf("The device %s attribute %s  is unsupported, %v", device.DeviceName, req.res.Name, err)
						continue
					}

					for _, msgBus := range d.msgBuses {
						msgBus.Publish(device.DeviceName, *result)
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

func (d *OPCUADriver) findEndpoint(device *device.Device) (string, error) {
	protocols := device.Protocols
	properties, ok := protocols[Protocol]
	if !ok {
		return "", fmt.Errorf("opcua protocol properties is not defined")
	}

	endpoint, ok := properties[Endpoint]
	if !ok {
		return "", fmt.Errorf("endpoint not found in the opcua protocol properties")
	}
	return fmt.Sprintf("%v", endpoint), nil
}

func (d *OPCUADriver) toRequest(deviceName string, index int, res device.DeviceResource) (*request, error) {
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

func getNodeID(attrs map[string]interface{}, id string) (string, error) {
	identifier, ok := attrs[id]
	if !ok {
		return "", fmt.Errorf("attribute %s does not exist", id)
	}

	return identifier.(string), nil
}
