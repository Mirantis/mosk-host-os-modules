# irqbalance module

The irqbalance module is designed to allow the cloud operator to install and configure the `irqbalance` service
on cluster machines using the day-2 operations API.

> Note: This module is implemented and validated against the following Ansible versions provided by MCC for Ubuntu 20.04, 22.04, 24.04:
> Ansible core 2.12.10 and Ansible collection 5.10.0.
>
> To verify the Ansible version in a specific Cluster release, refer to
> [Container Cloud documentation: Release notes - Cluster releases](https://docs.mirantis.com/container-cloud/latest/release-notes/cluster-releases.html).
> Use the *Artifacts > System and MCR artifacts* section of the corresponding Cluster release. For example, for
> [17.3.0](https://docs.mirantis.com/container-cloud/latest/release-notes/cluster-releases/17-x/17-3-x/17-3-0/17-3-0-artifacts.html#system-and-mcr-artifacts).

# Default irqbalance configuration

The default configuration file `/etc/default/irqbalance` can contain the following settings, as defined in the
[irqbalance documentation](https://github.com/Irqbalance/irqbalance/blob/master/misc/irqbalance.env):

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

# Setting empty values for the irqbalance parameters

When the cloud operator defines values for the irqbalance module in the `HOC` object, those values overwrite particular parameters
in the `/etc/default/irqbalance` file. If the operator does not define a value, the corresponding parameter in the `/etc/default/irqbalance`
configuration file keeps its current value.

For example, if you define `values.args` in the `HOC` object, this value overwrites the `IRQBALANCE_ARGS` parameter in `/etc/default/irqbalance`.
Otherwise, the `IRQBALANCE_ARGS` value remains the same in the configuration file.

If you need to provide an empty `IRQBALANCE_ARGS` value, you can define `values.args: ""` (empty string) in the `HOC` object.
Other parameters defined in `/etc/default/irqbalance` follow the same logic.

# Version 1.1.0 (latest)

The module allows installing, configuring, and enabling or disabling the `irqbalance` service on cluster machines.

Since v1.0.0, the following changes apply to the irqbalance module:

* Added the `oneshot` parameter.
* Changed the method of setting empty values for the irqbalance parameters for better usability:

  * When a parameter is not defined in `values` of the `HOC` object, the corresponding value remains the same in the irqbalance configuration file.
  * When a parameter is set to `""` (empty string) in `values` of the `HOC` object , the corresponding value in the `irqbalance` configuration file
    is also set to `""` (empty string).

The module accepts the following parameters, all of them are optional:

- `enabled`: Enables the `irqbalance` service. Defaults to `true`.
- `banned_cpulist`: Defines the `IRQBALANCE_BANNED_CPULIST` value. Do not define it if you do not want to update the current `IRQBALANCE_BANNED_CPULIST` value
  in the `irqbalance` configuration file. Mutually exclusive with `banned_cpus`.
- `banned_cpus`: Defines the `IRQBALANCE_BANNED_CPUS` value. Do not define it if you do not want to update the current  `IRQBALANCE_BANNED_CPUS` value in the
  `irqbalance` configuration file. Mutually exclusive with `banned_cpulist`. `IRQBALANCE_BANNED_CPUS` is deprecated in irqbalance v1.8.0.
- `args`: Defines the `IRQBALANCE_ARGS` value. Do not define it if you do not want to update the current `IRQBALANCE_ARGS` value in the `irqbalance`
  configuration file.
- `oneshot`: Defines the `IRQBALANCE_ONESHOT` value. Do not define it if you do not want to update the current `IRQBALANCE_ONESHOT` value in the `irqbalance`
  configuration file. `IRQBALANCE_ONESHOT` is commented out when `oneshot` is set to `false`, because setting `IRQBALANCE_ONESHOT` to any value leads to
  enablement of this functionality.
- `policy_script`: Defines the name of the irqbalance policy script, which is bash-compatible.
- `policy_script_filepath`: Defines the full file path to store the irqbalance policy script that can be used with the `--policyscript=<filePath>` argument.
  Do not define it if you do not want to write the policy script.
- `update_apt_cache`: Enables the update of `apt-cache` before installing the `irqbalance` service. Defaults to `true`.

> Note: `IRQBALANCE_BANNED_CPUS` is deprecated in irqbalance v1.8.0, which is used in Ubuntu 22.04, and is being replaced with `IRQBALANCE_BANNED_CPULIST`.
> For details, see [Release notes for irqbalance v1.8.0](https://github.com/Irqbalance/irqbalance/releases/tag/v1.8.0).

> Note: When you configure the policy script, at least the following parameters must be set: `args`, `policy_script`, and `policy_script_filepath`.
> Otherwise, the corresponding error message will be displayed in the status of the `HostOSConfiguration` object.

> Note: If an error message in the status of the `HostOSConfiguration` object contains `schema validation failed`,
> verify whether the types of used parameters are correct and whether the used combination of parameters is allowed.

> Note: If you enable the service without setting `banned_cpulist`, `banned_cpus`, `oneshot`, or `args`, the corresponding values
> in `/etc/default/irqbalance` will remain as they were before applying the current `HOC` configuration.

# Version 1.0.0 (deprecated)

> Note: The module version 1.0.0 is obsolete and not recommended for usage in production environments.

The module allows installing, configuring, and enabling or disabling the `irqbalance` service on cluster machines.
The module accepts the following parameters, all of them are optional:

- `enabled`: Enable the `irqbalance` service. Defaults to `true`.
- `banned_cpulist`: The `IRQBALANCE_BANNED_CPULIST` value. Leave empty to not update the current `IRQBALANCE_BANNED_CPULIST` value
  in the `irqbalance` configuration file. Mutually exclusive with `banned_cpus`.
- `banned_cpus`: The `IRQBALANCE_BANNED_CPUS` value. Leave empty to not update the current `IRQBALANCE_BANNED_CPUS` value
  in the `irqbalance` configuration file. `IRQBALANCE_BANNED_CPUS` is deprecated in irqbalance v1.8.0. Mutually exclusive with `banned_cpulist`.
- `args`: The `IRQBALANCE_ARGS` value. Leave empty to not update the current `IRQBALANCE_ARGS` value in the `irqbalance` configuration file.
- `policy_script`: The irqbalance policy script, which is bash-compatible.
- `policy_script_filepath`: The full file path name to store the irqbalance policy script that can be used with the `--policyscript=<filepath>` argument.
  Leave empty to not write the policy script.
- `update_apt_cache`: Enables the update of `apt-cache` before installing the `irqbalance` service. Defaults to `true`.

> Caution: When you configure the policy script, at least three parameters must be set: `args`, `policy_script`, and `policy_script_filepath`.
> Otherwise, the corresponding error message will be displayed in the status of the `HostOSConfiguration` object.

> Note: If an error message in the status of the `HostOSConfiguration` object contains `schema validation failed`,
> verify whether the types of used parameters are correct and whether the used combination of parameters is allowed.

> Note: If you enable the service without setting `banned_cpulist`, `banned_cpus`, `oneshot`, or `args`, the corresponding values
> in `/etc/default/irqbalance` will remain as they were before applying the current `HostOSConfiguration` configuration.

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

As a result of this configuration, no parameters will be set or overridden in the `irqbalance` configuration file.

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

As a result of this configuration:
- `IRQBALANCE_BANNED_CPULIST` and `IRQBALANCE_ARGS` will be set or overridden
- `IRQBALANCE_BANNED_CPUS` will be removed from the `irqbalance` configuration file
- `IRQBALANCE_ONESHOT` will be set to `True`.

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

As a result of this configuration:
- `IRQBALANCE_ARGS` will be set or overridden in the `irqbalance` configuration file
- The contents of `policy_script` will be written to `/etc/default/irqbalance-numa.sh`
- The `irqbalance` service will use the provided policy script

Refer to https://manpages.ubuntu.com/manpages/jammy/man1/irqbalance.1.html for the policy script description.
In particular, refer to the `numa_node` variable used in the example.

# Troubleshooting on the target host

Use the following troubleshooting commands for irqbalance on a host:

Verify the service status:

```sudo systemctl status irqbalance```

Verify the configuration:

```less /etc/default/irqbalance```

Verify the `init.d` script

```less /etc/init.d/irqbalance```

Verify logs:

```journalctl -u irqbalance*```

Verify statistics of interrupts:

```less -S /proc/interrupts```

Verify connections of NICs to NUMA nodes:

```cat /sys/class/net/<nic_name>/device/numa_node```

> Note: `numa_node` exists for a given NIC only if NUMA is configured on the host.

# irqbalance documentation

For information on the `irqbalance` service, refer to the official
[irqbalance documentation for Ubuntu 22.04](https://manpages.ubuntu.com/manpages/jammy/man1/irqbalance.1.html) and the
[Upstream GitHub project](https://github.com/Irqbalance/irqbalance/).
