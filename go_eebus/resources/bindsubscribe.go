package resources

type NodeManagementBindingData struct {
	BindingEntries []*BindSubscribeEntry `xml:"bindingEntries"`
}

type BindSubscribeEntry struct {
	ClientAddress FeatureAddressType `xml:"clientAddress"`
	ServerAddress FeatureAddressType `xml:"serverAddress"`
}

type NodeManagementSubscriptionData struct {
	SubscriptionEntries []*BindSubscribeEntry `xml:"subscriptionEntries"`
}

type NodeManagementBindingRequestCall struct {
	BindingRequest *BindingManagementRequestCallType `xml:"bindingRequest"`
}

type BindingManagementRequestCallType struct {
	ClientAddress     *FeatureAddressType `xml:"clientAddress"`
	ServerAddress     *FeatureAddressType `xml:"serverAddress"`
	ServerFeatureType string              `xml:"serverFeatureType"`
}

type NodeManagementSubscriptionRequestCall struct {
	SubscriptionRequest *SubscriptionManagementRequestCallType `xml:"subscriptionRequest"`
}

type SubscriptionManagementRequestCallType struct {
	ClientAddress     *FeatureAddressType `xml:"clientAddress"`
	ServerAddress     *FeatureAddressType `xml:"serverAddress"`
	ServerFeatureType string              `xml:"serverFeatureType"`
}

type ResultElement struct {
	ErrorNumber int    `xml:"errorNumber"`
	Description string `xml:"description"`
}

func ResultData(errorNumber int, description string) *ResultElement {
	return &ResultElement{
		ErrorNumber: errorNumber,
		Description: description,
	}
}
