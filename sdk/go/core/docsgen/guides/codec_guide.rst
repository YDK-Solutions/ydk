How do I work with instances of YANG data?
==========================================

.. contents:: Table of Contents

This document contains some examples of encoding and decoding yang data. To perform these operations, the :go:struct:`CodecService<ydk/services/CodecService>` is used.

The below approaches can be used to perform encoding and decoding of an interface Ipv4 loopback configuration. For these examples, the :go:struct:`InterfaceConfigurations<ydk/models/cisco_ios_xr/ifmgr_cfg/InterfaceConfigurations>` struct is used. Note that the ``ydk`` and ``cisco_ios_xr`` go packages need to be installed for this example.

Converting between JSON and XML
-------------------------------

To parse a JSON string representing yang data into a YDK go object and then to an XML string, the below approach can be used.

.. code-block:: go
    :linenos:

    package main

    import (
        "fmt"
        "github.com/CiscoDevNet/ydk-go/ydk/models/cisco_ios_xr/ifmgr_cfg"
        "github.com/CiscoDevNet/ydk-go/ydk/providers"
        "github.com/CiscoDevNet/ydk-go/ydk/services"
        encoding "github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format"
    )

    // Declare the JSON configuration
    const (
        ifJSON = `{
            "Cisco-IOS-XR-ifmgr-cfg:interface-configurations": {
                "interface-configuration": [
                {
                    "active": "act",
                    "interface-name": "Loopback0",
                    "description": "PRIMARY ROUTER LOOPBACK",
                    "Cisco-IOS-XR-ipv4-io-cfg:ipv4-network": {
                        "addresses": {
                            "primary": {
                                "address": "172.16.255.1",
                                "netmask": "255.255.255.255"
                            }
                        }
                    }
                }
            }
        }`
    )

    // execute main program.
    func main() {
        // Instantiate the codec service
        codec := services.CodecService{}

        // Instantiate codec providers with json and xml options
        jsonProvider := providers.CodecServiceProvider{}
        jsonProvider.Encoding = encoding.JSON
        xmlProvider := providers.CodecServiceProvider{}
        xmlProvider.Encoding = encoding.XML

        // Invoke the decode method to decode the JSON payload to a YDK go object
        interfaceConfigurations := codec.Decode(&jsonProvider, ifJSON)

        // Invoke the encode method to encode the YDK go object to an XML string
        ifXML = codec.Encode(&xmlProvider, &interfaceConfigurations)
        fmt.Println(ifXML)
    }


Converting to JSON
-------------------

To convert a YDK python object into a JSON string, the below approach can be used. Note that the attribute ``primary`` is an instance of a :ref:`presence type<presence-type>`, which is set to ``nil`` by default. So it needs to be assigned to a new instance of its type.

.. code-block:: go
    :linenos:

    package main

    import (
        "fmt"
        "github.com/CiscoDevNet/ydk-go/ydk/models/cisco_ios_xr/ifmgr_cfg"
        "github.com/CiscoDevNet/ydk-go/ydk/providers"
        "github.com/CiscoDevNet/ydk-go/ydk/services"
        encoding "github.com/CiscoDevNet/ydk-go/ydk/types/encoding_format"
    )

    // execute main program.
    func main() {
        // Instantiate the codec service
        codec := services.CodecService{}

        // Instantiate the provider with json option
        jsonProvider := providers.CodecServiceProvider{}
        jsonProvider.Encoding = encoding.JSON

        // Instantiate the InterfaceConfiguration instance
        interfaceConfig := ifmgr_cfg.InterfaceConfigurations_InterfaceConfiguration{}
        interfaceConfig.Active = "Act"
        interfaceConfig.InterfaceName = "Loopback0"
        interfaceConfig.Description = "PRIMARY ROUTER LOOPBACK"

        // Instantiate the Primary presence node
        interfaceConfig.Ipv4Network.Addresses.Primary.Address = "172.16.255.1"
        interfaceConfig.Ipv4Network.Addresses.Primary.Netmask = "255.255.255.255"

        // Instantiate the interface configuration structure
        interfaceConfigs := ifmgr_cfg.InterfaceConfigurations{}
        intefaceConfigs.InterfaceConfiguration = append(intefaceConfigs.InterfaceConfiguration, &interfaceConfig)
        
        // Invoke the encode method to encode the YDK go object to a JSON payload
        jsonPayload := codec.Encode(&jsonProvider, &interfaceConfigs)
        fmt.Println(jsonPayload)
    }
