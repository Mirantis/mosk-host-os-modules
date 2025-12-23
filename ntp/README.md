# ntp module

The ntp module allows the operator to manage ntp servers at runtime on cluster machines using the mechanism implemented in the day-2 operations API.

> Note: This module is implemented and validated against the following Ansible versions provided by MOSK for Ubuntu 22.04
> in the Cluster releases 16.3.0 and 17.3.0: Ansible core 2.12.10 and Ansible collection 5.10.0.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

# Version 1.0.0 (latest)

Using the ntp module 1.0.0, you can configure list of ntp servers.
The module contains the following input parameters:

- `ntp_servers`: List of NTP server.

# Configuration examples

Example of `HostOSConfiguration` with the ntp module 1.0.0 for configuration of ntp servers:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: ntp-200
      namespace: default
    spec:
      configs:
      - module: ntp
        moduleVersion: 1.0.0
        values:
          ntp_servers:
            - 0.ubuntu.pool.ntp.org
            - 1.ubuntu.pool.ntp.org
            - 2.ubuntu.pool.ntp.org
            - 3.ubuntu.pool.ntp.org
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```
