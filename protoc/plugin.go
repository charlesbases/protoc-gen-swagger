package protoc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/charlesbases/protoc-gen-swagger/conf"
	"github.com/charlesbases/protoc-gen-swagger/logger"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin .
func Plugin(fn func(p *Package) *pluginpb.CodeGeneratorResponse) {
	var buff = new(bytes.Buffer)
	if _, err := io.Copy(buff, os.Stdin); err != nil {
		logger.Fatal("read os.Stdin failed. ", err)
	}

	var req = new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(buff.Bytes(), req); err != nil {
		logger.Fatal("unmarshal os.Stdin failed. ", err)
	}
	if len(req.GetFileToGenerate()) == 0 {
		logger.Fatal("no file to generate")
	}

	parseArgs(req)

	if rsp, err := proto.Marshal(fn(parse(req))); err != nil {
		logger.Fatal(err)
	} else {
		os.Stdout.Write(rsp)
	}
}

// parseArgs 加载 protoc 传入的参数
func parseArgs(req *pluginpb.CodeGeneratorRequest) {
	for _, param := range strings.Split(req.GetParameter(), ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}

		switch param {
		// 解析基础配置文件
		case "confdir":
			conf.Parse(value)
		}
	}
}

// parse .
func parse(req *pluginpb.CodeGeneratorRequest) *Package {
	mark := Mark{Name: req.GetProtoFile()[0].GetPackage()}
	var p = mark.newPackage()

	var swg = sync.WaitGroup{}
	swg.Add(len(req.GetProtoFile()))

	for fidx := range req.GetProtoFile() {
		go func(file *descriptorpb.FileDescriptorProto) {
			// parse comment
			var cs = parseComments(file.SourceCodeInfo)

			// parse message
			for midx, protoMessage := range file.GetMessageType() {
				var paths = []int{COMMENT_PATH_MESSAGE, midx}

				for nidx, protoNested := range protoMessage.GetNestedType() {
					p.Messages = append(p.Messages, cs.parseMessageNested(protoNested, protoMessage.GetName(), append(paths, COMMENT_PATH_MESSAGE_MESSAGE, nidx)...))
				}

				for eidx, protoEnum := range protoMessage.GetEnumType() {
					p.Enums = append(p.Enums, cs.parseMessageEnum(protoEnum, protoMessage.GetName(), append(paths, COMMENT_PATH_MESSAGE_ENUM, eidx)...))
				}

				p.Messages = append(p.Messages, cs.parseMessage(protoMessage, paths...))
			}

			// parse enum
			for idx, protoEnum := range file.GetEnumType() {
				p.Enums = append(p.Enums, cs.parseEnum(protoEnum, COMMENT_PATH_ENUM, idx))
			}

			// parse service
			for idx, protoService := range file.GetService() {
				p.Services = append(p.Services, cs.parseService(protoService, COMMENT_PATH_SERVICE, idx))
			}

			swg.Done()
		}(req.GetProtoFile()[fidx])
	}

	swg.Wait()

	return p.sort()
}

// parseComments paarse comments in proto
func parseComments(infor *descriptorpb.SourceCodeInfo) comments {
	cs := make(map[string]*comment, 0)

	for _, location := range infor.GetLocation() {
		if location.GetLeadingComments() == "" && location.GetTrailingComments() == "" && len(location.GetLeadingDetachedComments()) == 0 {
			continue
		}

		detached := make([]string, 0)
		for _, val := range location.GetLeadingDetachedComments() {
			detached = append(detached, trim(val, "*", "\n"))
		}

		cs[fmt.Sprintf("%v", location.GetPath())] = &comment{
			leading:  trim(location.GetLeadingComments(), "*", "\n"),
			trailing: trim(location.GetTrailingComments(), "*", "\n"),
			detached: detached,
		}
	}
	return cs
}

// parseservice parse service in proto
func (cs comments) parseService(protoService *descriptorpb.ServiceDescriptorProto, paths ...int) *Service {
	mark := Mark{Name: protoService.GetName(), Desc: cs.comment(protoService.GetName(), paths...)}
	var service = mark.newService()

	for idx, protoRPC := range protoService.GetMethod() {
		method := cs.parseMethod(protoRPC, append(paths, COMMENT_PATH_SERVICE_METHOD, idx)...)
		if len(method.Path) == 0 {
			method.Path = methodPath(service.Name, method.Name)
		}
		service.Methods = append(service.Methods, method)
	}
	return service
}

// parseMethod parse method in service
func (cs comments) parseMethod(protoRPC *descriptorpb.MethodDescriptorProto, paths ...int) *ServiceMethod {
	var method = cs.mark(protoRPC.GetName(), paths...).newServiceMethod()

	method.RequestName = split(protoRPC.GetInputType())[1]
	method.ResponseName = split(protoRPC.GetOutputType())[1]
	return method
}

// parseMessage parse message in proto
func (cs comments) parseMessage(protoMessage *descriptorpb.DescriptorProto, paths ...int) *Message {
	mark := Mark{Name: protoMessage.GetName(), Desc: cs.comment(protoMessage.GetName(), paths...)}
	var message = mark.newMessage()

	for idx, field := range protoMessage.GetField() {
		message.Fields = append(message.Fields, cs.parseMessageField(protoMessage, field, append(paths, COMMENT_PATH_MESSAGE_FIELD, idx)...))
	}
	return message
}

// parseMessageNested parse message nested in message
func (cs comments) parseMessageNested(nested *descriptorpb.DescriptorProto, parent string, paths ...int) *Message {
	name := nestedName(parent, nested.GetName())
	mark := Mark{Name: name, Desc: cs.comment(name, paths...)}
	var message = mark.newMessage()

	for idx, field := range nested.GetField() {
		message.Fields = append(message.Fields, cs.parseMessageField(nested, field, append(paths, COMMENT_PATH_MESSAGE_FIELD, idx)...))
	}
	return message
}

// parseMessageEnum parse enum in message
func (cs comments) parseMessageEnum(protoEnum *descriptorpb.EnumDescriptorProto, parent string, paths ...int) *Enum {
	name := nestedName(parent, protoEnum.GetName())
	mark := Mark{Name: name, Desc: cs.comment(name, paths...)}
	var enum = mark.newEnum()

	for idx, enumField := range protoEnum.GetValue() {
		enum.Fields = append(enum.Fields, cs.parseEnumField(enumField, append(paths, COMMENT_PATH_ENUM_VALUE, idx)...))
	}
	return enum
}

// parseMessageField parse field in message
func (cs comments) parseMessageField(protoMessage *descriptorpb.DescriptorProto, protoField *descriptorpb.FieldDescriptorProto, paths ...int) *MessageField {
	var field = &MessageField{MessageName: protoMessage.GetName(), Description: cs.comment(protoField.GetName(), paths...)}

	// Json
	field.JsonName = protoField.GetName()
	field.JsonLabel = protoLabel2JsonLabel[protoField.GetLabel()]
	field.JsonType = protoType2JsonType[protoField.GetType()]
	field.JsonDefaultValue = jsonTypeDefaultValue[field.JsonType]

	// Proto
	field.ProtoName = protoField.GetName()
	field.ProtoLaber = protoField.GetLabel()
	field.ProtoType = protoField.GetType()
	field.ProtoNumber = protoField.GetNumber()

	switch field.JsonType {
	case JSON_TYPE_OBJECT:
		typename := split(protoField.GetTypeName())

		field.JsonType = typename[1]
		field.ProtoTypeName = typename[1]
		field.ProtoPackagePath = typename[0]
		field.ProtoShortName = protoField.GetTypeName()
	case JSON_TPYE_NUMBER, JSON_TYPE_STRING, JSON_TYPE_BOOLEAN:
		field.ProtoTypeName = descriptorpb.FieldDescriptorProto_Type_name[int32(field.ProtoType)]
	}

	return field
}

// parseEnum parse enum in proto
func (cs comments) parseEnum(protoEnum *descriptorpb.EnumDescriptorProto, paths ...int) *Enum {
	mark := Mark{Name: protoEnum.GetName(), Desc: cs.comment(protoEnum.GetName(), paths...)}
	var enum = mark.newEnum()

	for idx, enumField := range protoEnum.GetValue() {
		enum.Fields = append(enum.Fields, cs.parseEnumField(enumField, append(paths, COMMENT_PATH_ENUM_VALUE, idx)...))
	}
	return enum
}

// parseEnumField parse field in enum
func (cs comments) parseEnumField(protoEnumField *descriptorpb.EnumValueDescriptorProto, paths ...int) *EnumField {
	return &EnumField{
		Name:        protoEnumField.GetName(),
		Value:       protoEnumField.GetNumber(),
		Description: cs.comment(protoEnumField.GetName(), paths...),
	}
}
