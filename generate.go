package main

//go:generate jsonschemagen -x -c -f schema.json -o ./skelplate github.com/brainicorn/skelp/skelplate SkelplateDescriptor
///go:generate go-bindata -pkg skelplate -o ./skelplate/bindata.go ./skelplate/data
