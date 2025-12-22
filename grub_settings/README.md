# grub_settings module

The grub_settings module allows configuring most of GRUB2 variables, which are specified
in the `/etc/default/grub` configuration file by default.

By default, the module uses the drop-in `/etc/default/grub.d/99-grub_settings_hoc_module.cfg`
file to save the user-defined configuration. Drop-in configuration files take precedence
over `/etc/default/grub`. The user can override the drop-in file name using the `grub_cfg_filename`
parameter. The lexicographical order is used for file processing. For example, `100-myconf.cfg` is
processed before `90-foobar.conf`.

The module allows overriding any GRUB2 option. For example, if `/etc/default/grub` contains
`GRUB_CMD_LINE="loglevel=4 iommu=off"` but `/etc/default/grub.d/99-grub_settings_hoc_module.cfg` has
`GRUB_CMD_LINE="iommu=on debug"`, only the latter value is added to kernel parameters.

> Note: The grub_settings module creates a special file for LCM agent to request a subsequent reboot.
> This file has the text format and contains a line with the reboot reason. LCM agent reports
> to LCM controller that reboot is required for the corresponding LCM machine. You can disable
> creation of a reboot request by setting `disable_reboot_request` to `true`.
>
> To perform a reboot, create a
> [GracefulRebootRequest](https://docs.mirantis.com/mosk/latest/api/mgmt-api/lcm-api/graceful-reboot-request.html)
> object with a specific machine name.

## Supported GRUB2 parameters (under the `options` key)

- `grub_timeout` (integer) - defines `GRUB_TIMEOUT` in seconds.
- `grub_timeout_style` (string) - defines `GRUB_TIMEOUT_STYLE`. Allowed values: `menu`,
`countdown`, and `hidden`.
- `grub_hidden_timeout` (integer) - defines `GRUB_HIDDEN_TIMEOUT` in seconds.
- `grub_hidden_timeout_quiet` (boolean) - defines `GRUB_HIDDEN_TIMEOUT_QUIET`.
- `grub_recordfail_timeout` (integer) - defines `GRUB_RECORDFAIL_TIMEOUT` in seconds.
- `grub_default` (string) - defines `GRUB_DEFAULT` for the default kernel to boot.
- `grub_savedefault` (boolean) - defines `GRUB_SAVEDEFAULT`.
- `grub_cmdline_linux` (array of strings) - contains the list of options form the `GRUB_CMDLINE_LINUX`
parameter value. All of them are joined into a string. An empty list is allowed to override `GRUB_CMDLINE_LINUX`
with an empty value.
- `grub_cmdline_linux_default` (array of strings) - contains the list of options form the
`GRUB_CMDLINE_LINUX_DEFAULT` parameter value. All of them are joined into a string. An empty list
is allowed to override `GRUB_CMDLINE_LINUX` with an empty value.
- `grub_disable_os_prober` (boolean) - defines `GRUB_DISABLE_OS_PROBER`.
- `grub_disable_linux_recovery` (boolean) - defines `GRUB_DISABLE_LINUX_RECOVERY`.
- `grub_gfxmode` (string) - defines `GRUB_GFXMODE`. Only values in the
`1280x1024x16,800x600x24,640x480` format are allowed. For example, a screen resolution must be
divided by `,`.

Following parameters are available since 1.1.0 version:
- `grub_disable_recovery` (boolean) - defines `GRUB_DISABLE_RECOVERY`.
- `grub_preload_modules` (array of strings) - defines `GRUB_PRELOAD_MODULES`.

> Learn more:
>
> - [Ubuntu GRUB2 Settings](https://help.ubuntu.com/community/Grub2/Setup)
> - [Official GRUB2 docs](https://www.gnu.org/software/grub/manual/grub/html_node/Simple-configuration.html)


## Special module parameters

- `grub_cfg_filename` (string, optional) - name of a custom file under `/etc/default/grub.d`.
- `grub_reset_to_defaults` (boolean) - removes the drop-in configuration file with settings added
by the module and regenerates `grub.cfg`. Only the `true` value is allowed. Mutually exclusive with
all GRUB2 parameters specified above.
- `disable_reboot_request` (boolean) - creation of a special file for LCM agent to request
a subsequent reboot. If `true`, module does not create such a file and reboot does not occur. Default: `false`.

## Examples

Change some GRUB2 options without reboot:

```
---
values:
  options:
    grub_timeout: 10
    grub_hidden_timeout: 20
    grub_hidden_timeout_quiet: false
    grub_recordfail_timeout: 234
    grub_default: '1>2'
    grub_savedefault: false
    grub_cmdline_linux_default:
      - quiet
    grub_cmdline_linux:
      - 'cgroup_enable=memory'
      - 'debug'
      - 'intel_iommu=off'
    grub_disable_os_prober: true
    grub_disable_linux_recovery: false
    grub_gfxmode: '640x480x16'
    grub_timeout_style: menu
```

Use a custom grub configuration file name:

```
---
values:
  grub_cfg_filename: 9999-my-config.cfg
  options:
    grub_disable_os_prober: false
    grub_default: '3'
```

Change some GRUB2 options without a reboot request:

```
---
values:
  disable_reboot_request: true
  options:
    grub_timeout: 360
    grub_hidden_timeout_quiet: true
```

Reset the GRUB2 configuration to default without a reboot request:

```
---
values:
  grub_reset_to_defaults: true
  disable_reboot_request: true
```

Reset the GRUB2 configuration to default with a reboot request:

```
---
values:
  grub_reset_to_defaults: true
```

> Warning: You cannot use `grub_reset_to_defaults` with any GRUB2 configuration option,
> it is restricted by `schema.json` of the module. For example, the following incorrect
> configuration results in the `schema.json` validation failure in the `HostOSConfiguration`
> resource:
>
>```
>---
>values:
>  grub_reset_to_defaults: true
>  options:
>    grub_cmdline_linux:
>      - 'cgroup_enable=memory'
>      - 'debug'
>      - 'intel_iommu=off'
>```

