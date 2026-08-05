# NTP module

> **Note:** This README describes the latest version of the module. For version history and release notes, see the [changelog](changelog.md). Earlier versions are deprecated and not recommended for production environments.

The NTP module allows the operator to manage NTP servers at runtime on cluster machines using the mechanism implemented in the host operating system configuration API.

> Note: This module was initially implemented and validated against the following Ansible versions provided by MOSK for Ubuntu 22.04
> in the Cluster releases 16.3.0 and 17.3.0: Ansible core 2.12.10 and Ansible collection 5.10.0.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

## Configuration parameters

The module contains the following input parameters:

- `ntp_servers`: List of NTP servers.
- `stigHardening`: Optional. Enables fixes for the following DISA STIG tests.

  - UBTU-24-600160: Ubuntu 24.04 LTS must compare internal information system clocks at least every 24 hours with an authoritative time server.
  - UBTU-24-600180: Ubuntu 24.04 LTS must synchronize internal information system clocks to the authoritative time source when the time difference is greater than one second.

  Defaults to `False`. Set to `True` to enable these fixes. Before enabling, ensure that you are using an official, trusted time source approved by the Department of Defense (DoD) or the U.S. Federal Government.

## Configuration examples

Example of `HostOSConfiguration` with the NTP module for configuration of NTP servers:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: ntp-200
      namespace: default
    spec:
      configs:
      - module: ntp
        moduleVersion: 1.1.0
        values:
          ntp_servers:
            - 0.ubuntu.pool.ntp.org
            - 1.ubuntu.pool.ntp.org
            - 2.ubuntu.pool.ntp.org
            - 3.ubuntu.pool.ntp.org
          stigHardening: False
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```
