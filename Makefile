# Usage:
#   make test              race tests, all modules (-count=50)
#   make test M=05         race tests for one module (05-bounded-queue)
#   make bench M=05        benchmarks for one module
#   make bug N=04          run one bug-hunt drill
#   make progress          checked roadmap items per tier

M     ?=
N     ?=
COUNT ?= 50

PKGS := $(if $(M),./$(wildcard $(M)-*)/...,./...)
BUGS := ./interview/drills/bughunt
BUG  := $(wildcard interview/drills/bughunt/$(N)-*)

.PHONY: test test-fast bench vet lint fmt cover progress bug bugs bug-reset ci

test:
	go test -race -count=$(COUNT) -timeout 10m $(PKGS)

test-fast:
	go test -race -count=1 $(PKGS)

bench:
	go test -run='^$$' -bench=. -benchmem -cpu=1,2,4,8 $(PKGS)

vet:
	go vet $(PKGS)

lint:
	golangci-lint run $(PKGS)

fmt:
	gofmt -s -w .

cover:
	go test -race -coverprofile=cover.out $(PKGS) && go tool cover -func=cover.out | tail -1

progress:
	go run ./tools/progress

bug:
	go test -tags bughunt -race -count=1 -timeout 60s ./$(BUG)/

bugs:
	-go test -tags bughunt -race -count=1 -timeout 60s $(BUGS)/...

bug-reset:
	git checkout -- $(BUGS)

ci: vet lint test-fast
