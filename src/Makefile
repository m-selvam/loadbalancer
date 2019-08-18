BUILDPATH=$(CURDIR)
GO=$(shell which go)
GOINSTALL=$(GO) install
GOLINT=$(shell which golint)
GOFMT=$(GO) fmt
GOVET=$(GO) vet
build-lb:
	@echo "Building Distributed Load Balancer"
	@cd  $(BUILDPATH)/lb && $(GOINSTALL)
	@echo "done"
clean:
	@echo "Cleaning binary"
	@rm -rf  $(BUILDPATH)/pkg/*
	@echo "done"
fmt-lb:
	@echo "running  fmt on Distributed Load Balancer"
	@cd  $(BUILDPATH)/lb && $(GOFMT)
	@echo "done"
vet-lb:
	@echo "running  vet on Distributed Load Balancer"
	@cd  $(BUILDPATH)/lb && $(GOVET)
	@echo "done"
lint-lb:
	@echo "running  vet on Distributed Load Balancer"
	@cd  $(BUILDPATH)/lb && $(GOLINT)
	@echo "done"
all: clean fmt-lb vet-lb lint-lb build-lb