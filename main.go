package main

import (
	"github.com/charlesbases/protoc-gen-swagger/protoc"
	"github.com/charlesbases/protoc-gen-swagger/swagger"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protoc.Plugin(func(p *protoc.Package) *pluginpb.CodeGeneratorResponse {
		var rsp = new(pluginpb.CodeGeneratorResponse)

		// swagger api
		rsp.File = append(rsp.File, swagger.New(p).Generater())

		return rsp
	})
}
