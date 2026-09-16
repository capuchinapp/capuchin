.PHONY: require-version
require-version:
	if [ -n ${APP_VERSION} ] && [[ $$APP_VERSION == "v"* ]]; then echo "The app version may not start with v" && exit 1; fi
	if [ -z ${APP_VERSION} ]; then echo "Need to set APP_VERSION" && exit 1; fi;

.PHONY: build-docker-capuchin
build-docker-capuchin: require-version
	@docker buildx build --sbom=true --provenance=true \
		$(if $(PUSH),--push) \
		--label org.opencontainers.image.revision=$(shell git rev-parse HEAD) \
		--label org.opencontainers.image.version=$(APP_VERSION) \
		--label org.opencontainers.image.created=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) \
		-t ghcr.io/capuchinapp/capuchin:latest \
		-t ghcr.io/capuchinapp/capuchin:$(APP_VERSION) \
		-t ghcr.io/capuchinapp/capuchin:$(shell echo $(APP_VERSION) | cut -d '.' -f -2) \
		--build-arg APP_VERSION=$(APP_VERSION) \
		-f ./build/Dockerfile .

.PHONY: release
release:
	@./scripts/release.sh $(if $(DRY_RUN),--dry-run)
