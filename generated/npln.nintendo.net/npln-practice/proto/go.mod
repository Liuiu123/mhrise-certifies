module npln.nintendo.net/npln-practice/proto

go 1.27.1

replace github.com/lyft/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.2.1

require (
    github.com/envoyproxy/protoc-gen-validate v1.2.1
    google.golang.org/genproto/googleapis/api v0.0.0-20260904163448-b1c236e22ff4
    google.golang.org/protobuf v1.36.12
)