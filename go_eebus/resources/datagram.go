package resources

const SPECIFICATION_VERSION = "1.0.0"

type DeviceModel struct {
	DeviceType    string         `xml:"deviceType"`
	DeviceAddress string         `xml:"deviceAddress"`
	Description   string         `xml:"description"`
	Entities      []*EntityModel `xml:"entities"`
}

type EntityModel struct {
	EntityType    string          `xml:"entityType"`
	EntityAddress int             `xml:"entityAddress"`
	Description   string          `xml:"description"`
	Features      []*FeatureModel `xml:"features"`
}

type FeatureModel struct {
	FeatureType      string           `xml:"featureType"`
	FeatureAddress   int              `xml:"featureAddress"`
	Role             string           `xml:"role"`
	Description      string           `xml:"description"`
	Functions        []*FunctionModel `xml:"functions"`
	BindingTo        []string
	SubscriptionTo   []string
	MaxBindings      int
	MaxSubscriptions int
}

type FunctionModel struct {
	FunctionName string      `xml:"functionName"`
	ChangeNotify Notifier    `xml:"changeNotify"`
	Function     interface{} `xml:"function"`
}

type DatagramType struct {
	Header  *HeaderType  `xml:"header"`
	Payload *PayloadType `xml:"payload"`
}

type PayloadType struct {
	Cmd *CmdType `xml:"cmd"`
}

type CmdType struct {
	FunctionName string `xml:"functionName"`
	Function     string `xml:"function"`
}

type HeaderType struct {
	SpecificationVersion string              `xml:"specificationVersion"`
	AddressSource        *FeatureAddressType `xml:"addressSource"`
	AddressDestination   *FeatureAddressType `xml:"addressDestination"`
	MsgCounter           int                 `xml:"msgCounter"`
	CmdClassifier        string              `xml:"cmdClassifier"`
	Timestamp            string              `xml:"timestamp"`
	AckRequest           bool                `xml:"ackRequest"`
}

type ComissioningNewSkis struct {
	Skis    string `xml:"skis"`
	Devices string `xml:"devices"`
}

func (device *DeviceModel) CreateNodeManagement(isGateway bool) *FeatureModel {
	subscriptions := []string{}
	bindings := []string{}
	if isGateway {
		subscriptions = append(subscriptions, []string{"ActuatorSwitch", "MeasurementTemp", "MeasurementSolar", "MeasurementBattery", "Measurement", "NodeManagement", "ActuatorSwitch1", "ActuatorSwitch2"}...)
		bindings = append(bindings, []string{"ActuatorSwitch", "ActuatorSwitch1", "ActuatorSwitch2"}...)
	}
	return &FeatureModel{
		FeatureType:      "NodeManagement",
		FeatureAddress:   0,
		Role:             "special",
		SubscriptionTo:   subscriptions,
		BindingTo:        bindings,
		MaxBindings:      128,
		MaxSubscriptions: 128,
		Functions: []*FunctionModel{
			{
				FunctionName: "nodeManagementDetailedDiscoveryData",
				Function: &NodeManagementDetailedDiscovery{
					SpecificationVersionList: []*NodeManagementSpecificationVersionListType{
						{
							SpecificationVersion: SPECIFICATION_VERSION,
						},
					},
					DeviceInformation: &NodeManagementDetailedDiscoveryDeviceInformationType{
						Description: &NetworkManagementDeviceDescriptionDataType{
							DeviceAddress: &DeviceAddressType{
								Device: device.DeviceAddress,
							},
							DeviceType:  device.DeviceType,
							Description: device.Description,
						},
					},
					EntityInformation:  makeEntities(device),
					FeatureInformation: makeFeatures(device),
				},
			},
			{
				FunctionName: "nodeManagementBindingData",
				Function:     &NodeManagementBindingData{},
			},
			{
				FunctionName: "nodeManagementSubscriptionData",
				Function:     &NodeManagementSubscriptionData{},
			},
		},
	}
}
