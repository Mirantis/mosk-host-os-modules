# Host OS configuration ansible modules for MCC

This repository contains Ansible modules for Mirantis Container Cloud that are released to customers.

Modules manage various system configurations through Ansible playbooks, adhering to specific schemas and metadata requirements.

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

All modules and `index.yaml` are built [per-commit by Jenkins](https://ci.mcp.mirantis.net/job/kaas-bm-kaas-bm-host-os-modules-build/) using [`make-artifact.groovy`](https://gerrit.mcp.mirantis.com/plugins/gitiles/mcp/mcp-pipelines/+/refs/heads/master/make-artifact.groovy) pipeline, that runs the `Makefile` in a container.
 Merged modules are then avaiable on [`binary-dev-kaas-local` internal artifactory](https://artifactory.mcp.mirantis.net/ui/native/binary-dev-kaas-local/bm/bin/host-os-modules/) and in `master` branch of `artifact-metadata`.

Use `make promote` to promote latest modules version in the repository, so new non-development versions are set for every module and all dev versions are removed from `index.yaml`.

In time for release, move `artifact-metadata` items to `release` branch to release them onto <https://binary.mirantis.com/?prefix=bm/bin/host-os-modules/>.
