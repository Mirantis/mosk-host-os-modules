# irqbalance documentation

Refer to official irqbalance documentation for Ubuntu 22.04:
https://manpages.ubuntu.com/manpages/jammy/man1/irqbalance.1.html

Upstream project homepage:
https://github.com/Irqbalance/irqbalance/

# irqbalance configuration

Default configuration file `/etc/default/irqbalance` may look as the following:

```
# irqbalance is a daemon process that distributes interrupts across
# CPUs on SMP systems.  The default is to rebalance once every 10
# seconds.  This is the environment file that is specified to systemd via the
# EnvironmentFile key in the service unit file (or via whatever method the init
# system you're using has).

#
# IRQBALANCE_ONESHOT
#    After starting, wait for a minute, then look at the interrupt
#    load and balance it once; after balancing exit and do not change
#    it again.
#
#IRQBALANCE_ONESHOT=

#
# IRQBALANCE_BANNED_CPUS
#    64 bit bitmask which allows you to indicate which CPUs should
#    be skipped when reblancing IRQs.  CPU numbers which have their
#    corresponding bits set to one in this mask will not have any
#    IRQs assigned to them on rebalance.
#
#IRQBALANCE_BANNED_CPUS=

#
# IRQBALANCE_BANNED_CPULIST
#    The CPUs list which allows you to indicate which CPUs should
#    be skipped when reblancing IRQs. CPU numbers in CPUs list will
#    not have any IRQs assigned to them on rebalance.
#
#      The format of CPUs list is:
#        <cpu number>,...,<cpu number>
#      or a range:
#        <cpu number>-<cpu number>
#      or a mixture:
#        <cpu number>,...,<cpu number>-<cpu number>
#
#IRQBALANCE_BANNED_CPULIST=

#
# IRQBALANCE_ARGS
#    Append any args here to the irqbalance daemon as documented in the man
#    page.
#
#IRQBALANCE_ARGS=
```

Using the module, you can provide the following parameters, all of them are optional:

- `enabled`: Enable irqbalance service. 'true' by default.
- `banned_cpulist`: IRQBALANCE_BANNED_CPULIST value. Don't define it to not update current IRQBALANCE_BANNED_CPULIST in the irqbalance config file. Mutually exclusive with 'banned_cpus'.
- `banned_cpus`: IRQBALANCE_BANNED_CPUS value. Don't define it to not update current IRQBALANCE_BANNED_CPUS in the irqbalance config file. IRQBALANCE_BANNED_CPUS is deprecated in irqbalance v1.8.0. Mutually exclusive with 'banned_cpulist'.
- `args`: IRQBALANCE_ARGS value. Don't define it to not update current IRQBALANCE_ARGS in the irqbalance config file.
- `oneshot`: IRQBALANCE_ONESHOT value. Don't define it to not update current IRQBALANCE_ONESHOT in the irqbalance config file.
- `policy_script`: irqbalance policy script (bash compatible script).
- `policy_script_filepath`: Full file path name to store irqbalance policy script that can be used with '--policyscript=<filepath>' argument. Leave empty to not write policy script.
- `update_apt_cache`: Update apt cache before installing irqbalance. 'true' by default.

> Note. IRQBALANCE_BANNED_CPUS is deprecated in irqbalance v1.8.0 (that is used in Ubuntu 22.04), and is being replaced with IRQBALANCE_BANNED_CPULIST.
> For details, see https://github.com/Irqbalance/irqbalance/releases/tag/v1.8.0.

> Note. When you configure policy script, at least three parameters must be set: `args`, `policy_script` and `policy_script_filepath`.
> Otherwise, error message will be set in HOC object status.

> Note. If error message in HOC object status contains "schema validation failed", check:
> - types of used parameters;
> - whether used combination of parameters is allowed.

> Note. If you enable the service without setting `banned_cpulist`, `banned_cpus`, `oneshot` or `args`, the corresponding values
> in `/etc/default/irqbalance` will remain as they were before applying the current HOC configuration.

> Note. `IRQBALANCE_ONESHOT` is commented out when `oneshot` is set to `false`. It's because setting `IRQBALANCE_ONESHOT`
> to any value leads to enablement of this functionality.

# Configuration examples

## Example 1. Run irqbalance using defaults.

```
    spec:
      ...
      configs:
        ...
        - description: Example irqbalance configuration
          module: irqbalance
          moduleVersion: 1.1.0
          order: 1
          phase: "reconfigure"
          values: {}
```

In result, no parameters will be set/overridden in the irqbalance configuration file.

## Example 2. Run irqbalance and deny using certain CPU cores for IRQ balancing.

```
    spec:
      ...
      configs:
        ...
        - description: Example irqbalance configuration
          module: irqbalance
          moduleVersion: 1.1.0
          order: 1
          phase: "reconfigure"
          values:
            banned_cpulist: "0-15,31"
            oneshot: true
            args: "--journal"
```

In result, IRQBALANCE_BANNED_CPULIST and IRQBALANCE_ARGS parameters will be set/overridden,
IRQBALANCE_BANNED_CPUS parameter will be removed in the irqbalance configuration file,
IRQBALANCE_ONESHOT will be set to `True`.

## Example 3. Run irqbalance using policy script.

```
    spec:
      ...
      configs:
        ...
        - description: Example irqbalance configuration
          module: irqbalance
          moduleVersion: 1.1.0
          order: 1
          phase: "reconfigure"
          values:
            args: "--policyscript=/etc/default/irqbalance-numa.sh"
            policy_script: |
              #!/bin/bash

              # specifying  a -1 here forces irqbalance to consider an interrupt from a
              # device to be equidistant from all NUMA nodes.
              echo 'numa_node=-1'
            policy_script_filepath: "/etc/default/irqbalance-numa.sh"
```

In result,
- IRQBALANCE_ARGS parameter will be set/overridden in the irqbalance configuration file,
- contents of `policy_script` will be written to `/etc/default/irqbalance-numa.sh` file,
- irqbalance will use the provided policy script.

Refer to https://manpages.ubuntu.com/manpages/jammy/man1/irqbalance.1.html for policy script
description, in particular, `numa_node` variable used in the example.

# Troubleshooting on the target host

Check the service status:

```sudo systemctl status irqbalance```

Check the configuration:

```less /etc/default/irqbalance```

Check init.d script

```less /etc/init.d/irqbalance```

Check logs:

```journalctl -u irqbalance*```

Check interrupts statistics:

```less -S /proc/interrupts```

Check connections of NICs to NUMA nodes:

```cat /sys/class/net/<nic_name>/device/numa_node```

> Note. `numa_node` exists for a given NIC only if NUMA is configured on the host.
