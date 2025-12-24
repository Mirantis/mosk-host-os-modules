# auditd module

The `auditd` module allows configuring the auditd rules at runtime on cluster machines using the mechanism implemented in the day-2 operations API.

> **Note:** This module is implemented for the Ubuntu 22.04 and 24.04 host OS.
>
> **Note:** This module is implemented and validated against the following Ansible versions provided by MOSK for Ubuntu XX.04 in the Cluster release XXX: **Ansible Core XXX** and **Ansible Collection XXX**.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

---

# Version 1.0.0 (latest)

The module supports the following input parameters:

* **`purge`**: `bool`, optional, default: `false`
  Removes the auditd package and other module traces. If set to `true`, no other parameters take effect.

* **`enabled`**: `bool`, mandatory
  Enables or disables auditd.
  CIS 4.1.1.1, CIS 4.1.1.2.

* **`enabledAtBoot`**: `bool`, optional, default: `false`
  Configures GRUB so that processes capable of being audited can be audited even if they start up before auditd.
  CIS 4.1.1.3.

* **`backlogLimit`**: `int`, optional, default: `undefined`
  During boot, if `audit=1`, then the backlog holds 64 records.
  If more than 64 records are created during boot, auditd records may be lost and malicious activity could go undetected.
  CIS 4.1.1.4.

* **`maxLogFile`**: `int`, optional, default: `8`
  Configures the maximum size of the audit log file. Once the log reaches this size, it is rotated and a new log is started.
  CIS 4.1.2.1.

* **`maxLogFileAction`**: `string`, optional
  Allowed values: `rotate`, `keep_logs`, `compress`

  * Defines how to handle audit logs when the maximum file size is reached.
  * `keep_logs`: Rotate logs but never delete old ones.
  * `compress`: Same as `keep_logs`, plus a cron job compresses rotated log files, keeping up to **5 compressed files**.
  * CIS 4.1.2.2.

* **`maxLogFileKeep`**: `int`, optional, default: `5`
  Used when `maxLogFileAction: compress`. Defines the number of compressed log files to retain in `/var/log/auditd/`.

* **`mayHaltSystem`**: `bool`, optional, default: `false`
  Configures auditd to halt the system when audit logs are full. Applies the following settings:

  ```
  space_left_action = email
  action_mail_acct = root
  admin_space_left_action = halt
  ```

> WARNING: Expected behavior for `mayHaltSystem` is to LOCK host system when logs are full.
> USE WITH CAUTION!

CIS 4.1.2.3.

* **`presetRules`**: `string`, optional, default: `all,!immutable`
  A comma-separated list of preset rules:
  `docker`, `time-change`, `identity`, `system-locale`, `mac-policy`, `logins`,
  `session`, `perm-mod`, `access`, `privileged`, `mounts`, `delete`,
  `scope`, `actions`, `modules`, `immutable`.

  * Special value `all` enables all rules above.
  * Prefix a rule with `!` to exclude it (e.g., `all,!immutable`).

  **CIS controls:**

  * CIS 4.1.3 (time-change)
  * CIS 4.1.4 (identity)
  * CIS 4.1.5 (system-locale)
  * CIS 4.1.6 (mac-policy)
  * CIS 4.1.7 (logins)
  * CIS 4.1.8 (session)
  * CIS 4.1.9 (perm-mod)
  * CIS 4.1.10 (access)
  * CIS 4.1.11 (privileged)
  * CIS 4.1.12 (mounts)
  * CIS 4.1.13 (delete)
  * CIS 4.1.14 (scope)
  * CIS 4.1.15 (actions)
  * CIS 4.1.16 (modules)
  * CIS 4.1.17 (immutable)

  **Docker CIS controls:**

  * CIS 1.2.3
  * CIS 1.2.4
  * CIS 1.2.5
  * CIS 1.2.6
  * CIS 1.2.7
  * CIS 1.2.10
  * CIS 1.2.11

* **`customRules`**: `string`, optional, default: `undefined`
  Raw text content for a `60-custom.rules` file, applicable to any architecture.

---

# Configuration examples

Example `HostOSConfiguration` with the `auditd` module X.X.X:

```yaml
apiVersion: kaas.mirantis.com/v1alpha1
kind: HostOSConfiguration
metadata:
  name: auditd-200
  namespace: mosk
spec:
  configs:
    - module: auditd
      moduleVersion: X.X.X
      values:
        enabled: true
        enabledAtBoot: true
        maxLogFileAction: compress
        maxLogFileKeep: 5
        mayHaltSystem: true
  machineSelector:
    matchLabels:
      day2-auditd-module: "true"
```
