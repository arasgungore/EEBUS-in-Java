package resources;

import java.util.List;

public class Discovery {

    public static class NodeManagementDetailedDiscovery {
        public List<NodeManagementSpecificationVersionListType> specificationVersionList;
        public NodeManagementDetailedDiscoveryDeviceInformationType deviceInformation;
        public List<NodeManagementDetailedDiscoveryEntityInformationType> entityInformation;
        public List<NodeManagementDetailedDiscoveryFeatureInformationType> featureInformation;
    }

    public static class NodeManagementSpecificationVersionListType {
        public String specificationVersion;
    }

    public static class NodeManagementDetailedDiscoveryDeviceInformationType {
        public NetworkManagementDeviceDescriptionDataType description;
    }

    public static class NodeManagementDetailedDiscoveryEntityInformationType {
        public NetworkManagementEntityDescritpionDataType description;
    }

    public static class NodeManagementDetailedDiscoveryFeatureInformationType {
        public NetworkManagementFeatureInformationType description;
    }

    public static class NetworkManagementDeviceDescriptionDataType {
        public DeviceAddressType deviceAddress;
        public String deviceType;
        public String description;
    }

    public static class NetworkManagementEntityDescritpionDataType {
        public EntityAddressType entityAddress;
        public String entityType;
        public String description;
    }

    public static class NetworkManagementFeatureInformationType {
        public FeatureAddressType featureAddress;
        public String featureType;
        public String role;
        public FunctionPropertyType supportedFunction;
        public String description;
    }

    public static class FeatureAddressType {
        public String device;
        public int entity;
        public int feature;
    }

    public static class EntityAddressType {
        public String device;
        public int entity;
    }

    public static class DeviceAddressType {
        public String device;
    }

    public static class FunctionPropertyType {
        public String function;
        public String possibleOperations;
    }
}
