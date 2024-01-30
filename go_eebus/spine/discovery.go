package spine

import (
	"encoding/xml"

	"github.com/LMF-DHBW/go_eebus/resources"
)

func (conn *SpineConnection) StartDetailedDiscovery() {
	funct := conn.OwnDevice.Entities[0].Features[0].Functions[0].Function
	conn.SendXML(conn.OwnDevice.MakeHeader(0, 0, resources.MakeFeatureAddress("", 0, 0), "read", conn.MsgCounter, false), resources.MakePayload("nodeManagementDetailedDiscoveryData", funct))
	answer, ok := conn.RecieveTimeout(10)
	if ok {
		conn.Address = answer.Header.AddressSource.Device
		var Function *resources.NodeManagementDetailedDiscovery
		err := xml.Unmarshal([]byte(answer.Payload.Cmd.Function), &Function)
		if err == nil {
			conn.DiscoveryInformation = Function
			// Device discovery request correct
			for _, Entity := range conn.OwnDevice.Entities {
				for _, Feature := range Entity.Features {
					for _, FeatureInformation := range Function.FeatureInformation {
						destDevice := Function.DeviceInformation.Description.DeviceAddress.Device
						destEntity := FeatureInformation.Description.FeatureAddress.Entity
						destFeature := FeatureInformation.Description.FeatureAddress.Feature

						ownAddr := resources.FeatureAddressType{
							conn.OwnDevice.DeviceAddress,
							Entity.EntityAddress,
							Feature.FeatureAddress,
						}
						numBindings := conn.CountBindings(ownAddr)

						numSubscriptions := conn.CountSubscriptions(ownAddr)

						if Feature.BindingTo != nil &&
							resources.StringInSlice(FeatureInformation.Description.FeatureType, Feature.BindingTo) &&
							(FeatureInformation.Description.Role == "server" || FeatureInformation.Description.Role == "special") &&
							conn.OwnDevice.Entities[Entity.EntityAddress].Features[Feature.FeatureAddress].MaxBindings > numBindings {
							conn.sendBindingRequest(Entity.EntityAddress, Feature.FeatureAddress,
								resources.MakeFeatureAddress(destDevice, destEntity, destFeature),
								FeatureInformation.Description.FeatureType)
						}

						if Feature.SubscriptionTo != nil &&
							resources.StringInSlice(FeatureInformation.Description.FeatureType, Feature.SubscriptionTo) &&
							(FeatureInformation.Description.Role == "server" || FeatureInformation.Description.Role == "special") &&
							conn.OwnDevice.Entities[Entity.EntityAddress].Features[Feature.FeatureAddress].MaxSubscriptions > numSubscriptions {
							conn.sendSubscriptionRequest(Entity.EntityAddress, Feature.FeatureAddress,
								resources.MakeFeatureAddress(destDevice, destEntity, destFeature),
								FeatureInformation.Description.FeatureType)
						}
					}
				}
			}

		}
	}
}
