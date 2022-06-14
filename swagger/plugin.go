package swagger

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/charlesbases/protoc-gen-swagger/conf"
	"github.com/charlesbases/protoc-gen-swagger/logger"
	"github.com/charlesbases/protoc-gen-swagger/protoc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const SwaggerVersion = "2.0"

const DefaultAPIHost = "127.0.0.1"

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
	var service = conf.Get().Service
	if len(service) == 0 {
		service = p.Name
	}

	var s = &Swagger{
		name: p.Name + ".json",
		p:    p,

		Swagger: SwaggerVersion,
		Info: &Info{
			Title:       service,
			Version:     p.Version,
			Description: service,
		},
		Host:     apiHost(),
		BasePath: "",
		Schemes:  DefaultSchemes,
		Paths:    make(map[string]map[string]*API, 0),
	}

	s.parseDefinitions()
	s.parseServices()

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

// reflex return #/definitions/...
func (s *Swagger) reflex(defname string) *Definition {
	return &Definition{Reflex: refprefix + defname}
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
			api.parseParameter(s, m)

			s.push(m.Path, m.Method, api)
		}

		s.Tags = append(s.Tags, tag)
	}
}

const refprefix = "#/definitions/"

// parseDefinitions .
func (s *Swagger) parseDefinitions() {
	s.Definitions = make(map[string]*Definition, len(s.p.Messages)+len(s.p.Enums))

	// parse enums
	s.parseProtoEnum()

	// parse messages
	for _, mess := range s.p.Messages {
		s.parseProtoMessage(mess)
	}
}

// parseProtoEnum .
func (s *Swagger) parseProtoEnum() {
	for _, enum := range s.p.Enums {
		var def = &Definition{
			Name: enum.Name,
			Type: "string",
			Enum: make([]string, 0, len(enum.Fields)),
		}

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
func (s *Swagger) parseProtoMessage(mess *protoc.Message) {
	var def = &Definition{
		Name:        mess.Name,
		Type:        "object",
		Description: mess.Description,
	}
	fields := make(map[string]*Definition, 0)

	for _, mf := range mess.Fields {
		fields[mf.ProtoName] = s.parseProtoMessageField(mf)
	}

	def.Nesteds = fields

	s.Definitions[mess.Name] = def
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
			// 优先解析嵌套 message
			if _, found := s.Definitions[mf.ProtoTypeName]; !found {
				if mess, found := s.p.MessageDic[mf.ProtoTypeName]; found {
					s.parseProtoMessage(mess)
				}
			}

			if protoc.IsEntry(mf) {
				if entry, found := s.Definitions[mf.ProtoTypeName]; found && len(entry.Nesteds) != 0 {
					// if key, k_found := entry.Nesteds["key"]; k_found {
					//
					// }
					if val, v_found := entry.Nesteds["value"]; v_found {
						field.Entry = val
					}
				}
			} else {
				field.Reflex = s.reflex(mf.ProtoTypeName).Reflex
			}
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

type Position string

const (
	PositionFormData Position = "formData"
	PositionHeader   Position = "header"
	PositionQuery    Position = "query"
	PositionBody     Position = "body"
	PositionPath     Position = "path"
)

// parameterPosition .
func (api *API) parameterPosition(m *protoc.ServiceMethod) Position {
	if m.Consume == "multipart/form-data" {
		return PositionFormData
	}

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

// parseParameter .
func (api *API) parseParameter(s *Swagger, m *protoc.ServiceMethod) {
	api.parseParameterInHeader()
	api.parseParameterInPath(m)

	switch api.parameterPosition(m) {
	case PositionBody:
		api.parseParameterInBody(s, m)
	case PositionQuery:
		api.parseParameterInQuery(s, m)
	case PositionFormData:
		api.parseParamterInFormData(s, m)
	}
}

// parseParameterInHeader .
func (api *API) parseParameterInHeader() {
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

// parseParameter .
func (api *API) parseParameterInPath(m *protoc.ServiceMethod) {
	var uri = m.Path
	for len(uri) > 2 {
		l, r := strings.Index(uri, "{"), strings.Index(uri, "}")
		if l > 0 && r > 0 && r > l {
			api.Parameters = append(api.Parameters, &Parameter{
				In:       PositionPath,
				Name:     uri[l+1 : r],
				Type:     "string",
				Required: false,
			})
			uri = uri[r+1:]
		} else {
			return
		}
	}
}

// parseParameterInBody .
func (api *API) parseParameterInBody(s *Swagger, m *protoc.ServiceMethod) {
	api.Parameters = append(api.Parameters, &Parameter{
		In:          PositionBody,
		Name:        m.Name,
		Required:    false,
		Description: m.Description,
		Schema:      s.reflex(m.RequestName),
	})
}

// parseParameter .
func (api *API) parseParameterInQuery(s *Swagger, m *protoc.ServiceMethod) {
	if mess, found := s.Definitions[m.RequestName]; found {
		// message fields
		for name, field := range mess.Nesteds {
			switch field.Type {
			case "array":
				// repeated nesteds
				if len(field.Items.Reflex) != 0 {
					// query 中的 nesteds 只允许为 enum
					if def, found := s.Definitions[strings.TrimPrefix(field.Items.Reflex, refprefix)]; found && len(def.Enum) != 0 {
						api.Parameters = append(api.Parameters, &Parameter{
							In:          PositionQuery,
							Name:        name,
							Type:        field.Type,
							Required:    false,
							Description: field.Description,
							Items: &Definition{
								Type:    def.Type,
								Enum:    def.Enum,
								Default: def.Default,
							},
						})
					}
				} else {
					api.Parameters = append(api.Parameters, &Parameter{
						In:          PositionQuery,
						Name:        name,
						Type:        field.Type,
						Required:    false,
						Description: field.Description,
						Items: &Definition{
							Type: field.Items.Type,
						},
					})
				}
			default:
				// nesteds
				if len(field.Reflex) != 0 {
					// query 中的 nesteds 只允许为 enum
					if def, found := s.Definitions[strings.TrimPrefix(field.Reflex, refprefix)]; found && len(def.Enum) != 0 {
						api.Parameters = append(api.Parameters, &Parameter{
							In:          PositionQuery,
							Name:        name,
							Type:        def.Type,
							Required:    false,
							Enum:        def.Enum,
							Default:     def.Default,
							Description: def.Description,
						})
					}
				} else {
					api.Parameters = append(api.Parameters, &Parameter{
						In:          PositionQuery,
						Name:        name,
						Type:        field.Type,
						Required:    false,
						Description: field.Description,
					})
				}
			}
		}
	}
}

// parseParamterInFormData .
func (api *API) parseParamterInFormData(s *Swagger, m *protoc.ServiceMethod) {
	if mess, found := s.Definitions[m.RequestName]; found {
		// message fields
		for name, field := range mess.Nesteds {
			switch field.Type {
			case "array":
				// multipart/form-data 参数不支持 array
			default:
				// nesteds
				if len(field.Reflex) != 0 {
					// multipart/form-data 中的 nesteds 只允许为 enum
					if def, found := s.Definitions[strings.TrimPrefix(field.Reflex, refprefix)]; found && len(def.Enum) != 0 {
						api.Parameters = append(api.Parameters, &Parameter{
							In:          PositionFormData,
							Name:        name,
							Type:        def.Type,
							Required:    false,
							Enum:        def.Enum,
							Default:     def.Default,
							Description: def.Description,
						})
					}
				} else {
					if field.Format == "bytes" {
						api.Parameters = append(api.Parameters, &Parameter{
							In:          PositionFormData,
							Name:        name,
							Type:        "file",
							Required:    false,
							Description: field.Description,
						})
					} else {
						api.Parameters = append(api.Parameters, &Parameter{
							In:          PositionFormData,
							Name:        name,
							Type:        field.Type,
							Required:    false,
							Description: field.Description,
						})
					}
				}
			}
		}
	}
}
