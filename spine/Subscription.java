package spine;

import resources.*;

import javax.xml.bind.JAXBContext;
import javax.xml.bind.JAXBException;
import javax.xml.bind.Unmarshaller;
import java.io.StringReader;
import java.util.logging.Logger;

public class Subscription {

    public SpineConnection conn;

    public Subscription(SpineConnection conn) {
        this.conn = conn;
    }

    public void sendSubscriptionRequest(int entityAddress, int featureAddress, FeatureAddressType destinationAddr, String featureType) {
        FeatureAddressType clientAddr = resources.makeFeatureAddress(conn.ownDevice.getDeviceAddress(), entityAddress, featureAddress);
        conn.sendXML(
                conn.ownDevice.makeHeader(entityAddress, featureAddress, destinationAddr, "call", conn.msgCounter, true),
                resources.makePayload("nodeManagementSubscriptionRequestCall", new NodeManagementSubscriptionRequestCall(
                        new SubscriptionManagementRequestCallType(clientAddr, destinationAddr, featureType)
                )));

        // Receive response
        DatagramType answer = conn.receiveTimeout(10);
        if (answer != null) {
            Unmarshaller unmarshaller;
            try {
                unmarshaller = JAXBContext.newInstance(ResultElement.class).createUnmarshaller();
                ResultElement function = (ResultElement) unmarshaller.unmarshal(new StringReader(answer.getPayload().getCmd().getFunction()));
                if (function.getErrorNumber() == 0) {
                    Logger.getGlobal().info("Accepted subscription to: " + destinationAddr.getDevice());

                    BindSubscribeEntry newEntry = new BindSubscribeEntry(clientAddr, destinationAddr);
                    conn.bindSubscribeInfo.add(new BindSubscribeInfo("subscription", newEntry));
                    conn.bindSubscribeNotify("subscription", conn, newEntry);
                }
            } catch (JAXBException e) {
                e.printStackTrace();
            }
        }
    }

    public void processSubscriptionRequest(DatagramType datagram) {
        try {
            Unmarshaller unmarshaller = JAXBContext.newInstance(NodeManagementSubscriptionRequestCall.class).createUnmarshaller();
            NodeManagementSubscriptionRequestCall function = (NodeManagementSubscriptionRequestCall) unmarshaller.unmarshal(new StringReader(datagram.getPayload().getCmd().getFunction()));

            int entityAddr = function.getSubscriptionRequest().getServerAddress().getEntity();
            int featureAddr = function.getSubscriptionRequest().getServerAddress().getFeature();
            boolean isValidRequest = conn.ownDevice.getEntities().size() > entityAddr && conn.ownDevice.getEntities().get(entityAddr).getFeatures().size() > featureAddr;

            // Count the number of subscriptions
            int numSubscriptions = conn.countSubscriptions(function.getSubscriptionRequest().getServerAddress());

            if (isValidRequest && numSubscriptions < conn.ownDevice.getEntities().get(entityAddr).getFeatures().get(featureAddr).getMaxSubscriptions()) {
                Logger.getGlobal().info("Accept subscription from: " + function.getSubscriptionRequest().getClientAddress().getDevice());

                BindSubscribeEntry newEntry = new BindSubscribeEntry(
                        function.getSubscriptionRequest().getClientAddress(),
                        function.getSubscriptionRequest().getServerAddress()
                );

                conn.bindSubscribeInfo.add(new BindSubscribeInfo("subscription", newEntry));
                conn.bindSubscribeNotify("subscription", conn, newEntry);

                FeatureAddressType serverAddr = function.getSubscriptionRequest().getServerAddress();
                conn.sendXML(
                        conn.ownDevice.makeHeader(serverAddr.getEntity(), serverAddr.getFeature(), function.getSubscriptionRequest().getClientAddress(), "result", conn.msgCounter, false),
                        resources.makePayload("resultData", new ResultData(0, "positive acknowledgement for subscription request")));

                FunctionType funct = conn.ownDevice.getEntities().get(serverAddr.getEntity()).getFeatures().get(serverAddr.getFeature()).getFunctions().get(0);
                conn.sendXML(
                        conn.ownDevice.makeHeader(serverAddr.getEntity(), serverAddr.getFeature(), function.getSubscriptionRequest().getClientAddress(), "notify", conn.msgCounter, false),
                        resources.makePayload(funct.getFunctionName(), funct.getFunction()));

                // Save binding
            } else {
                AddressType ownAddr = datagram.getHeader().getAddressDestination();
                conn.sendXML(
                        conn.ownDevice.makeHeader(ownAddr.getEntity(), ownAddr.getFeature(), datagram.getHeader().getAddressSource(), "result", conn.msgCounter, false),
                        resources.makePayload("resultData", new ResultData(1, "negative acknowledgement for subscription request")));
            }
        } catch (JAXBException e) {
            e.printStackTrace();
        }
    }
}
