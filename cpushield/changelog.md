# 1.2.0

- Added native k0s orchestrator support with dynamic detection (starting from MOSK 2.33/27.1 release).
- Automated mapping of `system_cpus` to Kubelet `reservedSystemCPUs` and `cpuManagerPolicy: static`.
- Fixed Kubelet crash loops by forcefully clearing conflicting k0s hardcoded cgroups (`kubeletCgroups`, `kubeReservedCgroup`, `systemReservedCgroup`, `systemCgroups`).
- Added automatic deletion of `/var/lib/kubelet/cpu_manager_state` on CPU boundary changes.
- Added pre-flight check to safely block legacy KaaS deployments with k0s (MOSK 2.32/26.2 release, mgmt cluster only) as they do not support /etc/kubernetes/kubelet.conf.d directory passing to Kubelet.

# 1.1.0

- Added the `apply_settings_immediately` parameter that enables `systemctl daemon-reload` to immediately apply CPU/NUMA pinning settings for units from `systemd_units_to_pin`.
- Fixed the module behavior to not implicitly run `systemd daemon-reload` by Ansible, so that shielding settings take effect only after a node reboot or when `apply_settings_immediately` is explicitly enabled.

# 1.0.0

- Initial module version.
