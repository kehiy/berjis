check:
	golangci-lint run \
		--timeout=20m0s \
		--exclude='unused' \
		--enable=gofmt \
		--enable=unconvert \
		--enable=asciicheck \
		--enable=misspell \
		--enable=revive \
		--enable=decorder \
		--enable=reassign \
		--enable=usestdlibvars \
		--enable=nilerr \
		--enable=gosec \
		--enable=exportloopref \
		--enable=whitespace \
		--enable=goimports \
		--enable=gocyclo \
		--enable=lll
