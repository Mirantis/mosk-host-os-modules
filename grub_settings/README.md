# grub_settings day2 module

This module allows to configure most of Grub2 variables which are normally specified inside `/etc/default/grub` configuration file.

Module uses drop-in file `/etc/default/grub.d/99-grub_settings_hoc_module.cfg` by default to save a user-defined configuration. Drop-in config files take a precedence over `/etc/default/grub`. User can override drop-in filename via `grub_cfg_filename` parameter (see below). Note, that lexicographical order is used, e.g. file `100-myconf.cfg` will be processed **BEFORE** `90-foobar.conf`.

It allows to override any Grub2 option. F.e., if `/etc/default/grub` contains `GRUB_CMD_LINE="loglevel=4 iommu=off"`, but `/etc/default/grub.d/99-grub_settings_hoc_module.cfg` has `GRUB_CMD_LINE="iommu=on debug"`, only the latter value will be added to kernel parameters (e.g. no values joining for a variable used).

## Supported Grub2 parameters (under options key):
- `grub_timeout` - GRUB_TIMEOUT parameter, in seconds (integer);
- `grub_timeout_style` - GRUB_TIMEOUT_STYLE parameter. Allowed values: `menu`, `countdown`, `hidden`;
- `grub_hidden_timeout` - GRUB_HIDDEN_TIMEOUT parameter, in seconds (integer);
- `grub_hidden_timeout_quiet` - GRUB_HIDDEN_TIMEOUT_QUIET parameter, boolean, `true` or `false`;
- `grub_recordfail_timeout` - GRUB_RECORDFAIL_TIMEOUT parameter, in seconds (integer);
- `grub_default` - GRUB_DEFAULT parameter (default kernel to boot, string);
- `grub_savedefault` - GRUB_SAVEDEFAULT parameter, boolean, `true` or `false`;
- `grub_cmdline_linux` - a list of options to form GRUB_CMDLINE_LINUX parameter value. All of them are joined into a string. Empty list is allowed to override GRUB_CMDLINE_LINUX with empty value;
- `grub_cmdline_linux_default` - a list of options to form GRUB_CMDLINE_LINUX_DEFAULT parameter value. All of them are joined into a string. Empty list is allowed to override GRUB_CMDLINE_LINUX with empty value;
- `grub_disable_os_prober` - GRUB_DISABLE_OS_PROBER parameter, boolean, `true` or `false`;
- `grub_disable_linux_recovery` - GRUB_DISABLE_LINUX_RECOVERY parameter, boolean, `true` or `false`;
- `grub_gfxmode` - GRUB_GFXMODE parameter, string, only values in format `1280x1024x16,800x600x24,640x480` are allowed (e.g. screen resolutions divided by `,`).

Detailed information:
1. [Ubuntu Grub2 Settings](https://help.ubuntu.com/community/Grub2/Setup)
2. [Official Grub2 docs](https://www.gnu.org/software/grub/manual/grub/html_node/Simple-configuration.html)


## Special module parameters
- `grub_cfg_filename` - a name of custom file under `/etc/default/grub.d`, optional
- `grub_reset_to_defaults` - boolean, only `true` value is allowed. Mutually exclusive with all Grub2 settings parameters, specified above. Removes drop-in config file with settings added by module and regenerates `grub.cfg`;
- `disable_reboot_request` - boolean, `true` or `false`. If `true`, module will NOT create a special file for LCM agent for requesting a subsequent reboot.

## Examples

Change some Grub2 options without reboot:
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

Use a custom grub config filename:
```
---
values:
  grub_cfg_filename: 9999-my-config.cfg
  options:
    grub_disable_os_prober: false
    grub_default: '3'
```

Change some Grub2 options without reboot request:
```
---
values:
  disable_reboot_request: true
  options:
    grub_timeout: 360
    grub_hidden_timeout_quiet: true
```

Reset to default Grub2 configuration without reboot request:
```
---
values:
  grub_reset_to_defaults: true
  disable_reboot_request: true
```

Reset to default Grub2 configuration with reboot request:
```
---
values:
  grub_reset_to_defaults: true
```

**WRONG configuration** You cannot use `grub_reset_to_defaults` along with any Grub2 configuration option as it makes no sense and restricted by module's schema.json:
```
---
values:
  grub_reset_to_defaults: true
  options:
    grub_cmdline_linux:
      - 'cgroup_enable=memory'
      - 'debug'
      - 'intel_iommu=off'
```
This will result in failure of module's JSON schema validation in HostOSConfiguration resource.
