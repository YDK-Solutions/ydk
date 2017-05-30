package openconfig

const YdkOpenconfigModelsPath = "/usr/local/share/openconfig@0.1.1"

func OpenconfigAugmentLookupTables() map[string]string {
    caps := make(map[string]string)
    caps["openconfig-bgp"] = "2015-10-09"
    caps["openconfig-interfaces"] = "2015-11-20"
    caps["openconfig-routing-policy"] = "2015-10-09"

    return caps
}
