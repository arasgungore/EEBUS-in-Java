package ship;

import io.github.zeroone3010.yahueapi.JsonStringUtil;
import io.github.zeroone3010.yahueapi.domain.MessageType;

import javax.json.JsonObject;
import javax.json.JsonReader;
import javax.json.JsonValue;
import java.io.StringReader;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class ShipNode {

    public int serverPort;
    public String hostname;
    public boolean isGateway;
    public List<SMEInstance> smeList;
    public List<Request> requests;
    public ConnectionManagerSpine spineConnectionNotify;
    public CloseHandler spineCloseHandler;
    public String certName;
    public String devId;
    public String brand;
    public String devType;

    public ShipNode(String hostname, boolean isGateway, String certName, String devId, String brand, String devType) {
        this.serverPort = 0;
        this.hostname = hostname;
        this.isGateway = isGateway;
        this.smeList = new ArrayList<>();
        this.requests = new ArrayList<>();
        this.spineConnectionNotify = null;
        this.spineCloseHandler = null;
        this.certName = certName;
        this.devId = devId;
        this.brand = brand;
        this.devType = devType;
    }

    public void start() {
        int port = FreePort.getFreePort();
        this.serverPort = port;

        // Start server, Register Dns, and search for other DNS entries
        if (!this.isGateway) {
            startServer();
            registerDns();
        }

        browseDns();

        if (this.isGateway) {
            ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);
            scheduler.scheduleAtFixedRate(this::makeRequests, 0, 3, TimeUnit.SECONDS);
        }
    }

    public void makeRequests() {
        for (Request request : requests) {
            skis, _ := readSkis();
            if resources.StringInSlice(request.getSki(), skis) {
                connect(request.getPath(), request.getSki());
            }
        }
    }

    public void handleFoundService(zeroconf.ServiceEntry entry) {
        // If found service is not on the same device
        if (entry.getPort() != serverPort) {
            String ski = entry.getText().get(3).split("=")[1];
            if (resources.StringInSlice(ski, skis)) {
                connect(entry.getText().get(2).split("=")[1], ski);
            } else if (isGateway) {
                boolean requestAlreadyMade = false;
                for (Request req : requests) {
                    if (req.getSki().equals(ski)) {
                        requestAlreadyMade = true;
                        break;
                    }
                }
                if (!requestAlreadyMade) {
                    requests.add(new Request(entry.getText().get(2).split("=")[1], entry.getText().get(6).split("=")[1] + " " +
                            entry.getText().get(5).split("=")[1] + " " + entry.getText().get(1).split("=")[1], ski));
                }
            }
        }
    }

    public void newConnection(String role, WebSocketConnection conn, String ski) {
        String skiIsNew = "";
        List<String> skis = readSkis();
        if (!resources.StringInSlice(ski, skis) && isGateway) {
            skiIsNew = ski;
        }

        SMEInstance newSME = new SMEInstance(role, "INIT", conn, spineCloseHandler, ski);

        for (SMEInstance e : smeList) {
            if (e.getSki().equals(ski)) {
                conn.close();
                return;
            }
        }

        smeList.add(newSME);
        newSME.startCMI();
        spineConnectionNotify.connectionNotify(newSME, skiIsNew);
    }

    public void connect(String service, String ski) {
        WebSocketConnection conn = // create WebSocketConnection using the specified service
        newConnection("client", conn, ski);
    }

    public void startServer() {
        // Start the server
        HttpServer server = new HttpServer();
        server.startServer();
    }
}

class FreePort {
    public static int getFreePort() {
        // Implement logic to get a free port
        return 0;
    }
}

class HttpServer {
    public void startServer() {
        // Implement logic to start the HTTP server
    }
}
