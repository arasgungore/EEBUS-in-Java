package spine

import (
	"encoding/xml"
	"log"
	"strings"
	"time"

	"github.com/LMF-DHBW/go_eebus/resources"
	"github.com/LMF-DHBW/go_eebus/ship"
)

type Notifier func(resources.DatagramType, SpineConnection)
type BindSubscribeNotify func(string, *SpineConnection, *resources.BindSubscribeEntry)

type SpineConnection struct {
	SME                  *ship.SMEInstance
	Address              string
	MsgCounter           int
	OwnDevice            *resources.DeviceModel
	recieveChan          chan resources.DatagramType
	DiscoveryInformation *resources.NodeManagementDetailedDiscovery
	bindSubscribeNotify  BindSubscribeNotify
	bindSubscribeInfo    []*BindSubscribeInfo
	SubscriptionNofity   Notifier
	SubscriptionData     []*SubscriptionData
}

type SubscriptionData struct {
	EntityType     string
	FeatureType    string
	FeatureAddress resources.FeatureAddressType
	FunctionName   string
	CurrentState   string
}

type BindSubscribeInfo struct {
	BindSubscribe      string
	BindSubscribeEntry *resources.BindSubscribeEntry
}

func NewSpineConnection(SME *ship.SMEInstance, ownDevice *resources.DeviceModel, bindSubscribeNotify BindSubscribeNotify, SubscriptionNofity Notifier) *SpineConnection {
	return &SpineConnection{SME, "", 0, ownDevice, make(chan resources.DatagramType), nil, bindSubscribeNotify, nil, SubscriptionNofity, make([]*SubscriptionData, 0)}
}

func (conn *SpineConnection) SendXML(header *resources.HeaderType, payload *resources.PayloadType) {
	conn.MsgCounter++
	conn.SME.Send(resources.DatagramType{header, payload})
}

func (conn *SpineConnection) StartRecieveHandler() {
	conn.SME.Recieve(func(datagram resources.DatagramType) {
		entitiyAddr := datagram.Header.AddressDestination.Entity
		featureAddr := datagram.Header.AddressDestination.Feature
		deviceSource := datagram.Header.AddressSource.Device
		entitiySource := datagram.Header.AddressSource.Entity
		featureSource := datagram.Header.AddressSource.Feature
		isValidRequest := len(conn.OwnDevice.Entities) > entitiyAddr && len(conn.OwnDevice.Entities[entitiyAddr].Features) > featureAddr
		if isValidRequest {
			conn.MsgCounter = resources.Max(conn.MsgCounter, datagram.Header.MsgCounter)

			feature := conn.OwnDevice.Entities[entitiyAddr].Features[featureAddr]
			var function *resources.FunctionModel
			for _, v := range feature.Functions {
				if v.FunctionName == datagram.Payload.Cmd.FunctionName {
					function = v
					break
				}
			}
			switch datagram.Header.CmdClassifier {
			case "reply", "result":
				conn.recieveChan <- datagram
			case "read":
				if conn.requestAllowed("binding", datagram.Header) {
					conn.SendXML(
						conn.OwnDevice.MakeHeader(entitiyAddr, featureAddr, resources.MakeFeatureAddress(deviceSource, entitiySource, featureSource), "result", conn.MsgCounter, false),
						resources.MakePayload(function.FunctionName, function.Function))
				}
			case "write":
				if conn.requestAllowed("binding", datagram.Header) {
					function.ChangeNotify(datagram.Payload.Cmd.FunctionName, datagram.Payload.Cmd.Function, *datagram.Header.AddressDestination)
				}
			case "notify":
				if conn.requestAllowed("subscription", datagram.Header) {
					if len(conn.DiscoveryInformation.FeatureInformation) > featureSource && len(conn.DiscoveryInformation.EntityInformation) > entitiySource {
						featureInList := false
						for i := range conn.SubscriptionData {
							if conn.SubscriptionData[i].FeatureAddress == *datagram.Header.AddressSource && conn.SubscriptionData[i].FunctionName == datagram.Payload.Cmd.FunctionName {
								featureInList = true
								break
							}
						}
						if !featureInList {
							conn.SubscriptionData = append(conn.SubscriptionData, &SubscriptionData{
								EntityType:     conn.DiscoveryInformation.EntityInformation[entitiySource].Description.EntityType,
								FeatureType:    conn.DiscoveryInformation.FeatureInformation[featureSource].Description.FeatureType,
								FeatureAddress: *datagram.Header.AddressSource,
								FunctionName:   datagram.Payload.Cmd.FunctionName,
								CurrentState:   datagram.Payload.Cmd.Function,
							})
						} else {
							for i := range conn.SubscriptionData {
								if conn.SubscriptionData[i].FeatureAddress == *datagram.Header.AddressSource && conn.SubscriptionData[i].FunctionName == datagram.Payload.Cmd.FunctionName {
									conn.SubscriptionData[i].CurrentState = datagram.Payload.Cmd.Function
									break
								}
							}
						}
						conn.SubscriptionNofity(datagram, *conn)
					}
				}
			case "call":
				if datagram.Payload.Cmd.FunctionName == "nodeManagementBindingRequestCall" {
					conn.processBindingRequest(&datagram)
				} else if datagram.Payload.Cmd.FunctionName == "nodeManagementSubscriptionRequestCall" {
					conn.processSubscriptionRequest(&datagram)
				}
			case "comissioning":
				if conn.DiscoveryInformation.DeviceInformation.Description.DeviceType == "Gateway" {
					if datagram.Payload.Cmd.FunctionName == "saveSkis" {
						log.Println("Saving new SKIs")

						var Function *resources.ComissioningNewSkis
						err := xml.Unmarshal([]byte(datagram.Payload.Cmd.Function), &Function)
						if err == nil {
							ship.WriteSkis(strings.Split(Function.Skis, ";"), strings.Split(Function.Devices, ";"))
						}
					}
				}
			}
		}
	})
}

func (conn *SpineConnection) requestAllowed(bindSubscribe string, header *resources.HeaderType) bool {
	entitiyAddr := header.AddressDestination.Entity
	featureAddr := header.AddressDestination.Feature
	entitiySource := header.AddressSource.Entity
	featureSource := header.AddressSource.Feature
	if entitiyAddr == 0 && featureAddr == 0 {
		return true
	}
	for _, info := range conn.bindSubscribeInfo {
		if info.BindSubscribe == bindSubscribe &&
			entitiyAddr == info.BindSubscribeEntry.ServerAddress.Entity && featureAddr == info.BindSubscribeEntry.ServerAddress.Feature &&
			entitiySource == info.BindSubscribeEntry.ClientAddress.Entity && featureSource == info.BindSubscribeEntry.ClientAddress.Feature {
			return true
		}
	}
	return false
}

func (conn *SpineConnection) RecieveTimeout(seconds int) (resources.DatagramType, bool) {
	var res resources.DatagramType
	err := false
	select {
	case res = <-conn.recieveChan:
		err = true
	case <-time.After(time.Duration(seconds) * time.Second):
		err = false
	}
	return res, err
}

func (conn *SpineConnection) CountBindings(serverAddr resources.FeatureAddressType) int {
	numBindings := 0
	for _, bindSub := range conn.bindSubscribeInfo {
		if bindSub.BindSubscribe == "binding" && bindSub.BindSubscribeEntry.ServerAddress == serverAddr {
			numBindings++
		}
	}
	return numBindings
}

func (conn *SpineConnection) CountSubscriptions(serverAddr resources.FeatureAddressType) int {
	numSubscriptions := 0
	for _, bindSub := range conn.bindSubscribeInfo {
		if bindSub.BindSubscribe == "subscription" && bindSub.BindSubscribeEntry.ServerAddress == serverAddr {
			numSubscriptions++
		}
	}
	return numSubscriptions
}
