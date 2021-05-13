run:
	go run ./cmd/idrac


mod:
# real tab space or get error Makefile:2: *** missing separator.  Stop.
	go mod tidy
	go mod vendor
