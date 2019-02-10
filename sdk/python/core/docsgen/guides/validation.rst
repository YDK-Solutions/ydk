.. _validation:

How to Disable Validation
=========================

By default all the data sent to Yang enabled server are validated against Yang model constraints like data type, number range and string patterns.
This is one of the major features of YDK. However there is sometime a need to disable data validation in order to see how the server reacts to invalid data.
For this sole purpose the YDK has possibility to disable data validation. This feature currently implemented only for :py:class:`NetconfServiceProvider<ydk.providers.NetconfServiceProvider>`
through :py:class:`CrudService<ydk.services.CRUDService>` and :py:class:`NetconfService<ydk.services.NetconfService>`.

In order to disable validation of object/entity the :py:class:`Entity<ydk.types.Entity>` class member `ignore_validation` must be set to ``True``. 
Please note that this setting disables validation on the entire entity including leafs, and children entities. 
If the flag is set to non-top-level entity, the validation will be also disabled on all parent entities up to the top-level. 
Actually, it is important to note, when validation is disabled on non-top-level entity it must be a part of complete hierarchy of a top-level entity.

Here is an example of a negative test, which demonstrates how YDK validation could be disabled.

.. code-block:: python

    from __future__ import absolute_import
    from __future__ import print_function
    
    from ydk.errors    import YModelError, YServiceProviderError
    from ydk.providers import NetconfServiceProvider
    from ydk.services  import CRUDService
    
    from ydk.models.ydktest.ydktest_sanity import Runner
    
    def run_test(provider, uint8_value, disable_validation):
    
        crud = CRUDService()
    
        vstr = 'enabled'
        if disable_validation:
            vstr = 'disabled'
        print("\nSetting uint8 number to %s when validation is %s" % (uint8_value, vstr))

        runner = Runner()
        runner.ignore_validation = disable_validation   # Set 'ignore_validation' flag
        runner.ytypes.built_in_t.number8 = uint8_value
        try:
            crud.create(provider, runner)
            print("OK")
        except YServiceProviderError as err:
            print("NETCONF FAILED with Error:")
            print(err.message.split('</error-message>')[0].split('"en">')[1])
        except YModelError as err:
            print("YDK VALIDATION FAILED with YModelError:")
            print(err.message)
    
    if __name__ == '__main__':
    
        provider = NetconfServiceProvider( "127.0.0.1", "admin", "admin", 12022)
    
        run_test(provider, 88, False)
    
        run_test(provider, 888, False)
    
        run_test(provider, 888, True)   # Disable YDK validation

The script produces the following results::

    Setting uint8 number to 88 when validation is enabled
    OK
    
    Setting uint8 number to 888 when validation is enabled
    YDK VALIDATION FAILED with YModelError:
     Invalid value "888" in "number8" element. Path: /ydktest-sanity:runner/ytypes/built-in-t/number8
    
    Setting uint8 number to 888 when validation is disabled
    NETCONF FAILED with Error:
    "888" is not a valid value.
