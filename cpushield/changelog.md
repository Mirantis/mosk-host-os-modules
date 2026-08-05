# 1.1.0

- Added the `apply_settings_immediately` parameter that enables `systemctl daemon-reload` to immediately apply CPU/NUMA pinning settings for units from `systemd_units_to_pin`.
- Fixed the module behavior to not implicitly run `systemd daemon-reload` by Ansible, so that shielding settings take effect only after a node reboot or when `apply_settings_immediately` is explicitly enabled.

# 1.0.0

- Initial module version.
