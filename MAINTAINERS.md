# Maintainer instructions

To release a new version:

1. Update the version in `Makefile` and
   `cmd/pulumi-resource-kubernetes-proxy/main.go`.

2. Run `make python-sdk`.

3. Commit the changes and push to GitHub.

3. Tag the new version with `git tag -a $version -m $version`. Push the tag to
   GitHub.
