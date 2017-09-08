Serivces
========

.. go:package:: ydk/services
    :synopsis: Services API

.. contents:: Table of Contents

CRUD
----

.. go:struct:: CrudService

	Supports CRUD operations on entities.

.. function:: (provider *CrudService) Create(provider ServiceProvider, entity Entity)

	Create the entity.

	:param provider: An instance :go:struct:`ServiceProvider`
	:param entity: An instance of :go:struct:`Entity<ydk/types/Entity>`
	:return: whether the operation was successful or not
	:rtype: ``bool``

.. function:: (provider *CrudService) Update(provider ServiceProvider, entity Entity)

	Update the entity.

	:param provider: An instance :go:struct:`ServiceProvider`
	:param entity: An instance of :go:struct:`Entity<ydk/types/Entity>`
	:return: whether the operation was successful or not
	:rtype: ``bool``

.. function:: (provider *CrudService) Delete(provider ServiceProvider, entity Entity)

	Delete the entity.

	:param provider: An instance :go:struct:`ServiceProvider`
	:param entity: An instance of :go:struct:`Entity<ydk/types/Entity>`
	:return: whether the operation was successful or not
	:rtype: ``bool``

.. function:: (provider *CrudService) Read(provider ServiceProvider, filter Entity)

	Read the entity.

	:param provider: An instance :go:struct:`ServiceProvider`
	:param entity: An instance of :go:struct:`Entity<ydk/types/Entity>`
	:return: the entity as identified by the given filter
	:rtype: :go:struct:`Entity<ydk/types/Entity>`

.. function:: (provider *CrudService) ReadConfig(provider ServiceProvider, filter Entity)

	Read only config.

	:param provider: An instance :go:struct:`ServiceProvider`
	:param entity: An instance of :go:struct:`Entity<ydk/types/Entity>`
	:return: the entity as identified by the given filter
	:rtype: :go:struct:`Entity<ydk/types/Entity>`


Codec
-----

.. go:struct:: CodecService

	Supports encoding and decoding Go model API objects of type :go:struct:`Entity`

.. function:: (provider *CodecService) Encode(provider CodecServiceProvider, entity Entity)

	Encode converts entity object to XML/JSON payload

	:param provider: An instance :go:struct:`CodecServiceProvider<ydk/types/CodecServiceProvider>`
	:param entity: An instance of :go:struct:`Entity<ydk/types/Entity>`
	:return: encoded payload
	:rtype: A Go ``string``

.. function:: (provider *CodecService) Decode(provider CodecServiceProvider, payload string)

	Decode converts XML/JSON object to entity object

	:param provider: An instance :go:struct:`CodecServiceProvider<ydk/types/CodecServiceProvider>`
	:param payload: A Go ``string`` representing an encoded payload to decode
	:return: the decoded entity object
	:rtype: :go:struct:`Entity<ydk/types/Entity>`
