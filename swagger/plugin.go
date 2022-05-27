package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/charlesbases/protoc-gen-swagger/conf"
	"github.com/charlesbases/protoc-gen-swagger/logger"
	"github.com/charlesbases/protoc-gen-swagger/protoc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const SwaggerVersion = "2.0"

const DefaultAPIHost = "0.0.0.0"

var DefaultSchemes = []string{"http", "https"}

// apiHost .
func apiHost() string {
	if len(conf.Get().Host) != 0 {
		return conf.Get().Host
	}
	return DefaultAPIHost
}

// New .
func New(p *protoc.Package) *Swagger {
	var s = &Swagger{
		name: p.Name + ".json",
		p:    p,

		Swagger: SwaggerVersion,
		Info: &Info{
			Title:       p.Name,
			Version:     p.Version,
			Description: p.Name,
		},
		Host:     apiHost(),
		BasePath: "",
		Schemes:  DefaultSchemes,
		Paths:    make(map[string]map[string]*API, 0),
	}

	s.tidy()

	var swg = sync.WaitGroup{}
	swg.Add(2)

	go func() {
		s.parseServices()

		swg.Done()
	}()

	go func() {
		s.parseDefinitions()

		swg.Done()
	}()

	swg.Wait()
	return s
}

// Generater .
func (s *Swagger) Generater() *pluginpb.CodeGeneratorResponse_File {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		logger.Fatal(err)
	}

	var content = string(data)
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &s.name,
		Content: &content,
	}
}

// tidy .
func (s *Swagger) tidy() {
	s.defs = make(map[string]string, len(s.p.Messages)+len(s.p.Enums))

	for _, item := range s.p.Enums {
		s.defs[item.Name] = fmt.Sprintf("#/definitions/%s", item.Name)
	}
	for _, item := range s.p.Messages {
		s.defs[item.Name] = fmt.Sprintf("#/definitions/%s", item.Name)
	}
}

// reflex return #/definitions/...
func (s *Swagger) reflex(defname string) *Schema {
	if ref, found := s.defs[defname]; found {
		return &Schema{Reflex: ref}
	} else {
		return &Schema{}
	}
}

// parsePaths .
func (s *Swagger) parseServices() {
	for _, srv := range s.p.Services {
		var tag = &Tag{
			Name:        srv.Name,
			Description: srv.Description,
		}

		for _, m := range srv.Methods {
			api := &API{
				Tags:       []string{tag.Name},
				Summary:    m.Description,
				Consumes:   []string{m.Consume},
				Produces:   []string{m.Produce},
				Parameters: make([]*Parameter, 0),
				Responses:  make(map[string]*Parameter),
			}

			api.parseResponses(s, m)
			api.parseParameters(s, m)

			s.push(m.Path, m.Method, api)
		}

		s.Tags = append(s.Tags, tag)
	}
}

type Position string

const (
	PositionHeader Position = "header"
	PositionQuery  Position = "query"
	PositionBody   Position = "body"
)

// position .
func (api *API) position(m *protoc.ServiceMethod) Position {
	switch m.Method {
	case http.MethodGet:
		return PositionQuery
	default:
		return PositionBody
	}
}

// parseResponses .
func (api *API) parseResponses(s *Swagger, m *protoc.ServiceMethod) {
	api.Responses = map[string]*Parameter{
		"200": {
			Description: "successful",
			Schema:      s.reflex(m.ResponseName),
		},
	}
}

// parseParameters .
func (api *API) parseParameters(s *Swagger, m *protoc.ServiceMethod) {
	// Header
	{
		// Authorization
		if len(conf.Get().Header.Auth) != 0 {
			api.Parameters = append(api.Parameters, &Parameter{
				In:          PositionHeader,
				Name:        conf.Get().Header.Auth,
				Type:        "string",
				Required:    false,
				Description: "Authorization in Header",
			})
		}
	}

	switch api.position(m) {
	case PositionBody:
		api.Parameters = append(api.Parameters, &Parameter{
			In:          PositionBody,
			Name:        m.Name,
			Required:    true,
			Description: m.Description,
			Schema:      s.reflex(m.RequestName),
		})
	case PositionQuery:
		api.Parameters = append(api.Parameters, &Parameter{
			In:          PositionQuery,
			Name:        m.Name,
			Type:        "array",
			Required:    true,
			Description: m.Description,
			Schema:      s.reflex(m.RequestName),
		})
	}
}

// parseDefinitions .
func (s *Swagger) parseDefinitions() {
	s.Definitions = make(map[string]*Definition)

	s.parseProtoEnum()
	s.parseProtoMessage()
}

// parseProtoEnum .
func (s *Swagger) parseProtoEnum() {
	for _, enum := range s.p.Enums {
		var def = &Definition{Type: "string", Enum: make([]string, 0, len(enum.Fields))}

		// key list
		for _, field := range enum.Fields {
			def.Enum = append(def.Enum, field.Name)
		}

		// default
		if len(def.Enum) != 0 {
			def.Default = def.Enum[0]
		}

		// desc TODO enum desc + enum.field desc
		def.Description = enum.Description

		s.Definitions[enum.Name] = def
	}
}

// parseProtoMessage .
func (s *Swagger) parseProtoMessage() {
	for _, mess := range s.p.Messages {
		var def = &Definition{Type: "object", Description: mess.Description}
		fields := make(map[string]*Definition, 0)

		for _, mf := range mess.Fields {
			fields[mf.ProtoName] = s.parseProtoMessageField(mf)
		}

		def.Nesteds = fields
		s.Definitions[mess.Name] = def
	}
}

// parseProtoMessageField .
func (s *Swagger) parseProtoMessageField(mf *protoc.MessageField) *Definition {
	var field = new(Definition)
	if def, found := prototypes[mf.ProtoType]; found {
		field = def
	} else {
		switch mf.ProtoType {
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			field.Reflex = s.reflex(mf.ProtoTypeName).Reflex
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			field.Reflex = s.reflex(mf.ProtoTypeName).Reflex
		}
	}

	// proto laber
	switch mf.ProtoLaber {
	// repeated
	case descriptorpb.FieldDescriptorProto_LABEL_REPEATED:
		return &Definition{
			Type:  "array",
			Items: field,
		}
	default:
		return field
	}
}

// push api
func (s *Swagger) push(uri string, method string, api *API) {
	method = strings.ToLower(method)

	if apis, found := s.Paths[uri]; found {
		if _, found := apis[method]; found {
			logger.Fatalf("duplicate route. %s [%s]", uri, method)
		}

		apis[method] = api
	} else {
		var apis = make(map[string]*API, 0)
		apis[method] = api

		s.Paths[uri] = apis
	}
}
