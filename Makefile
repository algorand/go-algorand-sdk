lint:
	golint `go list ./... | grep -v /vendor/`
