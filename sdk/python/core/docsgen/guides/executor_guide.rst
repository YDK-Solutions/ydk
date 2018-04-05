How do I execute RPCs?
======================

.. contents:: Table of Contents

This document contains some examples of executing RPCs defined in yang. To perform these operations, the :py:class:`ExecutorService<ydk.services.ExecutorService>` is used.

The below approach can be used to execute a rollback RPC.

Executing a rollback RPC
------------------------

For this example, the :py:class:`Cisco_IOS_XR_cfgmgr_rollback_act.RollBackConfigurationLast<ydk.models.cisco_ios_xr.Cisco_IOS_XR_cfgmgr_rollback_act.RollBackConfigurationLast>` class is used. Note that the ``ydk`` and ``ydk-models-cisco-ios-xr`` python packages need to be installed for this example.


.. code-block:: python
    :linenos:

    # Import the rollback module
    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_cfgmgr_rollback_act

    # Import the executor service and netconf provider
    from ydk.services import ExecutorService
    from ydk.providers import NetconfServiceProvider

    # Create object
    roll_back_configuration_to = Cisco_IOS_XR_cfgmgr_rollback_act.RollBackConfigurationLast()

    # Force roll back for the five most recent changes
    roll_back_configuration_to.input.comment = "Forced programmatic rollback"
    roll_back_configuration_to.input.count = 5
    roll_back_configuration_to.input.force = True
    roll_back_configuration_to.input.label = "PRB-005"

    # Create executor service
    executor = ExecutorService()

    # Create a NetconfServiceProvider instance to connect to the device
    provider = NetconfServiceProvider(address='10.0.0.1',
                                         port=830,
                                         username='test',
                                         password='test')

    # Execute RPC on NETCONF device
    executor.execute_rpc(provider, roll_back_configuration_to)
