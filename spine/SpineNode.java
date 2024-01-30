package spine;

import resources.*;

import java.util.ArrayList;
import java.util.List;
import java.util.logging.Logger;

public class SpineNode {

    public ship.ShipNode shipNode;
    public List<SpineConnection> connections;
    public resources.DeviceModel deviceStructure;
    public List<BindSubscribe> bindings;
    public List<BindSubscribe> subscriptions;
    public Notifier subscriptionNofity;

    public SpineNode(String hostname, boolean isGateway, DeviceModel deviceModel, Notifier subscriptionNofity, String certName, String devId, String brand, String devType) {
        this.shipNode = new ship.ShipNode(hostname, isGateway, certName, devId, brand, devType);
        this.connections = new ArrayList<>();
        this.deviceStructure = deviceModel;
        this.bindings = new ArrayList<>();
        this.subscriptions = new ArrayList<>();
        this.subscriptionNofity = subscriptionNofity;
    }

    public void start() {
        shipNode.setSpineConnectionNotify(this::newConnection);
        shipNode.setSpineCloseHandler(this::closeHandler);
        shipNode.start();
    }

    public void newConnection(ship.SMEInstance SME, String newSki) {
        SpineConnection newSpineConnection = new SpineConnection(SME, deviceStructure, this::newBindSubscribe, subscriptionNofity);
        
        new Thread(() -> {
            newSpineConnection.startDetailedDiscovery();

            if (shipNode.isGateway()) {
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                List<String> skis = ship.readSkis();
                List<String> devices = ship.readSkis();
                newSpineConnection.sendXML(newSpineConnection.ownDevice.makeHeader(0, 0, resources.makeFeatureAddress("", 0, 0), "comissioning", newSpineConnection.msgCounter, false),
                        resources.makePayload("saveSkis", new resources.ComissioningNewSkis(String.join(";", skis), String.join(";", devices))));

                if (!newSki.isEmpty()) {
                    ship.writeSkis(skis.add(newSki), devices.add(newSpineConnection.address));

                    List<String> updatedSkis = ship.readSkis();
                    Logger.getGlobal().info("Sending new SKIs");
                    for (SpineConnection conn : connections) {
                        conn.sendXML(conn.ownDevice.makeHeader(0, 0, resources.makeFeatureAddress("", 0, 0), "comissioning", conn.msgCounter, false),
                                resources.makePayload("saveSkis", new resources.ComissioningNewSkis(String.join(";", updatedSkis), String.join(";", devices))));
                    }
                }
            }

            connections.add(newSpineConnection);
        }).start();

        newSpineConnection.startRecieveHandler();
    }

    public void newBindSubscribe(String bindSubscribe, SpineConnection conn, resources.BindSubscribeEntry entry) {
        if ("binding".equals(bindSubscribe)) {
            Logger.getGlobal().info("added binding");
            bindings.add(new BindSubscribe(conn, entry));

            NodeManagementBindingData ownBindings = (NodeManagementBindingData) deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(1).getFunction();
            ownBindings.getBindingEntries().add(entry);
            deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(1).setFunction(ownBindings);

        } else if ("subscription".equals(bindSubscribe)) {
            subscriptions.add(new BindSubscribe(conn, entry));

            NodeManagementSubscriptionData ownSubscriptions = (NodeManagementSubscriptionData) deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(2).getFunction();
            ownSubscriptions.getSubscriptionEntries().add(entry);
            deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(2).setFunction(ownSubscriptions);
        }

        for (BindSubscribe e : subscriptions) {
            if (e.getBindSubscribeEntry().getServerAddress().getFeature() == 0 && e.getBindSubscribeEntry().getServerAddress().getEntity() == 0) {
                e.getConn().sendXML(
                        e.getConn().ownDevice.makeHeader(0, 0, resources.makeFeatureAddress(
                                e.getBindSubscribeEntry().getClientAddress().getDevice(),
                                e.getBindSubscribeEntry().getClientAddress().getEntity(),
                                e.getBindSubscribeEntry().getClientAddress().getFeature()
                        ), "notify", e.getConn().msgCounter, false),
                        resources.makePayload("nodeManagementBindingData", deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(1).getFunction()));
                e.getConn().sendXML(
                        e.getConn().ownDevice.makeHeader(0, 0, resources.makeFeatureAddress(
                                e.getBindSubscribeEntry().getClientAddress().getDevice(),
                                e.getBindSubscribeEntry().getClientAddress().getEntity(),
                                e.getBindSubscribeEntry().getClientAddress().getFeature()
                        ), "notify", e.getConn().msgCounter, false),
                        resources.makePayload("nodeManagementSubscriptionData", deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(2).getFunction()));
            }
        }
    }

    public void closeHandler(ship.SMEInstance SME) {
        for (int i = 0; i < connections.size(); i++) {
            SpineConnection e = connections.get(i);
            if (e.getSME() == SME) {
                int offset = 0;
                for (int j = 0; j < bindings.size(); j++) {
                    BindSubscribe binding = bindings.get(j);
                    if (binding.getConn() == connections.get(i)) {
                        j += offset;
                        bindings.set(j, bindings.get(bindings.size() - 1));
                        bindings.remove(bindings.size() - 1);

                        NodeManagementBindingData ownBindings = (NodeManagementBindingData) deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(1).getFunction();
                        ownBindings.getBindingEntries().set(j, ownBindings.getBindingEntries().get(ownBindings.getBindingEntries().size() - 1));
                        ownBindings.getBindingEntries().remove(ownBindings.getBindingEntries().size() - 1);
                        deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(1).setFunction(ownBindings);
                        offset -= 1;
                    }
                }

                offset = 0;
                for (int j = 0; j < subscriptions.size(); j++) {
                    BindSubscribe subscription = subscriptions.get(j);
                    if (subscription.getConn() == connections.get(i)) {
                        j += offset;
                        subscriptions.set(j, subscriptions.get(subscriptions.size() - 1));
                        subscriptions.remove(subscriptions.size() - 1);

                        NodeManagementSubscriptionData ownSubscriptions = (NodeManagementSubscriptionData) deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(2).getFunction();
                        ownSubscriptions.getSubscriptionEntries().set(j, ownSubscriptions.getSubscriptionEntries().get(ownSubscriptions.getSubscriptionEntries().size() - 1));
                        ownSubscriptions.getSubscriptionEntries().remove(ownSubscriptions.getSubscriptionEntries().size() - 1);
                        deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(2).setFunction(ownSubscriptions);
                        offset -= 1;
                    }
                }

                for (BindSubscribe e1 : subscriptions) {
                    if (e1.getBindSubscribeEntry().getServerAddress().getFeature() == 0 && e1.getBindSubscribeEntry().getServerAddress().getEntity() == 0) {
                        e1.getConn().sendXML(
                                e1.getConn().ownDevice.makeHeader(0, 0, resources.makeFeatureAddress(
                                        e1.getBindSubscribeEntry().getClientAddress().getDevice(),
                                        e1.getBindSubscribeEntry().getClientAddress().getEntity(),
                                        e1.getBindSubscribeEntry().getClientAddress().getFeature()
                                ), "notify", e1.getConn().msgCounter, false),
                                resources.makePayload("nodeManagementBindingData", deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(1).getFunction()));
                        e1.getConn().sendXML(
                                e1.getConn().ownDevice.makeHeader(0, 0, resources.makeFeatureAddress(
                                        e1.getBindSubscribeEntry().getClientAddress().getDevice(),
                                        e1.getBindSubscribeEntry().getClientAddress().getEntity(),
                                        e1.getBindSubscribeEntry().getClientAddress().getFeature()
                                ), "notify", e1.getConn().msgCounter, false),
                                resources.makePayload("nodeManagementSubscriptionData", deviceStructure.getEntities().get(0).getFeatures().get(0).getFunctions().get(2).getFunction()));
                    }
                }

                connections.set(i, connections.get(connections.size() - 1));
                connections.remove(connections.size() - 1);

                shipNode.getSME()[i] = shipNode.getSME()[shipNode.getSME().length - 1];
                shipNode.getSME() = Arrays.copyOf(shipNode.getSME(), shipNode.getSME().length - 1);

                Logger.getGlobal().info("Connection closed!");
                break;
            }
        }
    }
}
