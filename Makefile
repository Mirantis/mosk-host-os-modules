ARTIFACTS_DIR ?= $(CURDIR)/_artifacts
MODULES_LIST ?= $(shell find * -maxdepth 0 -type d ! -name cmd ! -name $(shell basename $(ARTIFACTS_DIR)))
PROMOTE ?= ""

all: clean dirs build metadata sort-index index index-meta

.PHONY: cicd-build
cicd-build: all check-diff

.PHONY: promote
promote: promote-minor

.PHONY: promote-minor
promote-minor: PROMOTE=minor
promote-minor: all git-promote-commit

.PHONY: promote-major
promote-major: PROMOTE=major
promote-major: all git-promote-commit

define meta
	$(eval MODULE_NAME = $(1))
	$(eval VERSION = $(shell grep -E '^version:' $(MODULE_NAME)/metadata.yaml | awk {'print $$2'}))
	$(eval MODULE_TGZ = $(MODULE_NAME)-$(VERSION).tgz)
	$(eval SHASUM = $(shell shasum -a 256 $(ARTIFACTS_DIR)/$(MODULE_TGZ) | awk {'print $$1'}))
	$(eval ARTIFACT_METADATA_KEY = binary:bm:host-os-modules:$(MODULE_NAME))

	@printf 'key: $(ARTIFACT_METADATA_KEY)\n' > "$(ARTIFACTS_DIR)/$(MODULE_TGZ).metadata.yaml"
	@printf 'version: $(VERSION)\n' >> "$(ARTIFACTS_DIR)/$(MODULE_TGZ).metadata.yaml"
	@printf 'sha256sum: $(SHASUM)\n' >> "$(ARTIFACTS_DIR)/$(MODULE_TGZ).metadata.yaml"
endef

.PHONY: tgz
tgz:
	$(CURDIR)/cmd/module-builder module --promote=$(PROMOTE) --output=$(ARTIFACTS_DIR) $(MODULES_LIST)

.PHONY: sort-index
sort-index:
	$(CURDIR)/cmd/module-builder sort

.PHONY: index
index:
	cp index.yaml $(ARTIFACTS_DIR)

.PHONY: index-meta
index-meta:
	@printf 'key: binary:bm:host-os-modules:index\n' > "$(ARTIFACTS_DIR)/index.yaml.metadata.yaml"

.PHONY: metadata
metadata: tgz
	$(foreach dir,$(MODULES_LIST),$(call meta, $(dir)))

.PHONY: clean
clean:
	rm -rf $(ARTIFACTS_DIR) $(VENV_DIR)

.PHONY: dirs
dirs:
	mkdir -p "$(ARTIFACTS_DIR)"

.PHONY: vet
vet:
	@cd $(CURDIR)/cmd && go vet ./...

.PHONY: tidy
tidy:
	@cd $(CURDIR)/cmd && go mod tidy

.PHONY: generate
generate:
	@cd $(CURDIR)/cmd && go generate ./...

.PHONY: test
test:
	@cd $(CURDIR)/cmd && go test -v ./...

.PHONY: build
build: generate test vet
	cd $(CURDIR)/cmd; go build -ldflags "-s -w" -o ./module-builder ./

.PHONY: check-diff
check-diff:
	git diff --exit-code

.PHONY: git-promote-commit
git-promote-commit:
	! git diff --exit-code $(MODULES_LIST) index-dev.yaml
	! git diff --exit-code $(MODULES_LIST) index.yaml
	git add $(MODULES_LIST) index-dev.yaml index.yaml
	git commit -m "[promote] Release latest modules"


.PHONY: list-modules
list-modules:
	@echo $(MODULES_LIST)
