# EEBUS-in-Java

![EEBUS Logo](https://github.com/arasgungore/EEBUS-in-Java/blob/main/assets/eebus_logo.png)

Java implementation of the EEBUS protocol suite which provides support for the SHIP and SPINE protocols. The EEBUS protocols facilitate efficient communication among smart home IoT devices, allowing seamless integration regardless of device brand or type. This open-source framework aims to offer a future-proof solution applicable to various environments and applications.



## Table of Contents

- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [License](#license)
- [Acknowledgments](#acknowledgments)
- [Author](#author)



## Prerequisites

Before you begin, ensure you have the following requirements installed on your machine:

- [Java](https://www.java.com/)
- [Maven](https://maven.apache.org/) for building and managing dependencies



## Getting Started

To start using the project, follow these steps:

1. Build the project using Maven:

   ```bash
   mvn clean install
   ```

2. Run the main driver code or utilize the provided classes for your application.



## Usage

### General Usage

The framework can be utilized as follows:

```java
// Importing the framework
import eebus "github.com/LMF-DHBW/go-eebus"

// Configure EEBUS node
eebusNode = eebus.NewEebusNode("100.90.1.102", true, "gateway", "0001", "DHBW", "Gateway");

// IP-Adr, is gateway, ssl cert name, device ID, brand name, device type
eebusNode.Update = update; // set method called on subscription updates

// Function that creates device structure
buildDeviceModel(eebusNode);

// Start node
eebusNode.Start();
```


### Creating Device Structure

The device structure for an EEBUS node can be created as follows:

```java
eebusNode.DeviceStructure.DeviceType = "Generic";
eebusNode.DeviceStructure.DeviceAddress = "Switch1";
eebusNode.DeviceStructure.Entities = Arrays.asList(
    new resources.EntityModel(
        "Switch",
        0,
        Arrays.asList(
            eebusNode.DeviceStructure.createNodeManagement(false),
            new resources.FeatureModel(
                "ActuatorSwitch",
                1,
                "client",
                resources.ActuatorSwitch("button", "button for leds"),
                Arrays.asList("ActuatorSwitch")
            )
        )
    )
);

// Create node management again, in order to update discovery data
eebusNode.DeviceStructure.Entities.get(0).getFeatures().set(0, eebusNode.DeviceStructure.createNodeManagement(false));
```


### Accepting Requests from Gateway

Requests from the gateway can be accepted as follows:

```java
// Requests are saved in the following list
eebusNode.SpineNode.ShipNode.Requests;

int i = 0; // Select the first entry as an example

// Accept request by connecting with the device
Request req = eebusNode.SpineNode.ShipNode.Requests.get(i);
new Thread(() -> eebusNode.SpineNode.ShipNode.connect(req.getPath(), req.getSki())).start();

// Remove request from list
eebusNode.SpineNode.ShipNode.Requests.remove(i);
```


### Reading Subscription Messages

```java
public void update(resources.DatagramType data, spine.SpineConnection conn) {
    int entitySource = data.getHeader().getAddressSource().getEntity();
    int featureSource = data.getHeader().getAddressSource().getFeature();

    if (conn.DiscoveryInformation.FeatureInformation.get(featureSource).getDescription().getFeatureType().equals("Measurement")) {
        resources.MeasurementDataType Function = JAXB.unmarshal(new StringReader(data.getPayload().getCmd().getFunction()), resources.MeasurementDataType.class);
        // Function.getValue() contains measured value
    }
}
```


### Reading Feature States

```java
for (spine.SpineConnection e : eebusNode.SpineNode.Connections) {
    for (resources.FeatureInformation feature : e.SubscriptionData) {
        if ("ActuatorSwitch".equals(feature.getFeatureType()) && "LED".equals(feature.getEntityType())) {
            resources.FunctionElement state = JAXB.unmarshal(new StringReader(feature.getCurrentState()), resources.FunctionElement.class);
            // state.getFunction() contains the state (on/off)
        }
    }
}
```


### Sending Subscription Messages

```java
for (spine.SpineConnection e : eebusNode.SpineNode.Subscriptions) {
    e.send("notify", resources.makePayload("actuatorSwitchData",
            new resources.FunctionElement("on")));
}
```


### Sending Binding Messages

```java
for (spine.SpineConnection e : eebusNode.SpineNode.Bindings) {
    e.send("write", resources.makePayload("actuatorSwitchData",
            new resources.FunctionElement("toggle")));
}
```

Feel free to adjust these examples to fit your specific use case.



## License

This project is licensed under the [MIT License](LICENSE).



## Acknowledgments

This Java implementation is inspired by the Go-based EEBUS communication framework [go_eebus](https://github.com/LMF-DHBW/go_eebus).



## Author

👤 **Aras Güngöre**

- LinkedIn: [@arasgungore](https://www.linkedin.com/in/arasgungore)
- GitHub: [@arasgungore](https://github.com/arasgungore)
