package ydktest

// const YdkOpenconfigModelsPath = "/usr/local/share/openconfig@0.1.1"
const YdkYdktestModelsPath = "/usr/local/share/ydk/127.0.0.1:12022"

func YdktestAugmentLookupTables() map[string]string {
	caps := make(map[string]string)
	// caps["openconfig-bgp"] = "2015-10-09"
	// caps["openconfig-interfaces"] = "2015-11-20"
	// caps["openconfig-routing-policy"] = "2015-10-09"

	caps["ietf-netconf-with-defaults"] = "2011-06-01"
	caps["main"] = "2015-11-17"
	caps["main-aug1"] = "2015-11-17"
	caps["oc-pattern"] = "2015-11-17"
	caps["ydktest-filterread"] = "2015-11-17"
	caps["ydktest-sanity"] = "2015-11-17"
	caps["ydktest-sanity-augm"] = "2015-11-17"
	caps["ydktest-sanity-types"] = "2016-04-11"
	caps["ydktest-types"] = "2016-05-23"
	caps["openconfig-bgp"] = "2016-06-21"
	caps["openconfig-bgp-policy"] = "2016-06-21"
	caps["openconfig-bgp-types"] = "2016-06-21"
	caps["openconfig-interfaces"] = "2016-05-26"
	caps["openconfig-extensions"] = "2015-10-09"
	caps["openconfig-types"] = "2016-05-31"
	caps["openconfig-policy-types"] = "2016-05-12"
	caps["openconfig-routing-policy"] = "2016-05-12"
	caps["iana-crypt-hash"] = "2014-04-04"
	caps["iana-if-type"] = "2014-05-08"
	caps["ietf-inet-types"] = "2013-07-15"
	caps["ietf-interfaces"] = "2014-05-08"
	caps["ietf-netconf-acm"] = "2012-02-22"
	caps["ietf-netconf-monitoring"] = "2010-10-04"
	caps["ietf-netconf-notifications"] = "2012-02-06"
	caps["ietf-yang-library"] = "2016-04-09"
	caps["ietf-yang-types"] = "2013-07-15"
	caps["ydk"] = "2016-02-26"
	caps["ietf-netconf"] = "2011-06-01"

	return caps
}
