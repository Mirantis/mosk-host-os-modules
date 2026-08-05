# package module

> **Note:** This README describes the latest version of the module. For version history and release notes, see the [changelog](changelog.md). Earlier versions are deprecated and not recommended for production environments.

The package module allows the operator to configure additional Ubuntu mirrors and install required packages from these mirrors on cluster machines using the mechanism implemented in the day-2 operations API. Under the hood, this module is based on [apt](https://docs.ansible.com/ansible/2.9/modules/apt_module.html) and [apt_repository](https://docs.ansible.com/ansible/2.9/modules/apt_repository_module.html) Ansible modules.

> Note: This module was initially implemented and validated against the following Ansible versions provided by MOSK for Ubuntu 22.04 and 24.04: Ansible core 2.12.10 and Ansible collection 5.10.0.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

## Configuration parameters

The module contains the following input parameters:

- `dpkg_options`: Optional. Comma-separated list of `dpkg` options to be used during package installation or removal. Defaults to `force-confold,force-confdef`.
- `os_version`: Optional. Version of the Ubuntu operating system. Possible values are `20.04`, `22.04`, and `24.04`. Applies to machines with the specified Ubuntu version.
  If not provided, the Ubuntu version is not verified by the module.

  > Caution: Use the deprecated Ubuntu `20.04` only on existing clusters based on this Ubuntu release.
  > For any other use case, use the latest supported Ubuntu release.

- `packages`: Optional. Map with packages to be installed using the `packages[*].<paramName>` parameters described below.
- `packages[*].name`: Required. Package name.
- `packages[*].allow_downgrade`: Optional. Enables downgrading of an installed package. Mirantis recommends setting `yes` when `version` is specified. Defaults to `no`.
- `packages[*].allow_unauthenticated`: Optional. Parameter that enables management of packages from unauthenticated sources. Defaults to `no`.
- `packages[*].autoremove`: Optional. Parameter that enables removal of unused dependency packages. Defaults to `no`.
- `packages[*].purge`: Optional. Parameter that enables purging of configuration files if a package state is `absent`. Defaults to `no`.
- `packages[*].state`: Optional. Module state. Possible values: `present`, `absent`, `build-dep`, `latest`, `fixed`.
- `packages[*].version`: Optional. Package version to be installed and pinned using the apt_preferences pinning. Mirantis recommends setting `allow_downgrade` to `yes` when `version` is specified.
- `repositories`: Optional. Configuration map of repositories to be managed on machines using the `repositories[*].<paramName>` parameters described below.
- `repositories[*].codename`: Optional. Code name of the repository.
- `repositories[*].filename`: Required. Name of the file that stores the repository configuration.
- `repositories[*].key`: Optional. URL of the repository GPG key.
- `repositories[*].repo`: Required. URL of the repository.
- `repositories[*].state`: Optional. Module state. Possible values are `present` (default) or `absent`.
- `repositories[*].validate_certs`: Optional. Validator of the repository SSL certificate. Default is `true`.

## Configuration examples

Example of `HostOSConfiguration` with the `package` module for installation of a repository and package:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: package-200
      namespace: default
    spec:
      configs:
        - module: package
          moduleVersion: 1.4.0
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

Example of `HostOSConfiguration` with the `package` module for installation of a package with specific version:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: package-200
      namespace: default
    spec:
      configs:
        - module: package
          moduleVersion: 1.4.0
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

Example of `HostOSConfiguration` with the `package` module for removal of the previously configured repository and package:

```
    apiVersion: kaas.mirantis.com/v1alpha1
    kind: HostOSConfiguration
    metadata:
      name: package-200
      namespace: default
    spec:
      configs:
        - module: package
          moduleVersion: 1.4.0
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
