package resources

import (
	"encoding/xml"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		debug.PrintStack()
		os.Exit(1)
	}
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func xmlToString(in interface{}) string {
	bytes, err := xml.Marshal(in)
	CheckError(err)
	return string(bytes)
}

func timestampNow() string {
	current_time := time.Now()

	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.0Z",
		current_time.Year(), current_time.Month(), current_time.Day(),
		current_time.Hour(), current_time.Minute(), current_time.Second())
}

func makeEntities(device *DeviceModel) []*NodeManagementDetailedDiscoveryEntityInformationType {
	var all []*NodeManagementDetailedDiscoveryEntityInformationType
	entities := device.Entities
	for _, entity := range entities {
		all = append(all, &NodeManagementDetailedDiscoveryEntityInformationType{
			Description: &NetworkManagementEntityDescritpionDataType{
				EntityAddress: &EntityAddressType{
					Device: device.DeviceAddress,
					Entity: entity.EntityAddress,
				},
				EntityType:  entity.EntityType,
				Description: entity.Description,
			},
		})
	}
	return all
}

func makeFeatures(device *DeviceModel) []*NodeManagementDetailedDiscoveryFeatureInformationType {
	var all []*NodeManagementDetailedDiscoveryFeatureInformationType
	entities := device.Entities
	for _, entity := range entities {
		for _, feature := range entity.Features {
			all = append(all, &NodeManagementDetailedDiscoveryFeatureInformationType{
				Description: &NetworkManagementFeatureInformationType{
					FeatureAddress: &FeatureAddressType{
						Device:  device.DeviceAddress,
						Entity:  entity.EntityAddress,
						Feature: feature.FeatureAddress,
					},
					FeatureType: feature.FeatureType,
					Description: feature.Description,
					Role:        feature.Role,
				},
			})
		}
	}
	return all
}

func (device *DeviceModel) MakeHeader(entity int, feature int, addressDestination *FeatureAddressType, cmdClassifier string, msgCounter int, ackRequest bool) *HeaderType {
	return &HeaderType{
		SpecificationVersion: SPECIFICATION_VERSION,
		AddressSource: &FeatureAddressType{
			Device:  device.DeviceAddress,
			Entity:  entity,
			Feature: feature,
		},
		AddressDestination: addressDestination,
		MsgCounter:         msgCounter,
		CmdClassifier:      cmdClassifier,
		Timestamp:          timestampNow(),
		AckRequest:         ackRequest,
	}
}

func MakePayload(FunctionName string, Function interface{}) *PayloadType {
	return &PayloadType{
		Cmd: &CmdType{
			FunctionName,
			xmlToString(Function),
		},
	}
}

func MakeFeatureAddress(device string, entity int, feature int) *FeatureAddressType {
	return &FeatureAddressType{
		device,
		entity,
		feature,
	}
}
