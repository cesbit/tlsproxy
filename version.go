package main

// Version exposes the go-thingsdb version
const Version = "0.1.0"

//Publish module:
//
//  go mod tidy
//  go test ./...
//  git tag {VERSION}
//  git push origin {VERSION}
//  GOPROXY=proxy.golang.org go list -m github.com/thingsdb/go-thingsdb@{VERSION}
