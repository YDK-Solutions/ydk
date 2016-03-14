.. _presence-class:

Presence Classes
==================
Most classes exist only to contain child nodes. This is the default style.

Consider a class, Main, containing two classes, A and B:

.. code-block:: python

	def __init__(self):
	    self.A = Main.A()
	    self.B = Main.B()

A "presence class" is configuration data, representing a single bit of configuration data. The class acts as both a configuration knob and a means of organizing related configuration.

Consider the case where class A is a presence class, but B is not:

.. code-block:: python

	def __init__(self):
	    self.A = None


Here, A gets instantiated, but B does not. B must be manually instantiated.

For more information, `see RFC 6020 <https://tools.ietf.org/html/rfc6020#section-7.5.1>`_
