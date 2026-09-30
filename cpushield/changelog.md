# 1.2.0

- Added native k0s orchestrator support with dynamic detection. In MOSK, this feature applies starting with MOSK management 2.33.0 (MOSK 27.1).
- Automated the mapping of `system_cpus` to kubelet `reservedSystemCPUs` and `cpuManagerPolicy: static`.
- Fixed kubelet crash loops by forcefully clearing conflicting k0s hardcoded cgroups (`kubeletCgroups`, `kubeReservedCgroup`, `systemReservedCgroup`, `systemCgroups`).
- Added automatic deletion of `/var/lib/kubelet/cpu_manager_state` on CPU boundary changes.
- Added a pre-flight check that blocks deployment on k0s-based MOSK management clusters of the 2.32.0 (MOSK 26.2) release, because they do not support passing the `/etc/kubernetes/kubelet.conf.d` directory to kubelet.

# 1.1.0

- Added the `apply_settings_immediately` parameter that enables `systemctl daemon-reload` to immediately apply CPU/NUMA pinning settings for units from `systemd_units_to_pin`.
- Fixed the module behavior to not implicitly run `systemd daemon-reload` by Ansible, so that shielding settings take effect only after a node reboot or when `apply_settings_immediately` is explicitly enabled.

# 1.0.0

- Initial module version.
