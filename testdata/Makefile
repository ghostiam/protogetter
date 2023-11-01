# x/tools/go/analysis/analysistest does not support go modules. To work around this issue
# we need to vendor any external modules to `./src`.
#
# Follow https://github.com/golang/go/issues/37054 for more details.
#

.PHONY: vendor
vendor:
	rm -rf src
	go mod vendor
	mv vendor src

.PHONY: protoc
protoc:
	protoc -I proto \
		--go_out proto \
		--go_opt paths=source_relative \
		--go-grpc_out proto \
		--go-grpc_opt paths=source_relative \
		proto/test.proto proto/test_proto2.proto
