How do I execute RPCs?
======================

.. contents:: Table of Contents

This document contains some examples of executing RPCs defined in yang. To perform these operations, the :go:struct:`ExecutorService<ydk/services/ExecutorService>` is used.

The below approach can be used to execute a rollback RPC.

Executing a rollback RPC
------------------------

For this example, the :go:struct:`RollBackConfigurationLast<ydk/models/cisco_ios_xr/cfgmgr_rollback_act/RollBackConfigurationLast>` struct is used. Note that the ``ydk`` and ``cisco_ios_xr`` go packages need to be installed for this example.


.. code-block:: c
    :linenos:

    package main

    // Import the rollback module, executor service, and netconf provider
    import (
        "fmt"
        "github.com/CiscoDevNet/ydk-go/ydk/models/cisco_ios_xr/cfgmgr_rollback_act"
        "github.com/CiscoDevNet/ydk-go/ydk/providers"
        "github.com/CiscoDevNet/ydk-go/ydk/services"
    )

    // execute main program
    func main() {
        // Create object
        rollBackConfigTo := cfgmgr_rollback_act.RollBackConfigurationLast{}

        // Force roll back for the five most recent changes
        rollBackConfigTo.Input.Comment = "Forced programmatic rollback"
        rollBackConfigTo.Input.Count = 5
        rollBackConfigTo.Input.Force = true
        rollBackConfigTo.Input.Label = "PRB-005"

        // Create the executor service
        executor := services.ExecutorService{}

        // Create a NetconfServiceProvider instance to connect to the device
        provider := providers.NetconfServiceProvider{
            Address: "10.0.0.1",
            Username: "test",
            Password: "test",
            Port: 830,
            Protocol: "ssh"}
        provider.Connect()

        // Execute RPC on NETCONF device
        executor.ExecuteRpc(&provider, &rollBackConfigTo)
    }
