## Bundle Profiles

In order to generate model bundle, as it is described in the root directory README.md section `Generate and install model bundle`, it is necessary to create a bundle profile, which is JSON formatted file. The profile defines name, version, dependencies and YANG files included into the bundle. The profile must contain bundle description and bundle components parts. It optionally can contain dependency part.

### Bundle Description

This part must include:
- `"name"` of the bundle; the name also used to form part of the installation path of the bundle
- `"version"` of the bundle
- `"ydk_version"`, which refers to the [version](https://github.com/CiscoDevNet/ydk-gen/releases) of the YDK core package.

Other components of description part are optional. Here is simple example of description part:

```
{
    "name":"cisco-ios-xr",
    "version": "6.5.3",
    "ydk_version": "0.8.4",
    "Author": "Cisco",
    "Copyright": "Cisco",
    "Description": "Cisco IOS-XR Native Models From Git",
```

### Bundle Components

The `"models"` element of the profile describes sources of YANG models. It can contain combination of elements:
- `"dir"` - list of **relative** directory paths containing YANG files
- `"file"` - list of **relative** YANG file paths
- `"git"` - git repository, where YANG files are located

##### Example of use of local directory

In this example the bundle includes all YANG files located in local directory

```
    "models": {
        "dir": [
            "/yang/openconfig-nis"
        ]
    }
```

##### Example of use of specific files

In this example the bundle includes only selected YANG files

```
    "models": {
        "file": [
            "yang/ietf/ietf-inet-types@2013-07-15.yang",
            "yang/ietf/ietf-yang-types@2013-07-15.yang",
            "yang/ietf/iana-if-type.yang",
            "yang/ietf/ietf-netconf-acm@2012-02-22.yang",
            "yang/ietf/ietf-netconf-with-defaults@2011-06-01.yang",
            "yang/ietf/ietf-netconf@2011-06-01.yang",
            "yang/ietf/ietf-system.yang"
        ]
    }
```

##### Example of use of git sources

Each `"git"` source must specify `"url"` - git repository URL, and `"commits"` list. The specified URL must allow the repository
to be cloned without user intervention. Each element in `"commits"` list can specify:

- `"commitid"` - optional specification of a commit ID in string format. If not specified the HEAD revision is assumed. 
The further specified directories and files will be copied from the context of this commit.
- `"dir"` - optional list of **relative** directory paths within the git repository.
All `*.yang` files in specified directory **and any sub-directories** will be pulled into the generated bundle.
- `"file"` - optional list of **relative** `*.yang` file paths within the git repository.

```
    "models": {
        "git": [
            {
                "url": "https://github.com/YangModels/yang.git",
                "commits": [
                  {
                    "dir": [
                        "vendor/cisco/xr/653"
                    ]
                  }
                ]
            },
            {
                "url": "https://github.com/YangModels/yang.git",
                "commits": [
                  {
                    "commitid": "f6b4e2d59d4eedf31ae8b2fa3119468e4c38259c",
                    "dir": [
                        "experimental/openconfig/bgp",
                        "experimental/openconfig/policy"
                    ]
                  }
                ]
            }
        ]
    },
```

### Bundle Dependency

Sometimes the bundle can have dependency on other bundles. In this case the profile should define `"dependency"` list, where each list element defines  `"name"`, `"version"` and `"uri"` - location of dependent bundle profile.

```
    "dependency": [
      {
        "name": "ietf",
        "version": "0.1.5-post2",
        "core_version": "0.8.4",
        "uri": "file://profiles/bundles/ietf_0_1_5_post2.json"
      }
    ]
}
```

