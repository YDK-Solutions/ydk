.. _ref-howto_validation:

How to Disable Validation
=========================

By default all the data sent to Yang enabled server are validated against Yang model constraints like data type, number range and string patterns.
This is one of the major features of YDK. However there is sometime a need to disable data validation in order to see how the server reacts to invalid data.
For this sole purpose the YDK has possibility to disable data validation. This feature currently implemented only for :cpp:class:`NetconfServiceProvider<ydk::NetconfServiceProvider>`
through :cpp:class:`CrudService<ydk::CrudService>` and :cpp:class:`NetconfService<NetconfService>`.

In order to disable validation of object/entity the :cpp:class:`Entity<ydk::Entity>` class member `ignore_validation` must be set to ``true``. 
Please note that this setting disables validation on the entire entity including leafs, and children entities. 
If the flag is set to non-top-level entity, the validation will be also disabled on all parent entities up to the top-level. 
Actually, it is important to note, when validation is disabled on non-top-level entity it must be a part of complete hierarchy of a top-level entity.

Here is an example of a negative test, which disables validation of an entity containing data of invalid type.

.. code-block:: c++

  void int8_invalid_ignore_validation"()
  {
    NetconfServiceProvider provider{"10.10.10.10", "admin", "admin"};
    CrudService crud{};

    //CREATE
    auto r = ydktest_sanity::Runner();
    r.ignore_validation = true;
    r.ytypes->built_in_t->number8 = "test";
    CHECK_THROWS_WITH(crud.create(provider, r), Catch::Contains("\"test\" is not a valid value"));
  }

In response to this request the Netconf Server returns an error message, which is then captured and processed.

.. code-block:: sh

    [2019-02-08 08:43:54.025] [ydk] [error] RPC error occurred: <?xml version="1.0" encoding="UTF-8"?>
    <rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" message-id="4">
      <rpc-error>
        <error-type>application</error-type>
        <error-tag>invalid-value</error-tag>
        <error-severity>error</error-severity>
        <error-path xmlns:ydkut="http://cisco.com/ns/yang/ydktest-sanity">
        /rpc/edit-config/config/ydkut:runner/ydkut:ytypes/ydkut:built-in-t/ydkut:number8
      </error-path>
        <error-message xml:lang="en">"test" is not a valid value.</error-message>
        <error-info>
          <bad-element>number8</bad-element>
        </error-info>
      </rpc-error>
    </rpc-reply>
