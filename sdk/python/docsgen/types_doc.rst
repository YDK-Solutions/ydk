Using Types
***********
This document will explain and give examples to using the ydk types.
Types explained will include:

- Empty.
- Decimal64.
- FixedBitsDict.
- YList.

Example use of Empty type
=========================

- The leaf being configured (accept_route) under the module ydk.models.routing.routing_policy:
**accept_route:** accepts the route into the routing table. **type:** Empty

.. code-block:: python

    from ydk.models.routing.routing_policy import RoutingPolicy

    policy_definitions = RoutingPolicy.PolicyDefinitions()
    accept_leaf        = policy_definitions.policy_definition.statements.statement.actions.policy_definitions.policy_definition.statements.statement.actions.accept_route()

Example use of Decimal64 type
=============================

- The leaf being configured (restart_timer) under the ydk.models.bgp.bgp module:
**restart_timer:** Time interval in seconds after which the BGP session is re-established after being torn down due to exceeding the max-prefix limit. **type:** Decimal64

.. code-block:: python

    from ydk.models.bgp.bgp import Bgp

    neighbors   = Bgp.Neighbors()
    accept_leaf = neighbors.neighbor.afi_safis.afi_safi.ipv4_labelled_unicast.prefix_limit.config.neighbors.neighbor.afi_safis.afi_safi.ipv4_labelled_unicast.prefix_limit.config.restart_timer()

Example use of FixedBitsDict type
=================================

- The leaf being configured (restart_timer) under the ydk.models.ietf.ietf_netconf_acm module:
**access_operations:** Access operations associated with this rule.  This leaf matches if it has the value '*' or if the bit corresponding to the requested operation is set. **type:** str

.. code-block:: python

    from  from ydk.models.ietf.ietf_netconf_acm import Nacm

    rule_list   = Nacm.RuleList()
    rule        = rule_list.Rule()
    rule.parent = rule_list
    rule.rule_list.rule.access_operations['read'] = True

Example use of YList type
=========================

- The node being configured (rule_list) under the ydk.models.ietf.ietf_netconf_acm module:
**rule_list:** An ordered collection of access control rules. **type:** RuleList

.. code-block:: python

    from  from ydk.models.ietf.ietf_netconf_acm import Nacm

    rule_list   = Nacm.RuleList()
    rule        = rule_list.Rule()
    rule.parent = rule_list
