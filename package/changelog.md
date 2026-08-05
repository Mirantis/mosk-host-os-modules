# 1.4.0

- Added support for Ubuntu `24.04`.

# 1.3.0

- Added the `packages[*].version` parameter that allows installing and pinning a specific package version using the apt_preferences pinning.
- Added the `packages[*].allow_downgrade` parameter that enables downgrading of an installed package.

# 1.2.0

- Added the ability to configure additional Ubuntu mirrors and install packages from these mirrors on cluster machines using the `repositories` and extended `packages` parameters.

# 1.1.0

- Reworked the module input parameters and added JSON schema validation of module values.

# 1.0.0

- Initial module version.
