# package module

The package module allows the operator to configure additional Ubuntu mirrors and install required packages from these mirrors on cluster machines using the mechanism implemented in the day-2 operations API. Under the hood, this module is based on [apt](https://docs.ansible.com/ansible/2.9/modules/apt_module.html) and [apt_repository](https://docs.ansible.com/ansible/2.9/modules/apt_repository_module.html) Ansible modules.

> Note: This module is implemented and validated against the following Ansible versions provided by MCC for Ubuntu 22.04 in the Cluster releases 17.3.0: Ansible core 2.12.10 and Ansible collection 5.10.0.
>
> To verify the Ansible version in a specific Cluster release, refer to [Container Cloud documentation: Release notes - Cluster releases](https://docs.mirantis.com/container-cloud/latest/release-notes/cluster-releases.html).
> Use the *Artifacts > System and MCR artifacts* section of the corresponding Cluster release. For example, for
> [17.3.0](https://docs.mirantis.com/container-cloud/latest/release-notes/cluster-releases/17-x/17-3-x/17-3-0/17-3-0-artifacts.html#system-and-mcr-artifacts).

# Version 1.3.0 (latest)

Using the package module 1.3.0, you can configure additional Ubuntu mirrors and install packages from these mirrors on cluster machines with ability to specify and pin package versions. See documentation for the module version 1.2.0 below for more details.
New parameters comparing to 1.2.0 module version:

- `packages[*].allow_downgrade`: Optional. Parameter that enables downgrading of installed package. It is advised to set `yes` when `version` is specified. Defaults to `no`.
- `packages[*].version`: Optional. Package version to be installed and pinned via apt-preferences pinning. It is advised to set `allow_downgrade` to `yes` when `version` is specified.

# Configuration examples

Example of `HostOSConfiguration` with the `package` module 1.3.0 for installation of a package with specific version:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: package-200
      namespace: default
    spec:
      configs:
        - module: package
          moduleVersion: 1.3.0
          values:
            packages:
            - name: pinnedPackageName
              state: present
              version: 1.0.5-rc1
              allow_downgrade: yes
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```

# Version 1.2.0 (deprecated)

Using the package module 1.2.0, you can configure additional Ubuntu mirrors and install packages from these mirrors on cluster machines.
The module contains the following input parameters:

- `dpkg_options`: Optional. Comma-separated list of `dpkg` options to be used during package installation or removal. Defaults to `force-confold,force-confdef`.
- `os_version`: Optional. Version of the Ubuntu operating system. Possible values are `20.04` and `22.04`. Applies to machines with the specified Ubuntu version.
  If not provided, the Ubuntu version is not verified by the module.

  > Caution: Use the deprecated Ubuntu `20.04` only on existing clusters based on this Ubuntu release.
  > For any other use case, use the latest supported Ubuntu release `22.04`.

- `packages`: Optional. Map with packages to be installed using the `packages[*].<paramName>` parameters described below.
- `packages[*].name`: Required. Package name.
- `packages[*].allow_unauthenticated`: Optional. Parameter that enables management of packages from unauthenticated sources. Defaults to `no`.
- `packages[*].autoremove`: Optional. Parameter that enables removal of unused dependency packages. Defaults to `no`.
- `packages[*].purge`: Optional. Parameter that enables purging of configuration files if a package state is `absent`. Defaults to `no`.
- `packages[*].state`: Optional. Module state. Possible values: `present`, `absent`, `build-dep`, `latest`, `fixed`.
- `repositories`: Optional. Configuration map of repositories to be managed on machines using the `repositories[*].<paramName>` parameters described below.
- `repositories[*].codename`: Optional. Code name of the repository.
- `repositories[*].filename`: Required. Name of the file that stores the repository configuration.
- `repositories[*].key`: Optional. URL of the repository GPG key.
- `repositories[*].repo`: Required. URL of the repository.
- `repositories[*].state`: Optional. Module state. Possible values are `present` (default) or `absent`.
- `repositories[*].validate_certs`: Optional. Validator of the repository SSL certificate. Default is `true`.

# Configuration examples

Example of `HostOSConfiguration` with the `package` module 1.2.0 for installation of a repository and package:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: package-200
      namespace: default
    spec:
      configs:
        - module: package
          moduleVersion: 1.2.0
          values:
            dpkg_options: "force-confold,force-confdef"
            packages:
            - name: packageName
              state: present
            repositories:
            - filename: fileName
              key: https://example.org/packages/key.gpg
              repo: deb https://example.org/packages/ apt/stable/
              state: present
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```

Example of `HostOSConfiguration` with the `package` module 1.2.0 for removal of the previously configured repository and package:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: package-200
      namespace: default
    spec:
      configs:
        - module: package
          moduleVersion: 1.2.0
          values:
            packages:
            - name: packageName
              state: absent
            repositories:
            - filename: examplefile
              repo: deb https://example.org/packages/ apt/stable/
              state: absent
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```

# Versions 1.1.0 and 1.0.0 (deprecated)

> Note: The package module 1.1.0 and 1.0.0 versions are obsolete and not recommended for usage in production environments.

Using the package module version 1.0.0, you can install packages from already configured mirrors only. It cannot configure additional mirrors.

The module input values are a map of key-value pairs, where the key is a package name and the value is a package state (`present` or `absent`).

Example of `HostOSConfiguration` with the package module 1.0.0:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: package-100
      namespace: default
    spec:
      configs:
      - module: package
        moduleVersion: 1.0.0
        values:
          package1: present
          package2: absent
      machineSelector:
        matchLabels:
          day2-custom-label: "true"
```
