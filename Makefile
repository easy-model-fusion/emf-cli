APP ?= emf-cli
9 ?=
TAG ?= v1.0.0
BUILD ?= 0
BUILD_DATE = $(shell date +%FT%T)

GO111MODULE ?= on
GOPROXY ?= "https://proxy.golang.org,direct"
GOSUMDB ?= "sum.golang.org"
CGO_ENABLED ?= 0
GOINSECURE ?=
GONOSUMDB ?=
GO_OPT=GOPROXY=$(GOPROXY) GOINSECURE=$(GOINSECURE) GONOSUMDB=$(GONOSUMDB) GOSUMDB=$(GOSUMDB) GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS)

test-unit: # Test the code
	mkdir -p reporting
	$(GO_OPT) go test -p=1 -short -cover -coverpkg=$$($(GO_PACKAGE) | tr '\n' ',') -coverprofile=reporting/profile.out -json $$($(GO_PACKAGE)) > reporting/tests.json  || true
	go tool cover -html=reporting/profile.out -o reporting/coverage.html
	go tool cover -func=reporting/profile.out -o reporting/coverage.txt
	cat reporting/coverage.txt

build: # Build the executable
	$(GO_OPT) go build -a -trimpath -ldflags "-X main.Version=$(TAG)-$(BUILD) -X main.BuildDate=$(BUILD_DATE)" -o bin/$(APP)

run: # Run the executable
	bin/$(APP)