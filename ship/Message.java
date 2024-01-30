package ship;

import com.google.gson.Gson;
import com.google.gson.JsonSyntaxException;

import java.io.IOException;
import java.io.StringReader;
import java.io.StringWriter;

import javax.json.Json;
import javax.json.JsonObject;
import javax.json.JsonReader;
import javax.json.JsonWriter;

import io.github.zeroone3010.yahueapi.JsonStringUtil;
import io.github.zeroone3010.yahueapi.domain.MessageType;

public class Message {

    public static final int CMI_TIMEOUT = 10;

    public static class SMEInstance {
        public String role;
        public String connectionState;
        public WebSocketConnection connection;
        public CloseHandler closeHandler;
        public String ski;

        public SMEInstance(String role, WebSocketConnection connection, CloseHandler closeHandler, String ski) {
            this.role = role;
            this.connectionState = "";
            this.connection = connection;
            this.closeHandler = closeHandler;
            this.ski = ski;
        }

        public byte[] receiveTimeout(int seconds) {
            // Implement the logic for receiving messages with a timeout
            return null;
        }

        public void startCMI() {
            // Implement the logic for starting CMI
        }

        public void receiveOnce(Handler handleFunc) {
            // Implement the logic for receiving a single message
        }

        public void receive(DataHandler handleFunc) {
            // Implement the logic for receiving messages continuously
        }

        public void send(resources.DatagramType payload) {
            // Implement the logic for sending messages
        }
    }

    public static class CmiMessage {
        public int messageType;
        public int messageValue;

        public CmiMessage(int messageType, int messageValue) {
            this.messageType = messageType;
            this.messageValue = messageValue;
        }
    }

    public static class DataValue {
        public HeaderType header;
        public resources.DatagramType payload;

        public DataValue(HeaderType header, resources.DatagramType payload) {
            this.header = header;
            this.payload = payload;
        }
    }

    public static class HeaderType {
        public String protocollId;

        public HeaderType(String protocollId) {
            this.protocollId = protocollId;
        }
    }

    public static void main(String[] args) {
        // Test the conversion to and from JSON
        Message message = new Message();
        String jsonString = message.toJson();
        System.out.println("JSON: " + jsonString);

        Message parsedMessage = Message.fromJson(jsonString);
        System.out.println("MessageType: " + parsedMessage.getMessageType());
        System.out.println("Payload: " + parsedMessage.getDataValue().getPayload());
    }

    // Attributes as public fields for simplicity
    public int messageType;
    public DataValue dataValue;

    public Message() {
        this.messageType = 0;
        this.dataValue = null;
    }

    public Message(int messageType, DataValue dataValue) {
        this.messageType = messageType;
        this.dataValue = dataValue;
    }

    // Getter and setter methods for attributes

    public int getMessageType() {
        return messageType;
    }

    public void setMessageType(int messageType) {
        this.messageType = messageType;
    }

    public DataValue getDataValue() {
        return dataValue;
    }

    public void setDataValue(DataValue dataValue) {
        this.dataValue = dataValue;
    }

    // JSON serialization methods

    public String toJson() {
        Gson gson = new Gson();
        return gson.toJson(this);
    }

    public static Message fromJson(String json) {
        Gson gson = new Gson();
        return gson.fromJson(json, Message.class);
    }
}
