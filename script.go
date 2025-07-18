package jsrun

import (
	"github.com/nitsugaro/go-nstore"
	"github.com/nitsugaro/go-utils/encoding"
)

type Script struct {
	*nstore.Metadata

	Name       string `json:"name"`
	Type       string `json:"type"`
	CodeBase64 string `json:"code_base64"`
	rawCode    string `json:"-"`
}

func (s *Script) GetRawCode() (string, error) {
	if s.rawCode != "" {
		return s.rawCode, nil
	}

	bytes, err := encoding.DecodeBase64(s.CodeBase64)
	if err != nil {
		return "", err
	}

	s.rawCode = string(bytes)

	return s.rawCode, nil
}
