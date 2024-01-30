package resources;

import java.util.List;

public class BindSubscribe {

    public static class NodeManagementBindingData {
        public List<BindSubscribeEntry> bindingEntries;
    }

    public static class BindSubscribeEntry {
        public FeatureAddressType clientAddress;
        public FeatureAddressType serverAddress;
    }

    public static class NodeManagementSubscriptionData {
        public List<BindSubscribeEntry> subscriptionEntries;
    }

    public static class NodeManagementBindingRequestCall {
        public BindingManagementRequestCallType bindingRequest;
    }

    public static class BindingManagementRequestCallType {
        public FeatureAddressType clientAddress;
        public FeatureAddressType serverAddress;
        public String serverFeatureType;
    }

    public static class NodeManagementSubscriptionRequestCall {
        public SubscriptionManagementRequestCallType subscriptionRequest;
    }

    public static class SubscriptionManagementRequestCallType {
        public FeatureAddressType clientAddress;
        public FeatureAddressType serverAddress;
        public String serverFeatureType;
    }

    public static class ResultElement {
        public int errorNumber;
        public String description;
    }

    public static ResultElement resultData(int errorNumber, String description) {
        return new ResultElement() {{
            this.errorNumber = errorNumber;
            this.description = description;
        }};
    }

    public static class FeatureAddressType {
        // Define FeatureAddressType properties as needed
    }
}
