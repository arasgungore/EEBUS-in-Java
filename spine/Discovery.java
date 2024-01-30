package spine;

import resources.*;

public class Discovery {

    public SpineConnection conn;
    public String address;
    public int msgCounter;
    public NodeManagementDetailedDiscovery discoveryInformation;

    public Discovery(SpineConnection conn) {
        // Constructor implementation
        this.conn = conn;
        this.address = "";
        this.msgCounter = 0;
        this.discoveryInformation = null;
    }

    public void startDetailedDiscovery() {
        FunctionModel funct = conn.ownDevice.getEntities().get(0).getFeatures().get(0).getFunctions().get(0);
        conn.sendXML(conn.ownDevice.makeHeader(0, 0, resources.makeFeatureAddress("", 0, 0), "read", msgCounter, false),
                resources.makePayload("nodeManagementDetailedDiscoveryData", funct));

        DatagramType answer = conn.recieveTimeout(10);
        if (answer != null) {
            address = answer.getHeader().getAddressSource().getDevice();
            NodeManagementDetailedDiscovery Function = null;
            try {
                Function = xmlMapper.readValue(answer.getPayload().getCmd().getFunction(), NodeManagementDetailedDiscovery.class);
            } catch (IOException e) {
                e.printStackTrace();
            }

            if (Function != null) {
                discoveryInformation = Function;
                for (EntityModel entity : conn.ownDevice.getEntities()) {
                    for (FeatureModel feature : entity.getFeatures()) {
                        for (FeatureInformation featureInformation : Function.getFeatureInformation()) {
                            String destDevice = Function.getDeviceInformation().getDescription().getDeviceAddress().getDevice();
                            int destEntity = featureInformation.getDescription().getFeatureAddress().getEntity();
                            int destFeature = featureInformation.getDescription().getFeatureAddress().getFeature();

                            FeatureAddressType ownAddr = new FeatureAddressType(
                                    conn.ownDevice.getDeviceAddress(),
                                    entity.getEntityAddress(),
                                    feature.getFeatureAddress()
                            );

                            int numBindings = conn.countBindings(ownAddr);
                            int numSubscriptions = conn.countSubscriptions(ownAddr);

                            if (feature.getBindingTo() != null &&
                                    resources.stringInList(featureInformation.getDescription().getFeatureType(), feature.getBindingTo()) &&
                                    ("server".equals(featureInformation.getDescription().getRole()) ||
                                            "special".equals(featureInformation.getDescription().getRole())) &&
                                    conn.ownDevice.getEntities().get(entity.getEntityAddress())
                                            .getFeatures().get(feature.getFeatureAddress()).getMaxBindings() > numBindings) {
                                conn.sendBindingRequest(entity.getEntityAddress(), feature.getFeatureAddress(),
                                        resources.makeFeatureAddress(destDevice, destEntity, destFeature),
                                        featureInformation.getDescription().getFeatureType());
                            }

                            if (feature.getSubscriptionTo() != null &&
                                    resources.stringInList(featureInformation.getDescription().getFeatureType(), feature.getSubscriptionTo()) &&
                                    ("server".equals(featureInformation.getDescription().getRole()) ||
                                            "special".equals(featureInformation.getDescription().getRole())) &&
                                    conn.ownDevice.getEntities().get(entity.getEntityAddress())
                                            .getFeatures().get(feature.getFeatureAddress()).getMaxSubscriptions() > numSubscriptions) {
                                conn.sendSubscriptionRequest(entity.getEntityAddress(), feature.getFeatureAddress(),
                                        resources.makeFeatureAddress(destDevice, destEntity, destFeature),
                                        featureInformation.getDescription().getFeatureType());
                            }
                        }
                    }
                }
            }
        }
    }
}
