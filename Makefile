VERSION = 0.1.0

bin/pulumi-sdkgen-kubernetes-proxy:
	go build -o bin/pulumi-sdkgen-kubernetes-proxy ./cmd/pulumi-sdkgen-kubernetes-proxy

python-sdk: bin/pulumi-sdkgen-kubernetes-proxy
	rm -rf sdk
	bin/pulumi-sdkgen-kubernetes-proxy schema.json sdk
	cp README.md sdk/python/
	cd sdk/python/ && \
		sed -i.bak -e "s/\$${VERSION}/$(VERSION)/g" -e "s/\$${PLUGIN_VERSION}/$(VERSION)/g" setup.py && \
		rm setup.py.bak
