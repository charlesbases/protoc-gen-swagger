package protoc

// Mark .
type Mark struct {
	Name string `json:"-"`

	URI    string `json:"uri,omitempty"`
	Desc   string `json:"desc,omitempty"`
	Method string `json:"method,omitempty"`
	// Consume request ContentType. default: ContentTypeJson
	Consume string `json:"consume,omitempty"`
	// Produce response ContentType. default: ContentTypeJson
	Produce string `json:"produce,omitempty"`
}

// newPackage .
func (m *Mark) newPackage() *Package {
	return &Package{
		Name:       m.Name,
		Version:    version(),
		Services:   make([]*Service, 0),
		Enums:      make([]*Enum, 0),
		EnumDic:    make(map[string]*Enum, 0),
		Messages:   make([]*Message, 0),
		MessageDic: make(map[string]*Message, 0),
	}
}

// newService .
func (m *Mark) newService() *Service {
	return &Service{
		Name:        m.Name,
		Description: m.Desc,
		Methods:     make([]*ServiceMethod, 0),
	}
}

// newServiceMethod .
func (m *Mark) newServiceMethod() *ServiceMethod {
	return &ServiceMethod{
		Name:        m.Name,
		Path:        m.URI,
		Method:      m.Method,
		Consume:     m.Consume,
		Produce:     m.Produce,
		Description: m.Desc,
	}
}

// newMessage .
func (m *Mark) newMessage() *Message {
	return &Message{
		Name:        m.Name,
		Description: m.Desc,
		Fields:      make([]*MessageField, 0),
	}
}

// newEnum .
func (m *Mark) newEnum() *Enum {
	return &Enum{
		Name:        m.Name,
		Description: m.Desc,
		Fields:      make([]*EnumField, 0),
	}
}
