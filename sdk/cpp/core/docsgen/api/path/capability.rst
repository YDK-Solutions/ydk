Capability
==========


.. cpp:class:: ydk::path::Capability

    An instance of Capability consists of module name, revision, a the set of enabled features defined in this modules as well as the list of deviation modules.

    .. cpp:member:: std::string module

        Capability's module name.

    .. cpp:member:: std::string revision

        Capability's module revision.

    .. cpp:member:: std::vector<std::string> features

        List of features defined in this module that are enabled.

    .. cpp:member:: std::vector<std::string> deviations

        List of deviations that targeted for this module.

    .. cpp:function:: Capability(const std::string& mod, const std::string& rev)

    .. cpp:function:: Capability(const std::string& mod, const std::string& rev,\
                             const std::vector<std::string>& f, const std::vector<std::string>& d)

    .. cpp:function:: Capability(const Capability& cap)

    .. cpp:function:: Capability(Capability&& cap)

    .. cpp:function:: Capability& operator=(const Capability& cap)

    .. cpp:function:: Capability& operator=(Capability&& cap)

    .. cpp:function:: bool operator==(const Capability& cap)
