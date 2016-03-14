ydk.services module
===================

.. py:currentmodule:: ydk.services

services.py

The Services module.

Supported Services Include

	CRUDService: Provider Create/Read/Update/Delete API's
	
.. py:class:: ydk.services.CRUDService

	Bases: :class:`ydk.services.Service`
	
	CRUD Service class for supporting CRUD operations on entities.
	
	.. py:method:: create(provider, entity)
		
		Create the entity
		
		:param provider: An instance of ydk.providers.ServiceProvider
		:param entity: An instance of an entity class defined under the ydk.models package or subpackages.
		
		:return: None

		:raises ydk.errors.YPYDataValidationError: if validation error has occurred
		:raises ydk.errors.YPYError: if other error has occurred

		Possible Errors:

		* a server side error
		* there isn't enough information in the entity to prepare the message (eg. missing keys)


	.. py:method:: read(provider, read_filter, only_config=False)
		
		Read the entity or entities.
		
		:param provider: An instance of ydk.providers.ServiceProvider
		
		:param read_filter: A read_filter is an instance of an entity class. An entity class is a class defined under the ydk.models package that is not an Enum, Identity or a subclass of FixedBitsDict). Attributes of this entity class may contain values that act as match expressions or can be explicitly marked as to be read by assigning an instance of type `ydk.types.READ` to them.
		
		:param only_config: Flag that indicates that only the data that represents configuration data is to be fetched. Default is set to False i.e both oper and config data will be fetched.
		
		
		:return: The entity or list of entities as identified by the `read_filter`
		
		:raises ydk.errors.YPYDataValidationError: if validation error has occurred
		:raises ydk.errors.YPYError: if other error has occurred

		Possible errors could be

		* a server side error
		* if there isn't enough information in the entity to prepare the message (missing keys for example)
		 
		
	.. py:method:: update(provider, entity)
		
		Update the entity.
		
		Note:

		* An attribute of an entity class can be deleted by setting to an instance of `ydk.types.DELETE`.
		* An entity can only be updated if it exists on the server. Otherwise a `ydk.errors.YPYError` will be raised.
		
		:param provider: An instance of ydk.providers.ServiceProvider
		:param entity: An instance of an entity class defined under the ydk.models package or subpackages.
		
		:return: None
		
		:raises ydk.errors.YPYDataValidationError: if validation error has occurred
		:raises ydk.errors.YPYError: if other error has occurred

		Possible errors could be

		* a server side error
		* if there isn't enough information in the entity to the message (missing keys for example)
		          	  

	.. py:method:: delete(provider, entity)
	
		Delete the entity
		
		:param provider: An instance of ydk.providers.ServiceProvider
		:param entity: An instance of an entity class defined under the ydk.models package or subpackages.
		
		:return: None
		
		:raises ydk.errors.YPYDataValidationError: if validation error has occurred 
		:raises ydk.errors.YPYError: if other error has occurred

		Possible errors could be

		* a server side error
		* if there isn't enough information in the entity to the message (missing keys for example)
		