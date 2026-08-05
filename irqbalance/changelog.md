# 1.1.0

- Added the `oneshot` parameter.
- Changed the method of setting empty values for the irqbalance parameters for better usability:

  - When a parameter is not defined in `values` of the `HOC` object, the corresponding value remains the same in the irqbalance configuration file.
  - When a parameter is set to `""` (empty string) in `values` of the `HOC` object, the corresponding value in the `irqbalance` configuration file
    is also set to `""` (empty string).

# 1.0.0

- Initial module version.
