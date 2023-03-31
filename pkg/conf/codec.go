package conf

import (
	"encoding/json"
	"strings"
	"sync"
)

var once sync.Once

type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string

	GetEmptyPointer() interface{}
}

var registeredCodecs = make(map[string]Codec)

func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}

// GetCodec gets a registered Codec by content-subtype, or nil if no Codec is
// registered for the content-subtype.
//
// The content-subtype is expected to be lowercase.
func GetCodec(contentSubtype string) Codec {
	contentSubtype = strings.ToLower(contentSubtype)
	return registeredCodecs[contentSubtype]
}

type MapCodes struct{}

func (m MapCodes) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m MapCodes) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (m MapCodes) Name() string {
	return "json"
}

func (m MapCodes) GetEmptyPointer() interface{} {
	v := make(map[string]interface{})
	return &v
}

func init() {
	RegisterCodec(MapCodes{})
}
