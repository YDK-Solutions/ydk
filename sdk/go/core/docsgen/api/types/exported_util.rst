.. _types-utility:

Exported Utility Functions
==========================

.. go:package:: ydk/types
    :synopsis: YDK Types

This package contains YDK Go types. It provides built-in types specified in
`YANG RFC 6020 <https://tools.ietf.org/html/rfc6020>`_ and types used in YDK Go APIs.

.. attribute:: EntitySlice

    A slice of :go:struct:`Entity`

.. function:: (s EntitySlice) Len()

    :return: The length of given slice
    :rtype: ``int``

.. function:: (s EntitySlice) Less(i, j int)

    :return: If the Entity at index i is less than the one at index j
    :rtype: ``bool``

.. function:: (s EntitySlice) Swap(i, j int)

    Swaps the Entities at indices i and j

.. function:: GetRelativeEntityPath(current_node Entity, ancestor Entity, path string)

    :param current_node: (:go:struct:`Entity`)
    :param ancestor: (:go:struct:`Entity`)
    :param path: A Go ``string``
    :return: The relative entity path
    :rtype: A Go ``string``

.. function:: IsSet(Filter YFilter)

    :param Filter: :ref:`YFilter <y-filter>`
    :return: If the given yfilter is set (not equal to NotSet)
    :rtype: ``bool``

.. function:: EntityEqual(x, y Entity)

    :param x: (:go:struct:`Entity`)
    :param y: (:go:struct:`Entity`)
    :return: If the entities x and y and their children are equal in value
    :rtype: ``bool``
