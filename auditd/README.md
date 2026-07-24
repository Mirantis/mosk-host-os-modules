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
## Version 2.0.0 (latest)

The module supports the DISA STIG `Canonical Ubuntu 24.04 LTS STIG, V1R5` hardening for the `auditd` configuration.

The module contains the following input parameters:

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
  Configures the maximum size in MiB of the audit log file. Once the log reaches this size, it is rotated and a new log is started.
  CIS 4.1.2.1.

* **`maxLogFileAction`**: `string`, optional, default: `rotate`.
  Allowed values: `rotate`, `keep_logs`, `compress`

  * Defines how to handle audit logs when the maximum file size is reached.
  * `rotate`: Rotate logs, keep `maxLogFileKeep` files, delete oldest.
  * `keep_logs`: Rotate logs but never delete old ones.
  * `compress`: Same as `keep_logs`, plus a cron job compresses rotated log files, keeping up to **5 compressed files**.

  CIS 4.1.2.2.

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

  > **Warning:** The `mayHaltSystem` parameter locks the host system when logs reach capacity. Use this setting with extreme caution, as it will stop all system operations.

* **`runtimeLogUpload`**: `object`, optional, default: `undefined`.
  Configuration map for runtime uploading of local system audit logs to an external host with `auditd` running as receiver.

  * **`enabled`**: `bool`, mandatory if the `runtimeLogUpload` block is defined.
    Enables or disables the runtime log uploading.
  * **`server`**: `string`, mandatory if `enabled` is set to `true`.
    Remote receiver host.
  * **`port`**: `string`, optional, default: `"60"`.
    Port on the remote host where the auditd receiver is running.

* **`weeklyLogUpload`**: `object`, optional, default: `undefined`.
  Configuration map for offloading local system audit logs to an external or alternative log location on a weekly basis. Used to fulfill Defense Information Systems Agency (DISA) Security Technical Implementation Guide (STIG) compliance tracking rules for standalone nodes.
  For more details, see the [DISA STIG compliance configuration](#disa-stig-compliance-configuration) section.

  * **`enabled`**: `bool`, mandatory if `weeklyLogUpload` block is defined.
    Enables or disables the crontab script deployment inside `/etc/cron.weekly/` to offload local audit log tracks.
  * **`rsyncDest`**: `string`, mandatory if `enabled` is `true`.
    The destination directory target path. If remote `rsyncRemoteUser` and `rsyncRemoteHost` variables are configured, this represents the directory path located directly on the remote instance (e.g., `/home/storage/auditd_logs`). If remote variables are omitted, this is processed as a pure local system path wrapper (e.g., `/mnt/secure_backup`).
  * **`rsyncRemoteUser`**: `string`, optional, default: `""`.
    The SSH username token used to establish authorization scopes on the target log collection cluster node.
  * **`rsyncRemoteHost`**: `string`, optional, default: `""`.
    The remote receiver machine IP address or fully qualified network domain boundary.
  * **`rsyncAdditionalArgs`**: `string`, optional, default: `""`.
    Appends raw, explicit rsync string modifiers or runtime flags directly onto the automated execution query string (e.g., `--bwlimit=10m`).
  * **`rsyncSSHAdditionalArgs`**: `string`, optional, default: `""`.
    Appends raw configurations directly into the underlying SSH shell wrapper environment. Useful for specifying connection-specific arguments (e.g., `-p 22` or `-o ConnectTimeout=15`).

* **`presetRules`**: `string`, optional, default: `all,!stig,!immutable`.
  A comma-separated list of preset rules:
  `docker`, `time-change`, `identity`, `system-locale`, `mac-policy`, `logins`,
  `session`, `perm-mod`, `access`, `privileged`, `mounts`, `delete`,
  `scope`, `actions`, `modules`, `immutable`, `stig`.

  * Special value `all` enables all rules above.
  * Prefix a rule with `!` to exclude it (e.g., `all,!stig,!immutable`).

  > **Note:** The `stig` preset is developed and validated only for Ubuntu 24.04 host OS.

  Compared to the module version 1.0.0, the `perm-mod` preset is updated with STIG-compatible rules.

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
  Base64-encoded content for a `60-custom.rules` file, applicable to any architecture.

## Version 1.0.0

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
  Configures the maximum size in MiB of the audit log file. Once the log reaches this size, it is rotated and a new log is started.
  CIS 4.1.2.1.

* **`maxLogFileAction`**: `string`, optional, default: `rotate`.
  Allowed values: `rotate`, `keep_logs`, `compress`

  * Defines how to handle audit logs when the maximum file size is reached.
  * `rotate`: Rotate logs, keep `maxLogFileKeep` files, delete oldest.
  * `keep_logs`: Rotate logs but never delete old ones.
  * `compress`: Same as `keep_logs`, plus a cron job compresses rotated log files, keeping up to **5 compressed files**.

  CIS 4.1.2.2.

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
  Base64-encoded content for a `60-custom.rules` file, applicable to any architecture.

## Auditd customRules configuration

The module supports configuring additional `auditd` rules. To enable the rules:

1. Prepare a file with required rules. This file will be passed "as-is" to the `60-custom.rules` file on the target host.
   ```shell
   ~$ cat my_rules.txt
   -w /etc/apt -p wa -k fim-auditbeat
   -w /etc/docker -p wa -k fim-auditbeat
   ```

2. Encode the file content with the `base64` tool:
   ```shell
   ~$ base64 my_rules.txt
   ICAtdyAvZXRjL2FwdCAtcCB3YSAtayBmaW0tYXVkaXRiZWF0CiAgLXcgL2V0Yy9kb2NrZXIgLXAg
   d2EgLWsgZmltLWF1ZGl0YmVhdAo=
   ```

3. Pass the resulting base64-encoded text to the `customRules` field:
   ```yaml
   spec:
     configs:
     ...
       values:
         customRules: |
           ICAtdyAvZXRjL2FwdCAtcCB3YSAtayBmaW0tYXVkaXRiZWF0CiAgLXcgL2V0Yy9kb2NrZXIgLXAg
           d2EgLWsgZmltLWF1ZGl0YmVhdAo=
   ```

## DISA STIG compliance configuration

The module provides an automated way to configure the auditd daemon in strict compliance with DISA STIG requirements.
To enable the configuration, execute the following deployment steps.

### Enable DISA STIG-related rules

Ensure that `stig` and `perm-mod` presets are enabled in `presetRules`.

### Configure runtime log uploading

The module allows configuring the `auditd` daemon to push runtime log streams directly to a central log server using the native `audisp-remote` plugin architecture.

#### Prerequisites

The destination collector host must meet the following requirements:

* **Operating system target:** The destination collector host must not be running Ubuntu as its operating system. This requirement exists because the Ubuntu upstream `auditd` `.deb` packages are compiled without network listener support (`--with-listener=no`), preventing them from accepting remote connections.
* **Package distribution:** The destination collector host must have both the `auditd` package and `audispd-plugins` installed.
* **Server aggregation setup:** The `auditd` daemon on the destination collector host must be configured to bind to network interfaces and process incoming TCP operations. To do this, specify the following parameters in `/etc/audit/auditd.conf`:

  ```ini
  # /etc/audit/auditd.conf on the remote server
  tcp_listen_port = 60
  tcp_listen_queue = 5
  tcp_max_per_addr = 1
  tcp_client_max_idle = 0
  ```

#### Deployment workflow

To configure log uploading onto target hosts, add the `runtimeLogUpload` map to the `HostOSConfiguration` object. For example:

```yaml
spec:
  configs:
  ...
    values:
      runtimeLogUpload:
        enabled: true
        server: remote.server.ip.or.host
```

You can find a complete configuration example for remote log uploading with the `auditd` module 2.0.0 in [Configuration examples](#configuration-examples).

### Configure weekly log uploading

Using the weekly log uploading feature, you can create a cron task to push local `auditd` historical blocks to an auxiliary target infrastructure using `rsync`.

#### Prerequisites

* **Remote Endpoint Setup:** The remote target collection environment must have a functional OpenSSH daemon (`sshd`) operational and the `rsync` utility binary installed on its native `$PATH`.
* **Key-Pair Authentication:** A cryptographic SSH key-pair must be generated. The corresponding public token must be successfully appended to the destination profile's `authorized_keys` template on the receiver machine.

#### Deployment workflow

1. **Generate the Kubernetes secret artifact:**
   Create a secure Kubernetes generic secret container enclosing your deployment private key. The data object key name must map to `weeklyLogUploadRsyncPrivateKey`.

   ```bash
   kubectl create secret generic auditd-weekly-upload \
     --namespace default \
     --from-file=weeklyLogUploadRsyncPrivateKey=/home/ubuntu/.ssh/auditd_weekly_upload
   ```

2. **Map the secret in HostOSConfiguration:**
   Add a reference to the created secret in the `secretValues` section of the HostOSConfiguration. For more details on secret injection, refer to the [MOSK documentation: HostOSConfiguration and HostOSConfigurationModules concepts](https://docs.mirantis.com/mosk/latest/ops/bm-operations/host-os-conf/hoc-description.html).
3. **Configure HostOSConfiguration according to your requirements using the `weeklyLogUpload` subsection of the `values` section.**
   For details, see the configuration example for remote log uploading and
   weekly log uploading with the `auditd` module 2.0.0 in [Configuration examples](#configuration-examples).

## Configuration examples

Example of `HostOSConfiguration` for remote log uploading and weekly log uploading with the `auditd` module 2.0.0:

```yaml
apiVersion: kaas.mirantis.com/v1alpha1
kind: HostOSConfiguration
metadata:
  name: auditd-remote-upload
  namespace: default
spec:
  configs:
  - module: auditd
    moduleVersion: 2.0.0
    secretValues:
      name: auditd-weekly-upload
      namespace: default
    values:
      enabled: true
      runtimeLogUpload:
        enabled: true
        server: "remote.server.ip.or.host"
        port: "60"
      weeklyLogUpload:
        enabled: true
        rsyncRemoteUser: remote-username
        rsyncRemoteHost: remote.storage.server.ip.or.host
        rsyncDest: /home/debian/auditd_logs
        rsyncAdditionalArgs: --bwlimit=10m
        rsyncSSHAdditionalArgs: -p 22 -o ConnectTimeout=15
  machineSelector:
    matchLabels:
      auditd-label: 'true'
```

Example `HostOSConfiguration` with the `auditd` module 1.0.0:

```yaml
apiVersion: kaas.mirantis.com/v1alpha1
kind: HostOSConfiguration
metadata:
  name: auditd-100
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
        customRules: |
          ICAtdyAvZXRjL2FwdCAtcCB3YSAtayBmaW0tYXVkaXRiZWF0CiAgLXcgL2V0Yy9kb2NrZXIgLXAg
          d2EgLWsgZmltLWF1ZGl0YmVhdAo=
  machineSelector:
    matchLabels:
      day2-auditd-module: "true"
```
