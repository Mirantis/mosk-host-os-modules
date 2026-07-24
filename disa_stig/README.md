# disa_stig module

The `disa_stig` module supports configuring the host operating system at runtime to comply with the DISA STIG `Canonical Ubuntu 24.04 LTS STIG, V1R5` on cluster machines using the mechanism implemented in the day-2 operations API.

> **Note:** This module supports Ubuntu 24.04 host OS only.

> **Note:** This module is implemented and validated against the following Ansible versions provided by MOSK for Ubuntu 24.04 in the Cluster release XX.X.X: **Ansible Core X.XX.X** and **Ansible Collection X.X.X**.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

---

# Version 1.0.0 (latest)

> **WARNING:** Changes made by the module cannot be reverted, except for the external USB storage setting.

> **Note:** `lcm-ansible` can overwrite some changes made by the module. After each LCM operation, reapply the HostOSConfiguration object. For details, see [MOSK documentation: Retrigger a module configuration](https://docs.mirantis.com/mosk/latest/ops/bm-operations/host-os-conf/day2-crd-hoc-retrigger.html).

The module supports the following input parameters:

* **`enabled`**: `bool`, mandatory.
  Enables the module.
* **`disableUsbStorage`**: `bool`, optional, default is `undefined`.
  Disables any external USB storage according to `DISA STIG UBTU-24-300039`.
  * Set to `true` to comply with the DISA STIG requirement.
  * Set to `false` to revert the settings (for example, if an external USB storage is required for host maintenance).

---

# Configuration examples

Example `HostOSConfiguration` with the `disa_stig` module 1.0.0:

```yaml
apiVersion: kaas.mirantis.com/v1alpha1
kind: HostOSConfiguration
metadata:
  name: disa-stig-compliance
  namespace: default
spec:
  configs:
    - module: disa_stig
      moduleVersion: 1.0.0
      values:
        enabled: true
        disableUsbStorage: true
  machineSelector:
    matchLabels:
      disa-stig-module: "true"
```
