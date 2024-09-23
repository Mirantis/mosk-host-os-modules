# cpushield day2 module

**_NOTE:_** Only Ubuntu 22.04 and newer is supported with cgroup v2 only and systemd (e.g. in Unified Control Group hierarchy mode, without cgroup v1)

This module allows to configure cpu and NUMA node shielding. The module points systemd itself (PID 1 under init.scope), system.slice, user.slice and other unit files to use cores/NUMA nodes specified in parameters. Thus, other cores may be used exclusively, f.e., for pinning virtual machines vCPUs there (Nova vcpu_pin_set setting).

[AllowedCPUs](https://www.freedesktop.org/software/systemd/man/latest/systemd.resource-control.html#AllowedCPUs=) and [AllowedMemoryNodes](https://www.freedesktop.org/software/systemd/man/latest/systemd.resource-control.html#AllowedMemoryNodes=) Systemd unit file options are used for this purpose.

Reboot is required to properly move all system tasks to specified cores.

## Supported cpushield parameters

- `system_cpus` - list of CPU cores that will be used for system processes, required;
- `system_mem_numas` - list of NUMA nodes that will be used for system processes, optional;
- `disable_reboot_request` - boolean, `true` or `false`. If `true`, module will NOT create a special file for LCM agent for requesting a subsequent reboot (e.g. reboot will not occur).
- `enable` - enable or disable shielding systemd service (boolean). Default: `true`;
- `systemd_units_to_pin` - list of systemd units for which *AllowedCPUs* option will be set. See recommended value in an example below;
- `disable_old_shield_service` - boolean, `true` or `false`. If `true`, module will disable `old_shield_service_name` systemd service. Default: `false`;
- `old_shield_service_name` - name of systemd service that implements CPU shielding in older Ubuntu releases < 22.04. Default: `shield-cpus.service`.

Useful links:
1. [Manual configuration for older Ubuntu versions - cgroup v1](https://docs.mirantis.com/mosk/latest/deploy/deploy-openstack/advanced-config/advanced-compute/configure-cpu-isolation.html?highlight=cpu%20isolation)
2. [Shielding Linux Resources Book](https://documentation.suse.com/sle-rt/15-SP5/pdf/book-shielding_en.pdf)

## Examples

To pin system processes onto CPU cores 0, 1, 2 and 5 for Ubuntu 22.04 (cgroup v2 case) use:
```
---
values:
  system_cpus: '0-2, 5'
  systemd_units_to_pin:
  - system.slice
  - user.slice
  - kubepods.slice
```

To pin system processes onto CPU cores 10, 11 and disable old shielding service after Ubuntu 20.04 -> 22.04 upgrade:
```
---
values:
  disable_old_shield_service: true
  old_shield_service_name: my-old-shield-service.service
  system_cpus: '10, 11'
  systemd_units_to_pin:
  - system.slice
  - user.slice
  - kubepods.slice
```
