package protoc

import (
	"fmt"
)

// comment path
const (
	// tag numbers in FileDescriptorProto

	// COMMENT_PATH_PACKAGE package comment
	COMMENT_PATH_PACKAGE = 2
	// COMMENT_PATH_MESSAGE message comment
	COMMENT_PATH_MESSAGE = 4
	// COMMENT_PATH_ENUM enum comment
	COMMENT_PATH_ENUM = 5
	// COMMENT_PATH_SERVICE service comment
	COMMENT_PATH_SERVICE = 6
	// COMMENT_PATH_EXTENSION extension comment
	COMMENT_PATH_EXTENSION = 7
	// COMMENT_PATH_SYNTAX syntax comment
	COMMENT_PATH_SYNTAX = 12

	// tag numbers in DescriptorProto

	// COMMENT_PATH_MESSAGE_FIELD message.field
	COMMENT_PATH_MESSAGE_FIELD = 2
	// COMMENT_PATH_MESSAGE_MESSAGE message.nested
	COMMENT_PATH_MESSAGE_MESSAGE = 3
	// COMMENT_PATH_MESSAGE_ENUM message.enum
	COMMENT_PATH_MESSAGE_ENUM = 4
	// COMMENT_PATH_MESSAGE_EXTENSION message.ectension
	COMMENT_PATH_MESSAGE_EXTENSION = 6

	// tag numbers in EnumDescriptorProto

	// COMMENT_PATH_ENUM_VALUE enum value
	COMMENT_PATH_ENUM_VALUE = 2

	// tag numbers in ServiceDescriptorProto

	// COMMENT_PATH_SERVICE_METHOD service method
	COMMENT_PATH_SERVICE_METHOD = 2
)

type (
	comments map[string]*comment

	comment struct {
		leading  string
		trailing string
		detached []string
	}
)

// comment get comment by path
func (cs comments) comment(name string, paths ...int) string {
	if comment, found := cs[fmt.Sprintf("%v", paths)]; found && comment.leading != "" {
		return comment.leading
	}
	return name
}

// newPackage .
func newPackage(name string) *Package {
	return &Package{
		Name:       name,
		Version:    version(),
		Services:   make([]*Service, 0),
		Enums:      make([]*Enum, 0),
		EnumDic:    make(map[string]*Enum, 0),
		Messages:   make([]*Message, 0),
		MessageDic: make(map[string]*Message, 0),
	}
}

// newService .
func newService(name, desc string) *Service {
	return &Service{
		Name:        name,
		Description: desc,
		Methods:     make([]*ServiceMethod, 0),
	}
}

// newServiceMethod .
func newServiceMethod(name, desc string) *ServiceMethod {
	return &ServiceMethod{
		Name:        name,
		Method:      MethodPost,
		Description: desc,
	}
}

// newEnum .
func newEnum(name, desc string) *Enum {
	return &Enum{
		Name:        name,
		Description: desc,
		Fields:      make([]*EnumField, 0),
	}
}

// newMessage .
func newMessage(name, desc string) *Message {
	return &Message{
		Name:        name,
		Description: desc,
		Fields:      make([]*MessageField, 0),
	}
}
