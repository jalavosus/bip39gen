package outformat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/naoina/toml"
)

//go:generate stringer -type OutFormat -linecomment

type (
	OutFormat    uint
	FormatTextFn func(any, *bytes.Buffer)
)

const (
	JSON OutFormat = iota // json
	YAML                  // yaml
	TOML                  // toml
	Text                  // text
)

func FromString(s string) (f OutFormat) {
	switch strings.ToLower(s) {
	case JSON.String():
		f = JSON
	case YAML.String():
		f = YAML
	case TOML.String():
		f = TOML
	case Text.String():
		f = Text
	default:
		f = YAML
	}

	return
}

func (f OutFormat) Marshal(data any, rest ...any) (marshaled []byte) {
	var marshalErr error

	switch f {
	case JSON:
		marshaled, marshalErr = json.MarshalIndent(data, "", "  ")
	case YAML:
		marshaled, marshalErr = yaml.Marshal(data)
	case TOML:
		marshaled, marshalErr = toml.Marshal(data)
	case Text:
		format := rest[0].(func(any, *bytes.Buffer))
		marshaled = FormatText(data, format)
	}

	if marshalErr != nil {
		fmt.Println(marshalErr)
	}

	return
}

func FormatText(data any, format FormatTextFn) []byte {
	buf := new(bytes.Buffer)

	format(data, buf)

	return buf.Bytes()
}
