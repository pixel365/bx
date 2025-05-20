SPEC     := bx.spec
NAME     := $(shell rpm -q --qf "%{NAME}" --specfile $(SPEC))
VERSION  := $(shell rpm -q --qf "%{VERSION}" --specfile $(SPEC))
RELEASE  := $(shell rpm -q --qf "%{RELEASE}" --specfile $(SPEC))
DIST     := $(shell rpm --eval "%{?dist}")
SRPMDIR  := $(HOME)/rpmbuild/SRPMS
RPMSDIR  := $(HOME)/rpmbuild/RPMS
SPECDIR  := $(HOME)/rpmbuild/SPECS
SRPM     := $(SRPMDIR)/$(NAME)-$(VERSION)-$(RELEASE).src.rpm
RPM_ARCH := $(shell uname -m)
RPM      := $(RPMSDIR)/$(RPM_ARCH)/$(NAME)-$(VERSION)-$(RELEASE).$(RPM_ARCH).rpm

.PHONY: all fa fmt lint test build cover rpm srpm copr clean version

all: fa fmt lint test

fa:
	@fieldalignment -fix ./...

fmt:
	@goimports -w -local github.com/pixel365/bx .
	@gofmt -w .
	@golines -w .

lint:
	@golangci-lint run

test:
	@go $@ ./...

build:
	@go $@ -o ./bin/bx -ldflags="-s -w"

cover:
	go test -coverprofile=coverage.out ./... && go tool $@ -html=coverage.out

coverfn:
	go test -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out

doc:
	docsify serve docs


srpm:
	cp $(SPEC) $(SPECDIR)/ && \
	rpmbuild -bs $(SPECDIR)/$(SPEC)

rpm: srpm
	rpmbuild --rebuild $(SRPM)

copr: srpm
	copr-cli build $(NAME) $(SRPM)

clean:
	rm -f $(SRPM) $(RPM)

version:
	@echo Version: $(VERSION)-$(RELEASE)$(DIST)

