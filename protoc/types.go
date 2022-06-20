package protoc

import (
	"sort"
	"sync"

	"google.golang.org/protobuf/types/descriptorpb"
)

// json 字段标签
const (
	JSON_LABEL_OPTIONAL = "可选"
	JSON_LABEL_REQUIRED = "必须"
	JSON_LABEL_REPEATED = "重复"
)

// json 字段类型
const (
	JSON_TPYE_NUMBER  = `Number`
	JSON_TYPE_STRING  = `String`
	JSON_TYPE_BOOLEAN = `Boolean`
	JSON_TYPE_OBJECT  = `Object`
)

type Method string

const (
	MethodGet    Method = "GET"
	MethodPut    Method = "PUT"
	MethodPost   Method = "POST"
	MethodDelete Method = "DELETE"
)

var methods = map[Method]string{
	MethodGet:    "get",
	MethodPut:    "put",
	MethodPost:   "post",
	MethodDelete: "delete",
}

// String .
func (m Method) String() string {
	return string(m)
}

// LowerCase .
func (m Method) LowerCase() string {
	return methods[m]
}

var (
	jsonTypeDefaultValue = map[string]interface{}{
		JSON_TPYE_NUMBER:  0,
		JSON_TYPE_STRING:  "string",
		JSON_TYPE_BOOLEAN: false,
		JSON_TYPE_OBJECT:  nil,
	}

	protoType2JsonType = map[descriptorpb.FieldDescriptorProto_Type]string{
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:   JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT:    JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_INT32:    JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:   JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_INT64:    JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:   JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32:  JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:  JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32:   JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:   JSON_TPYE_NUMBER,
		descriptorpb.FieldDescriptorProto_TYPE_BYTES:    JSON_TYPE_STRING,
		descriptorpb.FieldDescriptorProto_TYPE_STRING:   JSON_TYPE_STRING,
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:     JSON_TYPE_BOOLEAN,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:     JSON_TYPE_OBJECT,
		descriptorpb.FieldDescriptorProto_TYPE_GROUP:    JSON_TYPE_OBJECT,
		descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:  JSON_TYPE_OBJECT,
	}

	protoLabel2JsonLabel = map[descriptorpb.FieldDescriptorProto_Label]string{
		descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL: JSON_LABEL_OPTIONAL,
		descriptorpb.FieldDescriptorProto_LABEL_REQUIRED: JSON_LABEL_REQUIRED,
		descriptorpb.FieldDescriptorProto_LABEL_REPEATED: JSON_LABEL_REPEATED,
	}
)

const ContentTypeJson = "application/json"

type (
	Package struct {
		enumLocker sync.RWMutex
		messLocker sync.RWMutex

		// Name Package.Name
		Name string
		// Version version
		Version string
		// Prefix uri prefix
		Prefix string
		// Services Service list
		Services []*Service
		// Enums Enum list
		Enums []*Enum
		// EnumDic Enum map
		EnumDic map[string]*Enum
		// Messages Message list
		Messages []*Message
		// MessageDic Message map
		MessageDic map[string]*Message
	}

	Service struct {
		Name        string
		Description string
		// Methods rpc list
		Methods []*ServiceMethod
	}

	// ServiceMethod service.rpc
	ServiceMethod struct {
		Name         string
		Path         string
		Method       Method
		Description  string
		Consume      string
		Produce      string
		RequestName  string
		ResponseName string
	}

	Enum struct {
		Name        string
		Description string
		Fields      []*EnumField
	}

	EnumField struct {
		Name        string
		Value       int32
		Description string
	}

	Message struct {
		Name        string
		Description string
		Fields      []*MessageField
	}

	MessageField struct {
		// MessageName Message.Name
		MessageName string
		// Description field description
		Description string

		ProtoName        string                                  // proto field name
		ProtoLaber       descriptorpb.FieldDescriptorProto_Label // proto 标签
		ProtoType        descriptorpb.FieldDescriptorProto_Type  // 隐式类型
		ProtoTypeName    string                                  // 显示类型
		ProtoFullName    string                                  // 包名.结构名
		ProtoPackagePath string                                  // 包路径
		ProtoNumber      int32                                   // 排序

		JsonName         string      // json field name
		JsonLabel        string      // json 标签
		JsonType         string      // json 类型
		JsonDefaultValue interface{} // json 数据默认值
	}
)

// sort .
func (p *Package) sort() *Package {
	var swg = sync.WaitGroup{}
	swg.Add(3)

	// Services
	go func() {
		if len(p.Services) != 0 {
			sort.Slice(p.Services, func(i, j int) bool {
				return ascending(p.Services[i].Name, p.Services[j].Name)
			})
		}

		swg.Done()
	}()

	// Messages
	go func() {
		if len(p.Messages) != 0 {
			sort.Slice(p.Messages, func(i, j int) bool {
				return ascending(p.Messages[i].Name, p.Messages[j].Name)
			})
		}

		swg.Done()
	}()

	// Enums
	go func() {
		if len(p.Enums) != 0 {
			sort.Slice(p.Enums, func(i, j int) bool {
				return ascending(p.Enums[i].Name, p.Enums[j].Name)
			})
		}

		swg.Done()
	}()

	swg.Wait()
	return p
}
