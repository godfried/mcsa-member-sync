-include PrivateRules.mak

MEMBAZ_EXPORT := $(shell pwd)/membaz-export.csv
EVERLYTIC_EXPORT := $(shell pwd)/everlytic-export.csv
MEMBAZ_MISSING := $(shell pwd)/membaz-missing.csv
EVERLYTIC_MISSING := $(shell pwd)/everlytic-missing.csv
EVERLYTIC_API_KEY ?= NOT_SET
EVERLYTIC_USERNAME ?= NOT_SET
MEMBAZ_PASSWORD ?= NOT_SET
MEMBAZ_USERNAME ?= NOT_SET

build: ## Build executables
	mkdir -p bin
	go build -v -o ./bin ./...
.PHONY: build

clean: ## Remove build artefacts
	rm -rf bin
.PHONY: clean

run: membaz-export everlytic-export find-missing ## Run full chain to add missing members to Everlytic & remove them from Membaz.
.PHONY: run

run-github: run
	if [ ! -s $(MEMBAZ_MISSING) ] ; then rm $(MEMBAZ_MISSING); fi
	if [ ! -s $(EVERLYTIC_MISSING) ] ; then rm $(EVERLYTIC_MISSING); fi
.PHONY: run-github

membaz-export: build ## Export members from Membaz
	bin/membaz-export -password "$(MEMBAZ_PASSWORD)"\
                  -username "$(MEMBAZ_USERNAME)"\
                  -destination $(MEMBAZ_EXPORT)
.PHONY: membaz-export

everlytic-export: build ## Export members from Everlytic
	bin/everlytic-export -api-key "$(EVERLYTIC_API_KEY)"\
                     -username "$(EVERLYTIC_USERNAME)"\
                     -destination $(EVERLYTIC_EXPORT)
.PHONY: everlytic-export

find-missing: build ## Find missing members in Membaz & Everlytic
	bin/find-missing -membaz-destination $(MEMBAZ_MISSING)\
                 -everlytic-destination $(EVERLYTIC_MISSING)\
                 -membaz-csv $(MEMBAZ_EXPORT)\
                 -everlytic-csv $(EVERLYTIC_EXPORT)
.PHONY: find-missing

everlytic-unsubscribe: build ## Unsubscribe members who are not present in Membaz
	bin/everlytic-unsubscribe -source $(MEMBAZ_MISSING)\
                              -api-key "$(EVERLYTIC_API_KEY)"\
                              -username "$(EVERLYTIC_USERNAME)"
.PHONY: everlytic-unsubscribe

test: ## Execute tests
	go test -v ./...
.PHONY: test

help:  ## Show this help.
	@echo "make targets:"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ": .*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""; echo "make vars (+defaults):"
	@grep -E '^[0-9a-zA-Z_-]+ \:=.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = " \\:= "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@grep -E '^[0-9a-zA-Z_-]+ \?=.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = " \\?= "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

vars:  ## Variables
	@echo "Current variable settings:"
	@echo "MEMBAZ_EXPORT=$(MEMBAZ_EXPORT)"
	@echo "EVERLYTIC_EXPORT=$(EVERLYTIC_EXPORT)"
	@echo "MEMBAZ_MISSING=$(MEMBAZ_MISSING)"
	@echo "EVERLYTIC_MISSING=$(EVERLYTIC_MISSING)"
	@echo "EVERLYTIC_USERNAME=$(EVERLYTIC_USERNAME)"
	@echo "EVERLYTIC_API_KEY=$(EVERLYTIC_API_KEY)"
	@echo "MEMBAZ_USERNAME=$(MEMBAZ_USERNAME)"
	@echo "MEMBAZ_PASSWORD=$(MEMBAZ_PASSWORD)"
.PHONY: vars
