.. _gnmi_service:

gNMI Service
============

.. cpp:class:: ydk::gNMIService

    Supports gNMI Set/Get/Subscribe operations on model API entities. It also allows to get gNMI server capabilities.

    .. cpp:function:: gNMIService()

        Constructs an instance of gNMIService.

    .. cpp:function:: bool set(gNMIServiceProvider & provider, Entity & entity)

        Create, update, or delete single entity in the server configuration.

        :param provider: (:cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>`) gNMI service provider instance.
        :param entity: (:cpp:class:`Entity<ydk::Entity>`) instance, which represents single container in device supported model.
                       The **Entity** instance must be annotated with :cpp:class:`YFilter<ydk::YFilter>`, which defines set operation:

                       * YFilter::replace - add new configuration or replace the whole configuration tree

                       * YFilter::update  - update or create configuration in existing tree

                       * YFilter::delete\_ - delete part or entire configuration tree

        :return: (``bool``) ``true`` if operation is successful, ``false`` otherwise.
        :raises: YServiceError if an error has occurred.

    .. cpp:function:: bool set(gNMIServiceProvider & provider, std::vector<Entity\*> entities)

        Create, update, or delete multiple entities in the server configuration.

        :param provider: (:cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>`) gNMI service provider instance.
        :param entities: (std::vector<ydk::Entity*>) multiple containers of the :cpp:class:`Entity<ydk::Entity>` instances encapsulated into ``std::vector<ydk::Entity*>``.
                       Each **Entity** instance must be annotated with :cpp:class:`YFilter<ydk::YFilter>`, which defines set operation:

                       * YFilter::replace - add new configuration or replace the whole configuration tree

                       * YFilter::update  - update or create configuration in existing tree

                       * YFilter::delete\_ - delete part or entire configuration tree

        :return: (``bool``) ``true``, if operation is successful, ``false`` otherwise.
        :raises: YServiceError if an error has occurred.

    .. cpp:function:: std::shared_ptr<Entity> get(gNMIServiceProvider & provider, Entity & read_filter, const std::string & read_mode)

        Read single entity from the server configuration.

        :param provider: (:cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>`) gNMI service provider instance.
        :param read_filter: (:cpp:class:`Entity<ydk::Entity>`) instance, which represents single container in device supported model.
        :param read_mode: (``std::string``) One of the values: ``CONFIG``, ``STATE``, ``OPERATIONAL``, or ``ALL``.
        :return: (std::shared_ptr<ydk::Entity>) An instance of :cpp:class:`Entity<ydk::Entity>` as identified by the **read_filter** or ``nullptr``, if operation fails.
        :raises: YServiceError if an error has occurred.

    .. cpp:function:: std::vector<std::shared_ptr<Entity>> get(gNMIServiceProvider & provider, std::vector<Entity\*> read_filters, const std::string & read_mode)

        Read multiple entities from the server configuration.

        :param provider: (:cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>`) gNMI service provider instance.
        :param read_filters: (std::vector<ydk::Entity\*>) multiple containers of the :cpp:class:`Entity<ydk::Entity>` instances encapsulated into ``std::vector<ydk::Entity*>``.
        :param read_mode: (``std::string``) One of the values: ``CONFIG``, ``STATE``, ``OPERATIONAL``, or ``ALL``.
        :return: The requested data encapsulated into **std::vector<std::shared_ptr<Entity>>** instance; if request fails - empty **std::vector**.
        :raises: YServiceError if an error has occurred.

    .. cpp:function:: subscribe(gNMIServiceProvider & provider, gNMISubscription & subscription, uint qos, std::string & mode, std::string & encoding, callback)

        Subscribe to telemetry updates.

        :param provider: (:cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>`) gNMI service provider instance.
        :param subscription: (:cpp:class:`gNMISubscription<ydk::gNMISubscription>`) An instance of structure, which represent the subscription.
        :param qos: (``uint``) QOS indicating the packet marking.
        :param mode: (``std::string``) Subscription mode: one of ``STREAM``, ``ONCE`` or ``POLL``.
        :param encoding: (``std::string``) Encoding method for the output: one of ``JSON``, ``BYTES``, ``PROTO``, ``ASCII``, or ``JSON_IETF``.
        :param callback: (``void*(const char*)``) Callback function, which is used to process the subscription data. The subscription data returned to the user as a string representation of protobuf ``SubscribeResponse`` message.
        :raises: YServiceError if an error has occurred.

    .. cpp:function:: subscribe(gNMIServiceProvider & provider, std::vector<ydk::gNMISubscription\*> & subscription, uint qos, std::string & mode, std::string & encoding, callback)

        Subscribe to telemetry updates.

        :param provider: (:cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>`) gNMI service provider instance.
        :param subscription: Set of (:cpp:class:`gNMISubscription<ydk::gNMISubscription>`) instances incapsulated into ``std::vector``, which represent the subscription.
        :param qos: (``uint``) QOS indicating the packet marking.
        :param mode: (``std::string``) Subscription mode: one of ``STREAM``, ``ONCE`` or ``POLL``.
        :param encoding: (``std::string``) Encoding method for the output: one of ``JSON``, ``BYTES``, ``PROTO``, ``ASCII``, or ``JSON_IETF``.
        :param callback: (``void*(const char*)``) Callback function, which is used to process the subscription data. The subscription data returned to the user as a string representation of protobuf ``SubscribeResponse`` message.
        :raises: YServiceError if an error has occurred.

    .. cpp:function:: std::string capabilities(ydk::gNMIServiceProvider & provider)
    
        Get gNMI server capabilities
        
        :param provider: (:cpp:class:`gNMIServiceProvider<ydk::gNMIServiceProvider>`) gNMI service provider instance.
        :return: (``std::string``) JSON encoded string, which represents gNMI server capabilities.
        :raises: YServiceError if an error has occurred.

.. cpp:class:: ydk::gNMISubscription

        Instance of this structure defines subscription for a single entity. Members of the structure are:
        
        * entity: (:cpp:class:`Entity<ydk::Entity>`) Instance of the subscription entity. This parameter must be set by the user.
        * subscription_mode: (``std::string``) Expected one of the following string values: ``TARGET_DEFINED``, ``ON_CHANGE``, or ``SAMPLE``; default value is ``ON_CHANGE``.
        * sample_interval: (``longlong``) Time interval in nanoseconds between samples in ``STREAM`` mode; default value is 60000000000 (1 minute).
        * suppress_redundant: (``bool``) Indicates whether values that not changed should be sent in a ``STREAM`` subscription; default value is ``false``
        * heartbeat_interval: (``longlong``) Specifies the maximum allowable silent period in nanoseconds when **suppress_redundant** is True. If not specified, the **heartbeat_interval** is set to 360000000000 (10 minutes) or **sample_interval** whatever is bigger.
