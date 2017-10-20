Capability
==========

.. go:package:: ydk/cgopath
    :synopsis: CGo Path API Capability

.. go:struct:: Capability(model, revision)

    An instance of Capability is defined by the module name and revision.

    :param model: (``CString``) Model name.
    :param revision: (``CString``) Model revision.

    .. note::

        If no revision is provided, use empty string:

        .. code-block:: go

            import "C"
            var cap1 C.Capability
            var cap2 C.Capability
            
            cap1 = C.CapabilityCreate(C.CString("openconfig-bgp"), C.CString(""))
            cap2 = C.CapabilityCreate(C.CString("openconfig-bgp"), C.CString("2015-10-09"))

            defer C.CapabilityFree(cap1)
            defer C.CapabilityFree(cap2)
