What are deviations?
====================

.. contents:: Table of Contents

Overview
---------

Not all devices faithfully support features defined in standard yang models. For a particular device, it could support only part of features or the feature it supported varies from the standard module. The YANG language defines the `deviation statement <https://tools.ietf.org/html/rfc6020#section-7.18.3>`_ to indicate such a thing.
For example, using a deviation statement in `cisco-xr-openconfig-telemetry-deviations.yang <https://github.com/YangModels/yang/blob/master/vendor/cisco/xr/631/cisco-xr-openconfig-telemetry-deviations.yang#L73>`_, the netconf server can indicate that the :go:struct:`openconfig_telemetry/TelemetrySystem_Subscriptions_Dynamic_Subscription<ydk/models.openconfig.openconfig_telemetry/TelemetrySystem_Subscriptions_Dynamic_Subscription>` container of the ``openconfig_telemetry`` model is not supported.

.. code-block:: c

    deviation /oc-telemetry:telemetry-system/oc-telemetry:subscriptions/oc-telemetry:dynamic/oc-telemetry:subscription {
        deviate not-supported;
    }


How to use deviation with YDK
------------------------------

When using YDK to program a device which has some unsupported features, YDK will raise error **before** sending payload to device. For instance, if the device advertises via a deviation that it does not support the ``subscription`` node as shown above, an error will be raised.

We can try to configure a dynamic subscription with the below app.

.. code-block:: c
    :linenos:

    package main

    import "github.com/CiscoDevNet/ydk-go/ydk/models/openconfig/openconfig_telemetry"

    func main() {
        telemetry := openconfig_telemetry.TelemetrySystem{}
        telemetry.Subscriptions = openconfig_telemetry.TelemetrySystem_Subscriptions
        telemetry.Subscriptions.Dynamic = openconfig_telemetry.TelemetrySystem_Subscriptions_Dynamic{}
        
        s := openconfig_telemetry.Subscriptions.Dynamic.Subscription{}
        s.SubscriptionId = 123
        s.State = openconfig_telemetry.Subscriptions.Dynamic.Subscription.State{}
        s.State.SubscriptionId = 123

        telemetry.Subscriptions.Dynamic.Subscription = append(telemetry.Subscriptions.Dynamic.Subscription, &s)

        // Call the CRUD create on the top-level telemetry object
        // (assuming you have already instantiated the service and provider)
        crud.Create(&provider, &telemetry)
    }

However the above app will fail and raise errors because of the deviation which is active on the ``subscription`` node.
