GO ?= GO111MODULE=on CGO_ENABLED=0 go

run:
	$(GO) run ${TARGET}

build:
	$(GO) build -o bootstrap ${TARGET}

test:
	$(GO) test ./... -count=1 -p=100

update-all:
	$(GO) get -u ./...

zip:
	rm .deploy.zip ; zip -r .deploy.zip **