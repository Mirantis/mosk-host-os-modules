# auditd module

The `auditd` module allows configuring the auditd rules at runtime on cluster machines using the mechanism implemented in the day-2 operations API.

> **Note:** This module supports Ubuntu 22.04 and 24.04 host OS.

> **Note:** This module is implemented and validated against the following Ansible versions provided by MOSK for Ubuntu 22.04 and Ubuntu 24.04 in the Cluster release 20.1.0: **Ansible Core 2.16.3** and **Ansible Collection 8.3.0**.
>
> To verify the Ansible version in a specific Cluster release, refer to the
> **Release artifacts > Management cluster artifacts > System and MCR artifacts**
> section of the required management Cluster release in
> [MOSK documentation: Release notes](https://docs.mirantis.com/mosk/latest/release-notes.html).

---

# Version 1.0.0 (latest)

The module supports the following input parameters:

* **`purge`**: `bool`, optional, default: `false`.
  Removes the auditd package and other module traces. If set to `true`, no other parameters take effect.

* **`enabled`**: `bool`, mandatory.
  Enables or disables auditd.
  CIS 4.1.1.1, CIS 4.1.1.2.

* **`enabledAtBoot`**: `bool`, optional, default: `false`.
  Configures GRUB to audit eligible processes even if they start up before auditd.
  CIS 4.1.1.3.

* **`backlogLimit`**: `int`, optional, default: `undefined`.
  If `audit=1` during boot, the backlog holds 64 records.
  Generating more than 64 records during boot may result in lost auditd records, allowing malicious activity to bypass detection.
  CIS 4.1.1.4.

* **`maxLogFile`**: `int`, optional, default: `8`.
  Configures the maximum size of the audit log file. Once the log reaches this size, it is rotated and a new log is started.
  CIS 4.1.2.1.

* **`maxLogFileAction`**: `string`, optional, default: `rotate`.
  Allowed values: `rotate`, `keep_logs`, `compress`

  * Defines how to handle audit logs when the maximum file size is reached.
  * `rotate`: Rotate logs, keep `maxLogFileKeep` files, delete oldest.
  * `keep_logs`: Rotate logs but never delete old ones.
  * `compress`: Same as `keep_logs`, plus a cron job compresses rotated log files, keeping up to **5 compressed files**.
  * CIS 4.1.2.2.

* **`maxLogFileKeep`**: `int`, optional, default: `5`.
  Used when `maxLogFileAction: compress`. Defines the number of compressed log files to retain in `/var/log/auditd/`.

* **`mayHaltSystem`**: `bool`, optional, default: `false`.
  Configures auditd to halt the system when audit logs are full. Applies the following settings:

  ```
  space_left_action = email
  action_mail_acct = root
  admin_space_left_action = halt
  ```
  CIS 4.1.2.3.

> Warning. The `mayHaltSystem` parameter locks the host system when logs reach capacity. Use this setting with extreme caution, as it will stop all system operations.

* **`weeklyLogOffload`**: `object`, optional, default: `undefined`.
  Configuration map for offloading local system audit logs to an external or alternative log location on a weekly basis. Used to fulfill DISA STIG compliance tracking rules for standalone nodes.
  **STIG UBTU-24-900950**.

  * **`enabled`**: `bool`, mandatory if `weeklyLogOffload` block is defined.
    Enables or disables the crontab script deployment inside `/etc/cron.weekly/` to offload local audit log tracks.
  * **`rsyncDest`**: `string`, mandatory if `enabled` is `true`.
    The destination directory target pathway. If remote `rsyncRemoteUser` and `rsyncRemoteHost` variables are configured, this represents the directory path located directly on the remote instance (e.g., `/home/storage/auditd_logs`). If remote variables are omitted, this is processed as a pure local system path wrapper (e.g., `/mnt/secure_backup`).
  * **`rsyncRemoteUser`**: `string`, optional, default: `""`.
    The SSH username token used to establish authorization scopes on the target log collection cluster node.
  * **`rsyncRemoteHost`**: `string`, optional, default: `""`.
    The remote receiver machine IP address or fully qualified network domain boundary.
  * **`rsyncAdditionalArgs`**: `string`, optional, default: `""`.
    Appends raw, explicit rsync string modifiers or runtime flags directly onto the automated execution query string (e.g., `--bwlimit=10m`).
  * **`rsyncSSHAdditionalArgs`**: `string`, optional, default: `""`.
    Appends raw configurations directly into the underlying SSH shell wrapper environment. Useful for specifying connection-specific arguments (e.g., `-p 22` or `-o ConnectTimeout=15`).

* **`presetRules`**: `string`, optional, default: `all,!immutable`.
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

* **`customRules`**: `string`, optional, default: `undefined`.
  Raw text content for a `60-custom.rules` file, applicable to any architecture.

# Weekly logs offloading configuration

The weekly logs offloading sub-feature automates compliance with DISA STIG mandate **UBTU-24-900950** by establishing an automated cron engine that pushes local `auditd` historical blocks to an auxiliary target infrastructure.

### Prerequisites

* **Remote Endpoint Setup:** The remote target collection environment must have a functional OpenSSH daemon (`sshd`) operational and the `rsync` utility binary installed on its native `$PATH`.
* **Key-Pair Authentication:** A cryptographic SSH key-pair must be generated. The corresponding public token must be successfully appended to the destination profile's `authorized_keys` template on the receiver machine.

### Deployment Workflow

1. **Generate the Kubernetes Secret Artifact:**
   Construct a secure Kubernetes generic secret container enclosing your deployment private key. The data object key name **must** map to `weeklyLogOffloadRsyncPrivateKey`.

   ```bash
   kubectl create secret generic auditd-weekly-offload \
     --namespace default \
     --from-file=weeklyLogOffloadRsyncPrivateKey=/home/ubuntu/.ssh/auditd_weekly_offload
   ```

2. **Map the Secret in HostOSConfiguration:**
   Reference the created secret under the `secretValues` section of your HostOSConfiguration. For more information regarding secret injection, refer to the [MOSK HostOSConfiguration Reference Guide](https://docs.mirantis.com/mosk/latest/ops/bm-operations/host-os-conf/hoc-description.html).
3. **Configure HostOSConfiguration according to your requirements using `weeklyLogOffload` subsection of `values` section.**
   See configuration example below.

# Configuration examples

Example `HostOSConfiguration` with the `auditd` module 1.0.0:

```yaml
apiVersion: kaas.mirantis.com/v1alpha1
kind: HostOSConfiguration
metadata:
  name: auditd-200
  namespace: mosk
spec:
  configs:
    - module: auditd
      moduleVersion: 1.0.0
      values:
        enabled: true
        enabledAtBoot: true
        maxLogFileAction: compress
        maxLogFileKeep: 5
        presetRules: "all,!delete,!immutable"
  machineSelector:
    matchLabels:
      day2-auditd-module: "true"
```

# Example of HostOSConfiguration for remote log offloading 

```yaml
apiVersion: kaas.mirantis.com/v1alpha1
kind: HostOSConfiguration
metadata:
  name: auditd-weekly-offload
  namespace: default
spec:
  configs:
  - module: auditd
    moduleVersion: X.X.X
    secretValues:
      name: auditd-weekly-offload
      namespace: default
    values:
      enabled: true
      weeklyLogOffload:
        enabled: true
        rsyncRemoteUser: "debian"
        rsyncRemoteHost: "172.19.120.25"
        rsyncDest: "/home/debian/auditd_logs"
        rsyncAdditionalArgs: "--bwlimit=10m"
        rsyncSSHAdditionalArgs: "-p 22 -o ConnectTimeout=15"
  machineSelector:
    matchLabels:
      auditd-label: 'true'
```