# disa_stig module

The `disa_stig` module allows configuring the auditd rules at runtime on cluster machines using the mechanism implemented in the day-2 operations API.

> **Note:** This module supports Ubuntu 24.04 host OS.

> **Note:** This module is implemented and validated against the following Ansible versions provided by MOSK for Ubuntu 24.04 in the Cluster release XX.X.X: **Ansible Core X.XX.X** and **Ansible Collection X.X.X**.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

---

# Version 1.0.0 (latest)

The module supports the following input parameters:

* **`enabled`**: `bool`, mandatory.
  Enables the module to do the changes.

List of supported compliance rules:

* UBTU-24-100850: SSH client must use FIPS 140-3 approved ciphers
* UBTU-24-100860: SSH client must use FIPS 140-3 validated MACs
* UBTU-24-100820: FIPS-140-3 ciphers
* UBTU-24-200640: Display the Standard Mandatory DOD Notice and Consent Banner
* UBTU-24-300023: Prevent remote hosts from connecting to the proxy display
* UBTU-24-400030: Implement smart card logins for MFA
* UBTU-24-600000: Terminate traffic after a period of inactivity
* UBTU-24-600010: Terminate after 10 minutes (600 seconds) of inactivity

> **Note:** SSH/SSHD rules are configured via separate /etc/ssh/ssh(/d)_config.d/01-stig.conf file. This could lead to false-negative results for corresponding benchmark's tests.

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
  machineSelector:
    matchLabels:
      day2-disa-stig-module: "true"
```
