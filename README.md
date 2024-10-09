# Host OS configuration ansible modules for MCC

This repository contains Ansible modules for Mirantis Container Cloud that are released to customers.

Modules manage various system configurations through Ansible playbooks, adhering to specific schemas and metadata requirements.

## TLDR

1. Make changes to a module or introduce a new one
1. **Do not change `version` field of the `metadata.yaml` files manually**
1. `make`
1. `git add -A && git commit`
1. `git push -u origin HEAD:refs/for/master`
1. Test the new CR
1. `make promote && git push -u origin HEAD:refs/for/master`

## Module structure

Each module directory contains specific files crucial for the operation of the OS Modules.

The common structure within each module is as follows:

- `main.yaml`: The primary playbook file that defines tasks to be executed.
- `metadata.yaml`: Provides metadata about the module like name, version, and relevant documentation URLs.
- `schema.json`: Defines the [JSON schema](https://json-schema.org/overview/what-is-jsonschema) for validating configurations specific to the module (i.e. restricted values).

## Module index

> Important - make sure to run `make` locally when updating modules, to keep `index.yaml` up to date.

Modules are indexed in `index.yaml` file that is stored in this repo. Before committing to gerrit, please run `make` to ensure fresh build and correct sha256 sums for module artifacts.

The directory `.githooks` contains Git hooks that can be used to enforce building modules and `index.yaml`.

**Do not change** `metadata.yaml`'s `version` field manually, it will be changed automatically after running the `make` before committing.

### Installing the Hooks

To install the hooks, run the following command:

```bash
cp ./.githooks/pre-commit ./.git/hooks
```

### Available Hooks

`pre-commit`: runs `make build tgz sort-index` before committing changes and compares
the output of the `git diff` command before and after.

## Build

Requirements:

- `go`
- `make`
- `shasum`

Modules and `index.yaml` are built using `cmd/module-builder.go` to ensure reproduceable tar.gz builds.

## MCC implementation details

Modules are installed and controlled through two CRs in the management cluster:

- [HostOSConfiguration](https://docs.mirantis.com/container-cloud/latest/api/bm/host-os-configuration.html)
- [HostOSConfigurationModule](https://docs.mirantis.com/container-cloud/latest/api/bm/host-os-configuration-modules.html)

Ansible module execution is implemented via existing LCM mechanism, by creating an additional `StateItem` for mapped `LCMMachines` in MCC Management Cluster.

## Modules provided by Container Cloud

Check a full list of [day2 modules implemented by Container Cloud](https://docs.mirantis.com/container-cloud/latest/operations-guide/operate-managed/operate-managed-bm/day2/mcc-day2-modules.html) on Mirantis documentation portal.

## Release process

All modules and `index.yaml` are built per-commit by Jenkins using pipeline, that runs the `Makefile` in a container. Merged modules are then avaiable on internal artifactory and in `master` branch of `artifact-metadata`.

Use `make promote` to promote latest modules version in the repository, so new non-development versions are set for every module and all dev versions are removed from `index.yaml`.

In time for release, move `artifact-metadata` items to `release` branch to release them onto <https://binary.mirantis.com/?prefix=bm/bin/host-os-modules/>.
