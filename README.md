merrors
=======

Simple multi-error `error` type for Go.

tl;dr:

* Sometimes you want multiple errors to be bundled into a single error.
* Whereas other libraries directly act on the error object, this package uses a builder to create errors. This significanly reduces the complexity of the error object itself.

# DESCRIPTION

<!-- INCLUDE(merror_example_test.go) -->
<!-- END INCLUDE -->
