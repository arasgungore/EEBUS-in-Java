package spine

import (
	"log"
	"strings"
	"time"

	"github.com/LMF-DHBW/go_eebus/resources"
	"github.com/LMF-DHBW/go_eebus/ship"
)

type SpineNode struct {
	ShipNode           *ship.ShipNode
	Connections        []*SpineConnection
	DeviceStructure    *resources.DeviceModel
	Bindings           []*BindSubscribe
	Subscriptions      []*BindSubscribe
	SubscriptionNofity Notifier
}

type BindSubscribe struct {
	Conn               *SpineConnection
	BindSubscribeEntry *resources.BindSubscribeEntry
}

func (bindSubscribe BindSubscribe) Send(msgClassifier string, payload *resources.PayloadType) {
	srv := bindSubscribe.BindSubscribeEntry.ServerAddress
	clt := bindSubscribe.BindSubscribeEntry.ClientAddress
	bindSubscribe.Conn.SendXML(
		bindSubscribe.Conn.OwnDevice.MakeHeader(srv.Entity, srv.Feature,
			resources.MakeFeatureAddress(clt.Device, clt.Entity, clt.Feature),
			msgClassifier, bindSubscribe.Conn.MsgCounter, false),
		payload)
}

func NewSpineNode(hostname string, isGateway bool, deviceModel *resources.DeviceModel, SubscriptionNofity Notifier, certName string, devId string, brand string, devType string) *SpineNode {
	return &SpineNode{ship.NewShipNode(hostname, isGateway, certName, devId, brand, devType), make([]*SpineConnection, 0), deviceModel, make([]*BindSubscribe, 0), make([]*BindSubscribe, 0), SubscriptionNofity}
}

func (spineNode *SpineNode) Start() {
	spineNode.ShipNode.SpineConnectionNotify = spineNode.newConnection
	spineNode.ShipNode.SpineCloseHandler = spineNode.closeHandler
	spineNode.ShipNode.Start()
}

func (spineNode *SpineNode) newConnection(SME *ship.SMEInstance, newSki string) {

	newSpineConnection := NewSpineConnection(SME, spineNode.DeviceStructure, spineNode.newBindSubscribe, spineNode.SubscriptionNofity)

	go func() {

		newSpineConnection.StartDetailedDiscovery()

		if spineNode.ShipNode.IsGateway {

			time.Sleep(time.Second / 10)
			skis, devices := ship.ReadSkis()
			newSpineConnection.SendXML(newSpineConnection.OwnDevice.MakeHeader(0, 0, resources.MakeFeatureAddress("", 0, 0), "comissioning", newSpineConnection.MsgCounter, false), resources.MakePayload("saveSkis", &resources.ComissioningNewSkis{
				Skis:    strings.Join(skis, ";"),
				Devices: strings.Join(devices, ";"),
			}))

			if newSki != "" {
				ship.WriteSkis(append(skis, newSki), append(devices, newSpineConnection.Address))

				skis, devices := ship.ReadSkis()
				log.Println("Sending new SKIs")
				for _, conn := range spineNode.Connections {
					conn.SendXML(conn.OwnDevice.MakeHeader(0, 0, resources.MakeFeatureAddress("", 0, 0), "comissioning", conn.MsgCounter, false), resources.MakePayload("saveSkis", &resources.ComissioningNewSkis{
						Skis:    strings.Join(skis, ";"),
						Devices: strings.Join(devices, ";"),
					}))
				}
			}
		}

		spineNode.Connections = append(spineNode.Connections, newSpineConnection)

	}()

	newSpineConnection.StartRecieveHandler()
}

func (spineNode *SpineNode) newBindSubscribe(bindSubscribe string, conn *SpineConnection, entry *resources.BindSubscribeEntry) {
	if bindSubscribe == "binding" {
		log.Println("added binding")
		spineNode.Bindings = append(spineNode.Bindings, &BindSubscribe{
			conn, entry,
		})
		// Add to binding list for bind information
		ownBindings := spineNode.DeviceStructure.Entities[0].Features[0].Functions[1].Function.(*resources.NodeManagementBindingData)
		ownBindings.BindingEntries = append(ownBindings.BindingEntries, entry)
		spineNode.DeviceStructure.Entities[0].Features[0].Functions[1].Function = ownBindings

	} else if bindSubscribe == "subscription" {
		spineNode.Subscriptions = append(spineNode.Subscriptions, &BindSubscribe{
			conn, entry,
		})
		// Add to subscription list for subscription information
		ownSubscriptions := spineNode.DeviceStructure.Entities[0].Features[0].Functions[2].Function.(*resources.NodeManagementSubscriptionData)
		ownSubscriptions.SubscriptionEntries = append(ownSubscriptions.SubscriptionEntries, entry)
		spineNode.DeviceStructure.Entities[0].Features[0].Functions[2].Function = ownSubscriptions
	}

	for _, e := range spineNode.Subscriptions {
		// Send with e.Conn from e.BindSubscribeEntry Address source to destination

		// Only send to right partners
		if e.BindSubscribeEntry.ServerAddress.Feature == 0 && e.BindSubscribeEntry.ServerAddress.Entity == 0 {
			e.Conn.SendXML(
				e.Conn.OwnDevice.MakeHeader(0, 0, resources.MakeFeatureAddress(e.BindSubscribeEntry.ClientAddress.Device, e.BindSubscribeEntry.ClientAddress.Entity, e.BindSubscribeEntry.ClientAddress.Feature), "notify", e.Conn.MsgCounter, false),
				resources.MakePayload("nodeManagementBindingData", spineNode.DeviceStructure.Entities[0].Features[0].Functions[1].Function))
			e.Conn.SendXML(
				e.Conn.OwnDevice.MakeHeader(0, 0, resources.MakeFeatureAddress(e.BindSubscribeEntry.ClientAddress.Device, e.BindSubscribeEntry.ClientAddress.Entity, e.BindSubscribeEntry.ClientAddress.Feature), "notify", e.Conn.MsgCounter, false),
				resources.MakePayload("nodeManagementSubscriptionData", spineNode.DeviceStructure.Entities[0].Features[0].Functions[2].Function))
		}
	}
}

// If connection is closed -> delete it from SPINE connection list
func (spineNode *SpineNode) closeHandler(SME *ship.SMEInstance) {
	for i, e := range spineNode.Connections {
		if e.SME == SME {
			offset := 0
			for j, binding := range spineNode.Bindings {
				if binding.Conn == spineNode.Connections[i] {
					j += offset
					spineNode.Bindings[j] = spineNode.Bindings[len(spineNode.Bindings)-1]
					spineNode.Bindings = spineNode.Bindings[:len(spineNode.Bindings)-1]

					ownBindings := spineNode.DeviceStructure.Entities[0].Features[0].Functions[1].Function.(*resources.NodeManagementBindingData)
					ownBindings.BindingEntries[j] = ownBindings.BindingEntries[len(ownBindings.BindingEntries)-1]
					ownBindings.BindingEntries = ownBindings.BindingEntries[:len(ownBindings.BindingEntries)-1]
					spineNode.DeviceStructure.Entities[0].Features[0].Functions[1].Function = ownBindings
					offset -= 1
				}
			}

			offset = 0
			for j, subscription := range spineNode.Subscriptions {
				if subscription.Conn == spineNode.Connections[i] {
					j += offset
					spineNode.Subscriptions[j] = spineNode.Subscriptions[len(spineNode.Subscriptions)-1]
					spineNode.Subscriptions = spineNode.Subscriptions[:len(spineNode.Subscriptions)-1]

					ownSubscriptions := spineNode.DeviceStructure.Entities[0].Features[0].Functions[2].Function.(*resources.NodeManagementSubscriptionData)
					ownSubscriptions.SubscriptionEntries[j] = ownSubscriptions.SubscriptionEntries[len(ownSubscriptions.SubscriptionEntries)-1]
					ownSubscriptions.SubscriptionEntries = ownSubscriptions.SubscriptionEntries[:len(ownSubscriptions.SubscriptionEntries)-1]
					spineNode.DeviceStructure.Entities[0].Features[0].Functions[2].Function = ownSubscriptions
					offset -= 1
				}
			}

			for _, e := range spineNode.Subscriptions {
				if e.BindSubscribeEntry.ServerAddress.Feature == 0 && e.BindSubscribeEntry.ServerAddress.Entity == 0 {
					e.Conn.SendXML(
						e.Conn.OwnDevice.MakeHeader(0, 0, resources.MakeFeatureAddress(e.BindSubscribeEntry.ClientAddress.Device, e.BindSubscribeEntry.ClientAddress.Entity, e.BindSubscribeEntry.ClientAddress.Feature), "notify", e.Conn.MsgCounter, false),
						resources.MakePayload("nodeManagementBindingData", spineNode.DeviceStructure.Entities[0].Features[0].Functions[1].Function))
					e.Conn.SendXML(
						e.Conn.OwnDevice.MakeHeader(0, 0, resources.MakeFeatureAddress(e.BindSubscribeEntry.ClientAddress.Device, e.BindSubscribeEntry.ClientAddress.Entity, e.BindSubscribeEntry.ClientAddress.Feature), "notify", e.Conn.MsgCounter, false),
						resources.MakePayload("nodeManagementSubscriptionData", spineNode.DeviceStructure.Entities[0].Features[0].Functions[2].Function))
				}
			}

			spineNode.Connections[i] = spineNode.Connections[len(spineNode.Connections)-1]
			spineNode.Connections = spineNode.Connections[:len(spineNode.Connections)-1]

			spineNode.ShipNode.SME[i] = spineNode.ShipNode.SME[len(spineNode.ShipNode.SME)-1]
			spineNode.ShipNode.SME = spineNode.ShipNode.SME[:len(spineNode.ShipNode.SME)-1]

			log.Println("Connection closed!")
			break
		}
	}
}
