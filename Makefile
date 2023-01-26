# Image URL to use all building/pushing image targets

# renovate: datasource=github-releases depName=cert-manager/cert-manager
CERT_MANAGER_VERSION ?= v1.11.0
TAG ?= "$(shell date +%Y%m%d%s)"
TAG := $(TAG)

# renovate: datasource=github-tags depName=helm/helm
HELM_VERSION ?= v3.10.2
CHART_VERSION = v0.5.0 # x-release-please-version

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn
ARCH?=amd64

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
# renovate: datasource=github-tags depName=kubernetes-sigs/kustomize
KUSTOMIZE_VERSION?=v4.5.7

.PHONY: integration-test #these tests should run on a real cluster!
integration-test:
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test.yaml

.PHONY: integration-test-local #these tests should run on a real cluster!
integration-test-local:
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test-local.yaml

.PHONY: cleanup-manifests
cleanup-manifests:
	rm -rf manifests

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: release-helm-manifests
release-helm-manifests: kustomize
	echo "building helm overlay"
	$(KUSTOMIZE) build ./helm/overlay > ./helm/chart/templates/rendered.yaml

.PHONY: helm-package
helm-package: build-release-manifests release-helm-manifests
	cd ./helm && helm package ./chart
	cd ./helm && mv keptn-lifecycle-toolkit-*.tgz ./chart/charts

.PHONY: build-release-manifests
build-release-manifests:
	$(MAKE) -C operator generate
	$(MAKE) -C klt-cert-manager generate
	$(MAKE) -C operator release-helm-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)
	$(MAKE) -C scheduler release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)
	$(MAKE) -C klt-cert-manager release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)

.PHONY: build-deploy-operator
build-deploy-operator:
	$(MAKE) -C operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)

	kubectl apply -f operator/config/rendered/release.yaml

.PHONY: build-deploy-scheduler
build-deploy-scheduler:
	$(MAKE) -C scheduler release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)
	kubectl create namespace keptn-lifecycle-toolkit-system --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f scheduler/config/rendered/release.yaml

.PHONY: build-deploy-certmanager
build-deploy-certmanager:
	$(MAKE) -C klt-cert-manager release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C klt-cert-manager push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C klt-cert-manager release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)
	kubectl create namespace keptn-lifecycle-toolkit-system --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f klt-cert-manager/config/rendered/release.yaml

.PHONY: build-deploy-dev-environment
build-deploy-dev-environment: build-deploy-certmanager build-deploy-operator build-deploy-scheduler
