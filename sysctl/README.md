# sysctl module

The sysctl module allows the operator to manage kernel parameters at runtime on cluster machines using the mechanism implemented in the day-2
operations API.
Under the hood, this module is based on the [sysctl](https://docs.ansible.com/ansible/2.9/modules/sysctl_module.html) Ansible module.

> Note: This module is implemented and validated against the following Ansible versions provided by MOSK for Ubuntu 20.04, 22.04, 24.04:
> Ansible core 2.12.10 and Ansible collection 5.10.0.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

# Version 1.2.0 (latest)

Using the sysctl module 1.2.0, you can configure kernel parameters using the common `/etc/sysctl.conf` file or using a standalone file with ability
to clean up changes. The module contains the following input parameters:

- `filename`: Optional. Name of the file that stores the provided kernel parameters.
- `cleanup_before`: Optional. Enables cleanup of the dedicated file name before setting new parameters.
- `state`: Optional. Module state. Possible values are `present` (default) or `absent`.
- `options`: List of key-value kernel parameters to be applied on the machine.

   > Caution: For integer or float values, the system accepts only strings. For example, `1` -> `"1"`, `1.01` -> `"1.01"`.

# Configuration examples

Example of `HostOSConfiguration` with the sysctl module 1.2.0 for configuration of kernel parameters:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: sysctl-200
      namespace: default
    spec:
      configs:
      - module: sysctl
        moduleVersion: 1.2.0
        values:
          filename: custom
          cleanup_before: true
          options:
            net.ipv4.ip_forward: "1"
          state: present
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```

Example of `HostOSConfiguration` with the sysctl module 1.2.0 for dropping previously configured kernel parameters:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: sysctl-200
      namespace: default
    spec:
      configs:
      - module: sysctl
        moduleVersion: 1.2.0
        values:
          filename: custom
          cleanup_before: true
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```

# Versions 1.1.0 and 1.0.0 (deprecated)

> Note: The sysctl module 1.1.0 and 1.0.0 versions are obsolete and not recommended for usage in production environments.

Using the sysctl module version 1.0.0, you can configure kernel parameters using the common `/etc/sysctl.conf` file without the ability to roll back changes.

> Caution: For integer or float values, the system accepts only strings. For example, `1` -> `"1"`, `1.01` -> `"1.01"`.

Example of `HostOSConfiguration` with the sysctl module 1.0.0:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: sysctl-100
      namespace: default
    spec:
      configs:
      - module: sysctl
        moduleVersion: 1.0.0
        values:
          net.ipv4.ip_forward: "1"
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```
