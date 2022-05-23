package protoc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// mark .
func (cs comments) mark(name string, paths ...int) *Mark {
	var mark = &Mark{Name: name, Desc: name, Method: http.MethodPost}

	if comment, found := cs[fmt.Sprintf("%v", paths)]; found && comment.leading != "" {
		if err := json.Unmarshal([]byte(comment.leading), mark); err != nil {
			mark.Desc = comment.leading
		}
	}

	// Method
	mark.Method = strings.ToUpper(mark.Method)

	// ContentType
	if len(mark.Consume) == 0 && mark.Method != http.MethodGet {
		mark.Consume = ContentTypeJson
	}
	if len(mark.Produce) == 0 {
		mark.Produce = ContentTypeJson
	}

	return mark
}
