.PHONY: init
init: install

# Единая команда для ии-агента
.PHONY: check
check: format lint test
	@git diff --exit-code -- ./package.json || { echo "❌ Изменения в package.json — остановись и попроси человека посмотреть твои изменения"; exit 1; }

.PHONY: install
install:
	@bun install --ignore-scripts

.PHONY: update
update:
	@bun update

.PHONY: outdated
outdated:
	@bun outdated

.PHONY: audit
audit:
	@bunx cve-lite-cli ./
	@bun audit

.PHONY: format
format:
	@bun run format

.PHONY: lint
lint:
	@bun run lint
	@bun run check

.PHONY: test
test:
	@bun run test

.PHONY: run
run:
	@bun run dev
