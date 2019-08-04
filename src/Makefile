BUILDPATH=$(CURDIR)
GO=$(shell which go)
GOINSTALL=$(GO) install
GOLINT=$(shell which golint)
GOFMT=$(GO) fmt
GOVET=$(GO) vet
build-lb:
	@echo "Building Distributed Load Balancer"
	@cd  $(BUILDPATH)/src && $(GOINSTALL)
	@echo "done"
clean:
	@echo "Cleaning binary"
	@rm -rf  $(BUILDPATH)/bin/*
	@echo "done"
fmt-lb:
	@echo "running  fmt on Distributed Load Balancer"
	@cd  $(BUILDPATH)/src/lb && $(GOFMT)
	@echo "done"
vet-lb:
	@echo "running  vet on Distributed Load Balancer"
	@cd  $(BUILDPATH)/src/lb && $(GOVET)
	@echo "done"
lint-lb:
	@echo "running  vet on Distributed Load Balancer"
	@cd  $(BUILDPATH)/src/lb && $(GOLINT)
	@echo "done"
all: clean fmt-lb vet-lb lint-lb build-lb