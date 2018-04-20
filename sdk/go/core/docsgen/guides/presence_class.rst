.. _presence-type:

What are presence types?
==========================

According to `RFC 6020 <https://tools.ietf.org/html/rfc6020#section-7.5.1>`_, YANG supports two styles of containers. One is for organizing hierarchy. Another type (called 'presence container') is for representing configuration data. For instance the existence of presence container ``ssh`` indicates the capability to log in to the device using ssh. Let us consider a presence node ``prefix-list-entries`` in `Cisco-IOS-XR-ipv4-acl-cfg.yang<https://github.com/YangModels/yang/blob/96883adbf612605f02271523d7eaa731ded46b61/vendor/cisco/xr/621/Cisco-IOS-XR-ipv4-acl-cfg.yang#L105>`_. This node is generated as a YDK struct shown below:

.. code-block:: c
    :linenos:

    package ipv4_acl_cfg

    // Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries
    // Sequence of entries forming a prefix list
    // This type is a presence type.
    type Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries struct {
        EntityData types.CommonEntityData
        YFilter yfilter.YFilter

        // A prefix list entry; either a description (remark) or a prefix to match
        // against. The type is slice of
        // Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries_PrefixListEntry.
        PrefixListEntry []Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries_PrefixListEntry
    }

    func (prefixListEntries *Ipv4AclAndPrefixList_Prefixes_Prefix_PrefixListEntries) GetEntityData() *types.CommonEntityData {
        prefixListEntries.EntityData.YFilter = prefixListEntries.YFilter
        prefixListEntries.EntityData.YangName = "prefix-list-entries"
        prefixListEntries.EntityData.BundleName = "cisco_ios_xr"
        prefixListEntries.EntityData.ParentYangName = "prefix"
        prefixListEntries.EntityData.SegmentPath = "prefix-list-entries"
        prefixListEntries.EntityData.CapabilitiesTable = cisco_ios_xr.GetCapabilities()
        prefixListEntries.EntityData.NamespaceTable = cisco_ios_xr.GetNamespaces()
        prefixListEntries.EntityData.BundleYangModelsLocation = cisco_ios_xr.GetModelsPath()

        prefixListEntries.EntityData.Children = make(map[string]types.YChild)
        prefixListEntries.EntityData.Children["prefix-list-entry"] = types.YChild{"PrefixListEntry", nil}
        for i := range prefixListEntries.PrefixListEntry {
            prefixListEntries.EntityData.Children[types.GetSegmentPath(&prefixListEntries.PrefixListEntry[i])] = types.YChild{"PrefixListEntry", &prefixListEntries.PrefixListEntry[i]}
        }
        prefixListEntries.EntityData.Leafs = make(map[string]types.YLeaf)
        return &(prefixListEntries.EntityData)
    }
