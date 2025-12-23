# cpushield module

> Warning: Ubuntu 22.04 and 24.04 are supported with cgroup v2 and systemd. For example, this
> configuration supports Unified Control Group hierarchy mode, which is not supported by cgroup v1.

The cpushield day-2 module allows configuring CPU and NUMA node shielding.
The module points `systemd` itself (PID 1 under `init.scope`), `system.slice`, `user.slice`,
and other unit files to use cores/NUMA nodes specified in parameters. Thus, you can use
other cores exclusively, for example, for pinning vCPUs of virtual machines (using the Nova
`vcpu_pin_set` parameter). For this purpose, use the
[AllowedCPUs](https://www.freedesktop.org/software/systemd/man/latest/systemd.resource-control.html#AllowedCPUs=) and [AllowedMemoryNodes](https://www.freedesktop.org/software/systemd/man/latest/systemd.resource-control.html#AllowedMemoryNodes=) systemd unit file options.

> Note: Assigning all system tasks to the specified cores requires a reboot.

> Tip: Mirantis does not recommend creating more than one `HostOSConfiguration` (HOC) object per machine
> with the cpushield module, because the systemd drop-in configuration file name is hardcoded to
> `99-shielding.conf`. CPU cores/NUMA nodes of a newer HOC object will completely regenerate this
> configuration file.
>
> To change cpushield settings, for example, to use other cores for system processes,
> either edit an existing HOC object, or remove the old one and create a new one from scratch
> to avoid confusion by multiple cpushield-containing objects.

## Supported cpushield parameters

> Note: The cpushield module creates a special file for LCM agent to request a subsequent reboot.
> This file has the text format and contains a line with the reboot reason. LCM agent reports
> to LCM controller that reboot is required for the corresponding LCM machine. You can disable
> creation of reboot request by setting `disable_reboot_request` to `true`.
>
> To perform the reboot, create a [GracefulRebootRequest](https://docs.mirantis.com/mosk/latest/api/mgmt-api/lcm-api/graceful-reboot-request.html)
> object with a specific machine name.

- `system_cpus` (required) - list of CPU cores to use for system processes.
- `system_mem_numas` (optional) - list of NUMA nodes to use for system processes.
- `disable_reboot_request` (optional, bool) - creation of a special file for LCM agent to request
  a subsequent reboot for the machine. If `true`, module does not create such a file.
  Default: `false`.
- `systemd_units_to_pin` (optional) - list of systemd units for which the `AllowedCPUs` option
  will be set. For the recommended value, see the example below.
- `disable_old_shield_service` (optional, bool) - disablement of the `old_shield_service_name`
  systemd service. Default: `false`.
- `apply_settings_immediately` (optional, boolean) - enables `systemctl daemon-reload` to immediately
  apply CPU/NUMA pinning settings for units from `systemd_units_to_pin`. Only processes and threads
  spawned after setting this option to `true` will adhere to new pinning. A system reboot is required
  to pin all existing system processes to the specified cores. Otherwise, the option does not apply to
  currently running processes. Enablement of this option may cause side effects such as container restarts
  or re-creations. Default: `false`.
- `old_shield_service_name` (optional, string) - name of the systemd service that implements
  CPU shielding in old Ubuntu releases < 22.04. Default: `shield-cpus.service`.

> See also:
> - [Manual configuration for older Ubuntu versions - cgroup v1](https://docs.mirantis.com/mosk/latest/deploy/deploy-openstack/advanced-config/advanced-compute/configure-cpu-isolation.html?highlight=cpu%20isolation)
> - [Shielding Linux Resources Book](https://documentation.suse.com/sle-rt/15-SP5/pdf/book-shielding_en.pdf)

## Examples

To pin system processes onto CPU cores 0, 1, 2, and 5 for Ubuntu 22.04 or 24.04 in the cgroup v2 use case:

```
---
values:
  system_cpus: '0-2, 5'
  systemd_units_to_pin:
  - system.slice
  - user.slice
  - kubepods.slice
```

To pin system processes onto CPU cores 10 and 11 as well as disable the old shielding service
after the Ubuntu 20.04 -> 22.04 upgrade:

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
