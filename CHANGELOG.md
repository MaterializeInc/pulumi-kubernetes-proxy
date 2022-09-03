# Changelog

All notable changes to this crate will be documented in this file.

The format is based on [Keep a Changelog], and this crate adheres to [Semantic
Versioning].

## 0.2.0 - 2022-09-06

* Backwards incompatible change: creating the provider is no longer sufficient
  to create the proxy listener. Instead, you must call the `startproxy`
  provider function, and use the port returned in the output in the rest of
  your calls. The returned port will be the port given in the `host_port`
  argument to the provider constructor, but using the `startproxy` return value
  will ensure pulumi creates a dependency such that anything using the
  function's output will have to happen after the proxy is fully initialized.

## 0.1.3 - 2021-07-06

* Correct version format in package URL.

## 0.1.2 - 2021-07-06

* Adjust release process to avoid hardcoding versions.

## 0.1.1 - 2021-07-06

* Fix naming of release tarballs so that `pulumi plugin install` works
  correctly.

## 0.1.0 - 2021-07-06

Initial release.

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
