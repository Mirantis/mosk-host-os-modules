# 2.0.0

- Added support for the DISA STIG `Canonical Ubuntu 24.04 LTS STIG, V1R5` hardening of the `auditd` configuration:
  - Added the `stig` preset to the `presetRules` parameter that applies the DISA STIG audit rules. Developed and validated only for the Ubuntu 24.04 host OS.
  - Updated the `perm-mod` preset with STIG-compatible rules.
  - Enforced strict access permissions on the auditd log files.
  - Added the `runtimeLogUpload` parameter that configures the `auditd` daemon to upload local audit logs to an external receiver host at runtime using the native `audisp-remote` plugin.
  - Added the `weeklyLogUpload` parameter that deploys a weekly cron task to offload local audit logs to an external or alternative log location using `rsync` (DISA STIG UBTU-24-900950).

# 1.0.0

- Initial module version.
