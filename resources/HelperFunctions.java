package resources;

import javax.xml.bind.JAXBContext;
import javax.xml.bind.JAXBException;
import javax.xml.bind.Marshaller;
import java.io.ByteArrayOutputStream;
import java.util.List;

public class HelperFunctions {

    public static boolean stringInSlice(String a, List<String> list) {
        for (String b : list) {
            if (b.equals(a)) {
                return true;
            }
        }
        return false;
    }

    public static void checkError(Exception err) {
        if (err != null) {
            System.out.println(err.getMessage());
            err.printStackTrace();
            System.exit(1);
        }
    }

    public static int max(int x, int y) {
        return Math.max(x, y);
    }

    public static String xmlToString(Object in) {
        try {
            JAXBContext context = JAXBContext.newInstance(in.getClass());
            Marshaller marshaller = context.createMarshaller();
            ByteArrayOutputStream baos = new ByteArrayOutputStream();
            marshaller.marshal(in, baos);
            return baos.toString();
        } catch (JAXBException e) {
            checkError(e);
            return "";
        }
    }

    public static String timestampNow() {
        return String.format("%tFT%<tT.0Z", System.currentTimeMillis());
    }

    public static List<NodeManagementDetailedDiscoveryEntityInformationType> makeEntities(DeviceModel device) {
        // Implement this method to create and return a list of NodeManagementDetailedDiscoveryEntityInformationType
        return List.of();
    }

    public static List<NodeManagementDetailedDiscoveryFeatureInformationType> makeFeatures(DeviceModel device) {
        // Implement this method to create and return a list of NodeManagementDetailedDiscoveryFeatureInformationType
        return List.of();
    }

    public static HeaderType makeHeader(DeviceModel device, int entity, int feature,
                                        FeatureAddressType addressDestination, String cmdClassifier,
                                        int msgCounter, boolean ackRequest) {
        return new HeaderType() {{
            specificationVersion = Datagram.SPECIFICATION_VERSION;
            addressSource = new FeatureAddressType() {{
                device = device.deviceAddress;
                entity = entity;
                feature = feature;
            }};
            addressDestination = addressDestination;
            msgCounter = msgCounter;
            cmdClassifier = cmdClassifier;
            timestamp = timestampNow();
            ackRequest = ackRequest;
        }};
    }

    public static PayloadType makePayload(String functionName, Object function) {
        return new PayloadType() {{
            cmd = new CmdType() {{
                functionName = functionName;
                function = xmlToString(function);
            }};
        }};
    }

    public static FeatureAddressType makeFeatureAddress(String device, int entity, int feature) {
        return new FeatureAddressType() {{
            this.device = device;
            this.entity = entity;
            this.feature = feature;
        }};
    }
}
