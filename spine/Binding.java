package spine;

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

public class Binding {

    public List<BindSubscribeInfo> bindSubscribeInfo;

    public void sendBindingRequest(int entityAddress, int featureAddress, FeatureAddressType destinationAddr, String featureType) {
        FeatureAddressType clientAddr = resources.MakeFeatureAddress(conn.getOwnDevice().getDeviceAddress(), entityAddress, featureAddress);
        conn.sendXML(
                conn.getOwnDevice().makeHeader(entityAddress, featureAddress, destinationAddr, "call", conn.getMsgCounter(), true),
                resources.makePayload("nodeManagementBindingRequestCall", new resources.NodeManagementBindingRequestCall(
                        new resources.BindingManagementRequestCallType(clientAddr, destinationAddr, featureType)
                ))
        );

        AnswerInfo answerInfo = conn.recieveTimeout(10);
        if (answerInfo.isOk()) {
            ResultElement function = JsonStringUtil.fromJson(answerInfo.getAnswer().getPayload().getCmd().getFunction(), ResultElement.class);
            if (function.getErrorNumber() == 0) {
                System.out.println("Binding to: " + destinationAddr.getDevice());

                BindSubscribeEntry newEntry = new BindSubscribeEntry(clientAddr, destinationAddr);
                conn.getBindSubscribeInfo().add(new BindSubscribeInfo("binding", newEntry));
                conn.bindSubscribeNotify("binding", conn, newEntry);
            }
        }
    }

    public void processBindingRequest(resources.DatagramType datagram) {
        NodeManagementBindingRequestCall function = JsonStringUtil.fromJson(datagram.getPayload().getCmd().getFunction(), NodeManagementBindingRequestCall.class);

        int entityAddr = function.getBindingRequest().getServerAddress().getEntity();
        int featureAddr = function.getBindingRequest().getServerAddress().getFeature();
        boolean isValidRequest = conn.getOwnDevice().getEntities().size() > entityAddr &&
                conn.getOwnDevice().getEntities().get(entityAddr).getFeatures().size() > featureAddr;

        int numBindings = conn.countBindings(function.getBindingRequest().getServerAddress());

        if (isValidRequest && conn.getOwnDevice().getEntities().get(entityAddr).getFeatures().get(featureAddr).getMaxBindings() > numBindings) {
            System.out.println("Binding to: " + function.getBindingRequest().getClientAddress().getDevice());

            BindSubscribeEntry newEntry = new BindSubscribeEntry(function.getBindingRequest().getClientAddress(), function.getBindingRequest().getServerAddress());
            conn.getBindSubscribeInfo().add(new BindSubscribeInfo("binding", newEntry));

            FeatureAddressType serverAddr = function.getBindingRequest().getServerAddress();
            conn.sendXML(
                    conn.getOwnDevice().makeHeader(serverAddr.getEntity(), serverAddr.getFeature(),
                            function.getBindingRequest().getClientAddress(), "result", conn.getMsgCounter(), false),
                    resources.makePayload("resultData", new ResultData(0, "positive acknowledgement for binding request"))
            );
        } else {
            FeatureAddressType ownAddr = datagram.getHeader().getAddressDestination();
            conn.sendXML(
                    conn.getOwnDevice().makeHeader(ownAddr.getEntity(), ownAddr.getFeature(),
                            datagram.getHeader().getAddressSource(), "result", conn.getMsgCounter(), false),
                    resources.makePayload("resultData", new ResultData(1, "negative acknowledgement for binding request"))
            );
        }
    }
}
