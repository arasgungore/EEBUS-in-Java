import java.util.ArrayList;
import java.util.List;

public class Main {

    public static void main(String[] args) {
        // Replace placeholders with actual values
        String hostname = "yourHostname";
        boolean isGateway = true;  // or false
        String certName = "yourCertName";
        String devId = "yourDevId";
        String brand = "yourBrand";
        String devType = "yourDevType";

        EebusNode eebusNode = new EebusNode(hostname, isGateway, certName, devId, brand, devType);

        // Set the updater function
        eebusNode.update = (datagram, conn) -> {
            // Implement your logic for handling updates here
            System.out.println("Received update: " + datagram);
        };

        // Start the EebusNode
        eebusNode.start();
    }
}

class EebusNode {
    boolean isGateway;
    SpineNode spineNode;
    List<SpineConnection> connections;
    DeviceModel deviceStructure;
    Updater update;

    public EebusNode(String hostname, boolean isGateway, String certName, String devId, String brand, String devType) {
        this.isGateway = isGateway;
        this.spineNode = new SpineNode(hostname, isGateway, certName, devId, brand, devType);
        this.connections = new ArrayList<>();
        this.deviceStructure = new DeviceModel();
    }

    public void start() {
        spineNode.start();
    }

    public void subscriptionNotify(DatagramType datagram, SpineConnection conn) {
        if (update != null) {
            update.update(datagram, conn);
        }
    }
}
