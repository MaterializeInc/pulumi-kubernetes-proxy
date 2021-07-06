VERSION ?= $(patsubst v%,%,$(shell git describe))

bin/pulumi-sdkgen-kubernetes-proxy: cmd/pulumi-sdkgen-kubernetes-proxy/*.go
	go build -o bin/pulumi-sdkgen-kubernetes-proxy ./cmd/pulumi-sdkgen-kubernetes-proxy

python-sdk: bin/pulumi-sdkgen-kubernetes-proxy
	rm -rf sdk
	bin/pulumi-sdkgen-kubernetes-proxy $(VERSION)
	cp README.md sdk/python/
	cd sdk/python/ && \
		sed -i.bak -e "s/\$${VERSION}/$(VERSION)/g" -e "s/\$${PLUGIN_VERSION}/$(VERSION)/g" setup.py && \
		rm setup.py.bak
