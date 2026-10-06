module mhrise-npln-lab

go 1.27.1

require (
	github.com/golang/protobuf v1.5.4
	golang.org/x/net v0.57.0
	golang.org/x/sys v0.47.0
	golang.org/x/text v0.40.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260825221802-da73d73af1c5
	google.golang.org/grpc v1.83.1
	google.golang.org/protobuf v1.36.12
)

require (
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/lyft/protoc-gen-star/v2 v2.0.4 // indirect
	github.com/lyft/protoc-gen-validate v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/tools v0.47.0 // indirect
	google.golang.org/genproto v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260904163448-b1c236e22ff4 // indirect
	npln.nintendo.net/npln-practice/proto v0.0.0 // indirect
)

replace github.com/lyft/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.2.1

replace npln.nintendo.net/npln-practice/proto => ./generated/npln.nintendo.net/npln-practice/proto
