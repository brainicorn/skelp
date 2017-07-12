package main

//go:generate rm ./skelplate/schema_accessor.go
//go:generate jsonschemagen -x -c -f schema.json -o ./skelplate github.com/brainicorn/skelp/skelplate SkelplateDescriptor
//go:generate go run ./md/markdown.go ./md/
