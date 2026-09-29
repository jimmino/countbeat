GO ?= go
DAGGER ?= dagger

WORK_DIR   := $(shell pwd)

GITEA_SDK_TEST_URL ?= http://localhost:3000
GITEA_SDK_TEST_USERNAME ?= test01
GITEA_SDK_TEST_PASSWORD ?= test01234
GITEA_SDK_TEST_EMAIL_DOMAIN ?= testdata.noreply.gitea.com
GITEA_SDK_SEEDED_ADMIN_USERNAME ?= gitea-admin
GITEA_SDK_SEEDED_ADMIN_PASSWORD ?= gitea-admin
GITEA_DAGGER_GO_IMAGE ?= docker.io/library/golang:1-alpine
GITEA_DAGGER_STABLE_TAG ?=
GITEA_DAGGER_MAIN_NIGHTLY_TAG ?= docker.gitea.com/gitea:main-nightly-rootless
GITEA_DAGGER_ADMIN_USER ?= $(GITEA_SDK_SEEDED_ADMIN_USERNAME)
GITEA_DAGGER_ADMIN_PASSWORD ?= $(GITEA_SDK_SEEDED_ADMIN_PASSWORD)
GITEA_DAGGER_STABLE_TAG_ARG := $(if $(GITEA_DAGGER_STABLE_TAG),--gitea-tag=$(GITEA_DAGGER_STABLE_TAG),)
GITEA_DAGGER_MAIN_NIGHTLY_TAG_ARG := $(if $(GITEA_DAGGER_MAIN_NIGHTLY_TAG),--gitea-tag=$(GITEA_DAGGER_MAIN_NIGHTLY_TAG),)
GITEA_DAGGER_ADMIN_ARG := $(if $(GITEA_DAGGER_ADMIN_USER),--admin-user=$(GITEA_DAGGER_ADMIN_USER),) $(if $(GITEA_DAGGER_ADMIN_PASSWORD),--admin-password=$(GITEA_DAGGER_ADMIN_PASSWORD),)

PACKAGE := gitea.dev/sdk

GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3
GITEA_VET_PACKAGE ?= code.gitea.io/gitea-vet@v0.2.1

GITEA_VERSION ?= 1.26-nightly
GITEA_DL := https://dl.gitea.com/gitea/$(GITEA_VERSION)/gitea-$(GITEA_VERSION)-
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
  GITEA_DL := $(GITEA_DL)linux-

  UNAME_P := $(shell uname -p)
  ifeq ($(UNAME_P),unknown)
   GITEA_DL := $(GITEA_DL)amd64
  endif
  ifeq ($(UNAME_P),x86_64)
   GITEA_DL := $(GITEA_DL)amd64
  endif
  ifneq ($(filter %86,$(UNAME_P)),)
   GITEA_DL := $(GITEA_DL)386
  endif
  ifneq ($(filter arm%,$(UNAME_P)),)
    GITEA_DL := $(GITEA_DL)arm-5
  endif
endif
ifeq ($(UNAME_S),Darwin)
  GITEA_DL := $(GITEA_DL)darwin-10.12-amd64
endif

.PHONY: all
all: clean test build

.PHONY: help
help:
	@echo "Make Routines:"
	@echo " - \"\"              run \"make clean test build\""
	@echo " - build             build sdk"
	@echo " - clean             clean"
	@echo " - fmt               format the code"
	@echo " - fmt-check         verify formatting"
	@echo " - generate          regenerate compatibility wrappers"
	@echo " - lint              run formatting checks and golangci-lint"
	@echo " - vet               examines Go source code and reports"
	@echo " - test-unit         run unit tests"
	@echo " - test-integration-stable  run integration tests (need a running gitea)"
	@echo " - test-integration-main-nightly run main-nightly integration tests"
	@echo " - test-unit-dagger  run unit tests with Dagger"
	@echo " - test-integration-dagger-stable  run stable integration tests with Dagger"
	@echo " - test-integration-dagger-main-nightly run main-nightly integration tests with Dagger"
	@echo " - test-integration-dagger run Dagger-backed integration tests"
	@echo " - test-dagger       run unit and integration tests with Dagger"
	@echo " - test              run unit and integration tests"
	@echo " - test-instance     start a gitea instance for test"


.PHONY: clean
clean:
	rm -r -f test
	$(GO) clean -i ./...

.PHONY: fmt
fmt:
	$(GO) run $(GOLANGCI_LINT_PACKAGE) fmt ./...

.PHONY: fmt-check
fmt-check:
	@before=$$(mktemp); after=$$(mktemp); \
	git diff --no-ext-diff --binary -- . > $$before; \
	$(GO) run $(GOLANGCI_LINT_PACKAGE) fmt ./...; \
	git diff --no-ext-diff --binary -- . > $$after; \
	if ! cmp -s $$before $$after; then \
		echo "Please run 'make fmt' and commit the result:"; \
		git diff --color=always -- .; \
		rm -f $$before $$after; \
		exit 1; \
	fi; \
	rm -f $$before $$after

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: vet
vet:
	# Default vet
	$(GO) vet $(PACKAGE)
	# Custom vet
	$(GO) get $(GITEA_VET_PACKAGE)
	$(GO) build code.gitea.io/gitea-vet
	$(GO) vet -vettool=gitea-vet $(PACKAGE)

.PHONY: lint
lint: fmt-check
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run --timeout 5m ./...

.PHONY: test-unit
test-unit:
	$(GO) test -race -cover -coverprofile coverage.out .

.PHONY: test-integration-stable
test-integration-stable:
	@export GITEA_SDK_TEST_URL=${GITEA_SDK_TEST_URL}; export GITEA_SDK_TEST_USERNAME=${GITEA_SDK_TEST_USERNAME}; export GITEA_SDK_TEST_PASSWORD=${GITEA_SDK_TEST_PASSWORD}; export GITEA_SDK_TEST_ADMIN_USERNAME=${GITEA_SDK_SEEDED_ADMIN_USERNAME}; export GITEA_SDK_TEST_ADMIN_PASSWORD=${GITEA_SDK_SEEDED_ADMIN_PASSWORD}; \
	if [ -z "$$(curl --noproxy "*" "${GITEA_SDK_TEST_URL}/api/v1/version" 2> /dev/null)" ]; then \echo "No test-instance detected!"; exit 1; else \
	    $(GO) test -race ./tests/... -skip '^TestMainNightly'; \
	fi

.PHONY: test-integration-main-nightly
test-integration-main-nightly:
	@export GITEA_SDK_TEST_URL=${GITEA_SDK_TEST_URL}; export GITEA_SDK_TEST_USERNAME=${GITEA_SDK_TEST_USERNAME}; export GITEA_SDK_TEST_PASSWORD=${GITEA_SDK_TEST_PASSWORD}; export GITEA_SDK_TEST_ADMIN_USERNAME=${GITEA_SDK_SEEDED_ADMIN_USERNAME}; export GITEA_SDK_TEST_ADMIN_PASSWORD=${GITEA_SDK_SEEDED_ADMIN_PASSWORD}; \
	if [ -z "$$(curl --noproxy "*" "${GITEA_SDK_TEST_URL}/api/v1/version" 2> /dev/null)" ]; then \echo "No test-instance detected!"; exit 1; else \
	    $(GO) test -race ./tests/... -run '^TestMainNightly'; \
	fi

.PHONY: test
test: test-unit test-integration-stable test-integration-main-nightly

.PHONY: test-unit-dagger
test-unit-dagger:
	$(DAGGER) -m ./dagger call unit --source=. --go-image=$(GITEA_DAGGER_GO_IMAGE)

.PHONY: test-integration-dagger-stable
test-integration-dagger-stable:
	$(DAGGER) -m ./dagger call integration-stable --source=. --go-image=$(GITEA_DAGGER_GO_IMAGE) $(GITEA_DAGGER_STABLE_TAG_ARG) $(GITEA_DAGGER_ADMIN_ARG)

.PHONY: test-integration-dagger-main-nightly
test-integration-dagger-main-nightly:
	$(DAGGER) -m ./dagger call integration-main-nightly --source=. --go-image=$(GITEA_DAGGER_GO_IMAGE) $(GITEA_DAGGER_MAIN_NIGHTLY_TAG_ARG) $(GITEA_DAGGER_ADMIN_ARG)

.PHONY: test-integration-dagger
test-integration-dagger: test-integration-dagger-stable test-integration-dagger-main-nightly

.PHONY: test-dagger
test-dagger: test-unit-dagger test-integration-dagger

.PHONY: test-instance
test-instance:
	rm -f -r ${WORK_DIR}/test 2> /dev/null; \
	mkdir -p ${WORK_DIR}/test/conf/ ${WORK_DIR}/test/data/
	wget ${GITEA_DL} -O ${WORK_DIR}/test/gitea-main; \
	chmod +x ${WORK_DIR}/test/gitea-main; \
	echo "[security]" > ${WORK_DIR}/test/conf/app.ini; \
	echo "INTERNAL_TOKEN = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NTg4MzY4ODB9.LoKQyK5TN_0kMJFVHWUW0uDAyoGjDP6Mkup4ps2VJN4" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "INSTALL_LOCK   = true" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "SECRET_KEY     = 2crAW4UANgvLipDS6U5obRcFosjSJHQANll6MNfX7P0G3se3fKcCwwK3szPyGcbo" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "PASSWORD_COMPLEXITY = off" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[database]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "DB_TYPE = sqlite3" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[repository]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "ROOT = ${WORK_DIR}/test/data/" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "[server]" >> ${WORK_DIR}/test/conf/app.ini; \
	echo "ROOT_URL = ${GITEA_SDK_TEST_URL}" >> ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/gitea-main migrate -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/gitea-main admin user create --username=${GITEA_SDK_SEEDED_ADMIN_USERNAME} --password=${GITEA_SDK_SEEDED_ADMIN_PASSWORD} --email=${GITEA_SDK_SEEDED_ADMIN_USERNAME}@${GITEA_SDK_TEST_EMAIL_DOMAIN} --admin=true --must-change-password=false -c ${WORK_DIR}/test/conf/app.ini; \
	${WORK_DIR}/test/gitea-main web -c ${WORK_DIR}/test/conf/app.ini

.PHONY: bench
bench:
	$(GO) test -run=XXXXXX -benchtime=10s -bench=. || exit 1

.PHONY: build
build:
	$(GO) build
