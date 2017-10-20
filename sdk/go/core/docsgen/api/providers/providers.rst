Serivce Providers
=================

.. go:package:: ydk/providers
    :synopsis: Service Providers API

.. contents:: Table of Contents

Open Daylight
-------------

.. go:struct:: OpendaylightServiceProvider

	OpenDaylightServiceProvider A service provider to be used to communicate with an OpenDaylight instance: https://www.opendaylight.org

	.. attribute:: Path

		A Go ``string`` that represents the path

	.. attribute:: Address

		A Go ``string`` that represents the IP address of the device supporting a open daylight interface

	.. attribute:: Username

		A Go ``string`` that represents the username to log in to the device

	.. attribute:: Password

		A Go ``string`` that represents the password to log in to the device
	
	.. attribute:: Port

		An ``int`` that represents the device port used to access the open daylight interface.

	.. attribute:: EncodingFormat

		An instance of :go:struct:`EncodingFormat<ydk/types/EncodingFormat>`

	.. attribute:: Protocol

		A Go ``string`` that represents protocol used to connect to the device

	.. attribute:: Private

		An instance of :go:struct:`COpenDaylightServiceProvider<ydk/types/COpenDaylightServiceProvider>`

	.. attribute:: ProvidersHolder

		A slice ([] :go:struct:`ServiceProvider<ydk/types/ServiceProvider>`)

	.. attribute:: State

		An instance of :go:struct:`State<ydk/types/State>` which represents the error state of this provider

.. function:: (provider *OpenDaylightServiceProvider) GetPrivate()

	:return: private pointer for OpenDaylightServiceProvider
	:rtype: interface{}

.. function:: (provider *OpenDaylightServiceProvider) Connect()

	Connect to OpenDaylightServiceProvider using Path/Address/Username/Password/Port

.. function:: (provider *OpenDaylightServiceProvider) GetNodeIDs()

	:return: OpenDaylightServiceProvider Node IDs
	:rtype: []string

.. function:: (provider *OpenDaylightServiceProvider) GetNodeProvider(nodeID string)

	:return: Node provider by ID
	:rtype: :go:struct:`ServiceProvider<ydk/types/ServiceProvider>`

.. function:: (provider *OpenDaylightServiceProvider) GetState()

	:return: error state from OpenDaylightServiceProvider
	:rtype: :go:struct:`*State<ydk/types/State>`

.. function:: (provider *OpenDaylightServiceProvider) Disconnect()

	Disconnect from OpenDaylightServiceProvider


Netconf
-------

.. go:struct:: NetconfServiceProvider

	NetconfServiceProvider Implementation of ServiceProvider for the NETCONF protocol: https://tools.ietf.org/html/rfc6241

	.. attribute:: Repo

		An instance of :go:struct:`Repository<ydk/types/Repository>` This attribute represents the repository of YANG models.

	.. attribute:: Address

		A Go ``string`` that represents the IP address of the device supporting a netconf interface

	.. attribute:: Username

		A Go ``string`` that represents the username to log in to the device

	.. attribute:: Password

		A Go ``string`` that represents the password to log in to the device

	.. attribute:: Port

		An ``int`` that represents the device port used to access the netconf interface.

	.. attribute:: Protocol

		A Go ``string`` that represents protocol used to connect to the device

	.. attribute:: Private

		An instance of :go:struct:`CServiceProvider<ydk/types/CServiceProvider>`

	.. attribute:: State

		An instance of :go:struct:`State<ydk/types/State>`

.. function:: (provider *NetconfServiceProvider) GetPrivate()

	:return: private pointer for NetconfServiceProvider
	:rtype: interface{}

.. function:: (provider *NetconfServiceProvider) Connect()
	
	Connect to NetconfServiceProvider using Repo/Address/Username/Password/Port

.. function:: (provider *NetconfServiceProvider) GetState()

	:return: error state from NetconfServiceProvider
	:rtype: :go:struct:`*State<ydk/types/State>`

.. function:: (provider *NetconfServiceProvider) Disconnect()

	Disconnect from NetconfServiceProvider


Restconf
--------

.. go:struct:: RestconfServiceProvider

	RestconfServiceProvider Implementation of ServiceProvider for the RESTCONF protocol: https://tools.ietf.org/html/draft-ietf-netconf-restconf-18

	.. attribute:: Path

		A Go ``string`` that represents the path

	.. attribute:: Address

		A Go ``string`` that represents the IP address of the device supporting a restconf interface

	.. attribute:: Username

		A Go ``string`` that represents the username to log in to the device

	.. attribute:: Password

		A Go ``string`` that represents the password to log in to the device

	.. attribute:: Port

		An ``int`` that represents the device port used to access the restconfs interface.

	.. attribute:: Encoding

		An instance of :go:struct:`EncodingFormat<ydk/types/EncodingFormat>`

	.. attribute:: StateURLRoot

		A Go ``string``. This attribute provides backwards compatibility with older drafts of restconf RFC, this can be "/operational" or "/data"

	.. attribute:: ConfigURLRoot

		A Go ``string``. This attribute provides backwards compatibility with older drafts of restconf RFC, this can be "/config" or "/data" (which is the default)

	.. attribute:: Private

		An instance of :go:struct:`CServiceProvider<ydk/types/CServiceProvider>`

	.. attribute:: State

		An instance of :go:struct:`State<ydk/types/State>`

.. function:: (provider *RestconfServiceProvider) GetPrivate()
	
	:return: private pointer for RestconfServiceProvider
	:rtype: interface{}

.. function:: (provider *RestconfServiceProvider) Connect()

	Connect to RestconfServiceProvider using Path/Address/Username/Password/Port

.. function:: (provider *RestconfServiceProvider) GetState()
	
	:return: error state from RestconfServiceProvider
	:rtype: :go:struct:`*State<ydk/types/State>`

.. function:: (provider *RestconfServiceProvider) Disconnect
	
	Disconnect from RestconfServiceProvider


Codec
-----

.. go:struct:: CodecServiceProvider

	CodecServiceProvider Encode and decode to XML/JSON format

	.. attribute:: Repo

		An instance of :go:struct:`Repository<ydk/types/Repository>` This attribute represents the repository of YANG models.

	.. attribute:: Encoding

		An instance of :go:struct:`EncodingFormat<ydk/types/EncodingFormat>`

	.. attribute:: RootSchemaTable

		An instance of map[string] :go:struct:`RootSchemaNode<ydk/types/RootSchemaNode>`

	.. attribute:: State

		An instance of :go:struct:`State<ydk/types/State>`

.. function:: (provider *CodecServiceProvider) Initialize()

	Initialize the CodecServiceProvider

.. function:: (provider *CodecServiceProvider) GetEncoding()

	:return: the encoding format for CodecServiceProvider
	:rtype: :go:struct:`EncodingFormat<ydk/types/EncodingFormat>`

.. function:: (provider *CodecServiceProvider) GetState()

	:return: error state from CodecServiceProvider
	:rtype: :go:struct:`*State<ydk/types/State>`

.. function:: (provider *CodecServiceProvider) GetRootSchemaNode(entity types.Entity)

	:return: root schema node for entity
	:rtype: :go:struct:`RootSchemaNode<ydk/types/RootSchemaNode>`

