.. _presence-type:

What are presence types?
==========================

According to `RFC 6020 <https://tools.ietf.org/html/rfc6020#section-7.5.1>`_, YANG supports two styles of containers. One is for organizing hierarchy. Another type (called 'presence container') is for representing configuration data. For instance the existence of presence container ``ssh`` indicates the capability to log in to the device using ssh. Let us consider a presence node ``prefix-list-entries`` in `Cisco-IOS-XR-ipv4-acl-cfg.yang<https://github.com/YangModels/yang/blob/96883adbf612605f02271523d7eaa731ded46b61/vendor/cisco/xr/621/Cisco-IOS-XR-ipv4-acl-cfg.yang#L105>`_. This node is generated as a `YDK struct <https://github.com/CiscoDevNet/ydk-go/blob/a636555e554800cec197cfc0f4e51d345798cfe3/ydk/models/cisco_ios_xr/ipv4_acl_cfg/ipv4_acl_cfg.go#L901>`_ shown below:

.. code-block:: c
    :linenos:

    // Ipv4AclAndPrefixList_Prefixes_Prefix
    // Name of a prefix list
    type Ipv4AclAndPrefixList_Prefixes_Prefix struct {
        EntityData types.CommonEntityData
        YFilter yfilter.YFilter

        // This attribute is a key. Prefix list name - max 32 characters. The type is
        // string.
        PrefixListName interface{}

        // Sequence of entries forming a prefix list.
        PrefixListEntries Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries
    }

Since the existence of container ``prefix-list-entries`` itself represents configuration data, YDK does not instantiate an instance of class :go:struct:`Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries<ydk/models/cisco_ios_xr/ipv4_acl_cfg/Ipv4AclAndPrefixList/Prefixes/Prefix/PrefixListEntries>` and assign it to the ``PrefixListEntries`` leaf. The user needs to manually instantiate and assign this object as follows:

.. code-block:: c
    :linenos:
    
    pref := Ipv4AclAndPrefixList_Prefixes_Prefix{}
    pref.PrefixListEntries = Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries{}
