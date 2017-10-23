package main

//go:generate rm -f ./skelplate/schema_accessor.go
//go:generate jsonschemagen -xcf schema.json -o ./skelplate github.com/brainicorn/skelp/skelplate SkelplateDescriptor
//go:generate go run ./md/markdown.go ./md/
