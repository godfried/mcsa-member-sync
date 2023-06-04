-include PrivateRules.mak

MEMBAZ_EXPORT=$(shell pwd)/membaz-export.csv
EVERLYTIC_EXPORT=$(shell pwd)/everlytic-export.csv
MEMBAZ_MISSING=$(shell pwd)/membaz-missing.csv
EVERLYTIC_MISSING=$(shell pwd)/everlytic-missing.csv

build:
	mkdir -p bin
	go build -v -o ./bin ./...

clean:
	rm -rf bin

run: membaz-export everlytic-export find-missing

membaz-export: build
	bin/membaz-export -password "$(MEMBAZ_PASSWORD)"\
                  -username "$(MEMBAZ_USERNAME)"\
                  -destination $(MEMBAZ_EXPORT)

everlytic-export: build
	bin/everlytic-export -api-key "$(EVERLYTIC_API_KEY)"\
                     -username "$(EVERLYTIC_USERNAME)"\
                     -destination $(EVERLYTIC_EXPORT)

find-missing: build
	bin/find-missing -membaz-destination $(MEMBAZ_MISSING)\
                 -everlytic-destination $(EVERLYTIC_MISSING)\
                 -membaz-csv $(MEMBAZ_EXPORT)\
                 -everlytic-csv $(EVERLYTIC_EXPORT)

test:
	go test -v ./...