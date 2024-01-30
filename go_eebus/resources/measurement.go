package resources

type TimePeriodType struct {
	StartTime string `xml:"startTime"`
	EndTime   string `xml:"endTime"`
}

type MeasurementDataType struct {
	ValueType        string         `xml:"valueType"`
	Timestamp        string         `xml:"timestamp"`
	Value            float64        `xml:"value"`
	EvaluationPeriod TimePeriodType `xml:"evaluationPeriod"`
	ValueSource      string         `xml:"valueSource"`
	ValueTendency    string         `xml:"valueTendency"`
	ValueState       string         `xml:"valueState"`
}

type MeasurementDescriptionDataType struct {
	MeasurementType string `xml:"measurementType"`
	Unit            string `xml:"unit"`
	ScopeType       string `xml:"scopeType"`
	Label           string `xml:"label"`
	Description     string `xml:"description"`
}

func Measurement(MeasurementType string, Unit string, ScopeType string, Label string, Description string) []*FunctionModel {
	return []*FunctionModel{
		{
			FunctionName: "measurementData",
			Function:     &MeasurementDataType{},
		},
		{
			FunctionName: "measurementDescription",
			Function: &MeasurementDescriptionDataType{
				MeasurementType,
				Unit,
				ScopeType,
				Label,
				Description,
			},
		},
	}
}
