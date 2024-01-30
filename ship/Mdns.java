package ship;

import io.github.zeroone3010.yahueapi.discovery.HueBridgeDiscoveryService;
import io.github.zeroone3010.yahueapi.discovery.HueBridgeDiscoveryServiceImpl;
import io.github.zeroone3010.yahueapi.discovery.HueBridgeDiscoveryListener;
import io.github.zeroone3010.yahueapi.discovery.HueBridgeDiscoveryEvent;
import io.github.zeroone3010.yahueapi.domain.HueBridge;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

public class Mdns {

    private ShipNode shipNode;

    public Mdns(ShipNode shipNode) {
        this.shipNode = shipNode;
    }

    public void browseDns() {
        System.out.println("Browsing for entries");
        HueBridgeDiscoveryService discoveryService = new HueBridgeDiscoveryServiceImpl();
        discoveryService.addDiscoveryListener(new HueBridgeDiscoveryListener() {
            @Override
            public void onDiscovery(HueBridgeDiscoveryEvent event) {
                handleFoundService(event.getBridge());
            }
        });

        discoveryService.startDiscovery();

        // Let it run for a while
        try {
            TimeUnit.SECONDS.sleep(60);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        discoveryService.stopDiscovery();
    }

    public void registerDns() {
        try {
            String port = Integer.toString(shipNode.getServerPort());
            String id = "DEVICE-EEB01-" + shipNode.getDevId() + ".local.";

            String[] txtRecord = {
                    "txtvers=1",
                    "id=" + id,
                    "path=wss://" + shipNode.getHostname() + ":" + port,
                    "SKI=" + shipNode.getSki(),
                    "register=true",
                    "brand=" + shipNode.getBrand(),
                    "type=" + shipNode.getDevType()
            };

            System.out.println("Registering: " + String.join(", ", txtRecord));

            HueBridgeDiscoveryService discoveryService = new HueBridgeDiscoveryServiceImpl();
            HueBridge bridge = new HueBridge(shipNode.getHostname(), port, shipNode.getDevId(), txtRecord);
            discoveryService.registerBridge(bridge);

            // Shutdown server after 2 minutes
            TimeUnit.SECONDS.sleep(120);

            discoveryService.unregisterBridge(bridge);
        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }
    }

    private void handleFoundService(HueBridge bridge) {
        // Implement the logic to handle the found service (Hue Bridge in this case)
    }
}
