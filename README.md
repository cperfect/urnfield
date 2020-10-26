Library for using URN fields in go structs or func params
====================================================

WARNING WIP - use at your own risk

See https://tools.ietf.org/html/rfc8141

The purpose of this is to allow identification across contexts where the actual data to map to isn't available in full and to support lazy loading.

The intent is that once a URN string is created and set then it is immutable (except mayb the resolvers component?)