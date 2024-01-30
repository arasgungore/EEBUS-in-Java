package go_eebus

import (
	"github.com/LMF-DHBW/go_eebus/resources"
	"github.com/LMF-DHBW/go_eebus/spine"
)

type EebusNode struct {
	isGateway       bool
	SpineNode       *spine.SpineNode
	DeviceStructure *resources.DeviceModel
	Update          Updater
}

type Updater func(resources.DatagramType, spine.SpineConnection)

func NewEebusNode(hostname string, isGateway bool, certName string, devId string, brand string, devType string) *EebusNode {
	deviceModel := &resources.DeviceModel{}
	newEebusNode := &EebusNode{isGateway, nil, deviceModel, nil}
	newEebusNode.SpineNode = spine.NewSpineNode(hostname, isGateway, deviceModel, newEebusNode.SubscriptionNofity, certName, devId, brand, devType)
	return newEebusNode
}

func (eebusNode *EebusNode) SubscriptionNofity(datagram resources.DatagramType, conn spine.SpineConnection) {
	if eebusNode.Update != nil {
		eebusNode.Update(datagram, conn)
	}
}

func (eebusNode *EebusNode) Start() {
	eebusNode.SpineNode.Start()
}
