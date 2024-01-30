package resources

type NodeManagementDetailedDiscovery struct {
	SpecificationVersionList []*NodeManagementSpecificationVersionListType            `xml:"specificationVersionList"`
	DeviceInformation        *NodeManagementDetailedDiscoveryDeviceInformationType    `xml:"deviceInformation"`
	EntityInformation        []*NodeManagementDetailedDiscoveryEntityInformationType  `xml:"entityInformation"`
	FeatureInformation       []*NodeManagementDetailedDiscoveryFeatureInformationType `xml:"featureInformation"`
}

type NodeManagementSpecificationVersionListType struct {
	SpecificationVersion string `xml:"specificationVersion"`
}

type NodeManagementDetailedDiscoveryDeviceInformationType struct {
	Description *NetworkManagementDeviceDescriptionDataType `xml:"description"`
}

type NodeManagementDetailedDiscoveryEntityInformationType struct {
	Description *NetworkManagementEntityDescritpionDataType `xml:"description"`
}

type NodeManagementDetailedDiscoveryFeatureInformationType struct {
	Description *NetworkManagementFeatureInformationType `xml:"description"`
}

type NetworkManagementDeviceDescriptionDataType struct {
	DeviceAddress *DeviceAddressType `xml:"deviceAddress"`
	DeviceType    string             `xml:"deviceType"`
	Description   string             `xml:"description"`
}

type NetworkManagementEntityDescritpionDataType struct {
	EntityAddress *EntityAddressType `xml:"entityAddress"`
	EntityType    string             `xml:"entityType"`
	Description   string             `xml:"description"`
}

type NetworkManagementFeatureInformationType struct {
	FeatureAddress    *FeatureAddressType   `xml:"featureAddress"`
	FeatureType       string                `xml:"featureType"`
	Role              string                `xml:"role"`
	SupportedFunction *FunctionPropertyType `xml:"supportedFunction"`
	Description       string                `xml:"description"`
}

type FeatureAddressType struct {
	Device  string `xml:"device"`
	Entity  int    `xml:"entity"`
	Feature int    `xml:"feature"`
}

type EntityAddressType struct {
	Device string `xml:"device"`
	Entity int    `xml:"entity"`
}

type DeviceAddressType struct {
	Device string `xml:"device"`
}

type FunctionPropertyType struct {
	Function           string `xml:"function"`
	PossibleOperations string `xml:"possibleOperations"`
}
