package resources

type Notifier func(string, string, FeatureAddressType)

type FunctionElement struct {
	Function string `xml:"function"`
}

type DescriptionElement struct {
	Label       string `xml:"label"`
	Description string `xml:"description"`
}

func ActuatorSwitch(label string, description string, ChangeNotify Notifier) []*FunctionModel {
	return []*FunctionModel{
		{
			FunctionName: "actuatorSwitchData",
			Function: &FunctionElement{
				Function: "off",
			},
			ChangeNotify: ChangeNotify,
		},
		{
			FunctionName: "actuatorSwitchDescriptionData",
			Function: &DescriptionElement{
				label,
				description,
			},
		},
	}
}
