package languages

import (
	"image/color"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestHexStringToRGBA64(t *testing.T) {
	for _, tc := range []struct {
		hexString string
		rgba64    color.RGBA64
		err       error
	}{
		{"#fffbec", color.RGBA64{255, 251, 236, 0xff}, nil},
		{"#5e5", color.RGBA64{85, 238, 85, 0xff}, nil},
		{"", color.RGBA64{0x0, 0x0, 0x0, 0xff}, errInvalidHexLength},
	} {
		parsed, err := HexStringToRGBA64(tc.hexString)
		if err != nil {
			if tc.err != nil {
				assert.Equal(t, tc.err, err)
			} else {
				t.Error(err)
			}
		}

		assert.Equal(t, tc.rgba64, parsed)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	for _, tc := range []struct {
		raw         []byte
		unmarshaled Language
	}{
		{
			[]byte(`---
type: programming
color: "#814CCC"
extensions:
- ".bsl"
- ".os"
tm_scope: source.bsl
ace_mode: text
language_id: 0`),
			Language{
				ID:            0,
				Name:          "1C Enterprise",
				Type:          Programming,
				Color:         color.RGBA64{129, 76, 204, 255},
				Extensions:    []string{".bsl", ".os"},
				AceMode:       "text",
				TextmateScope: "source.bsl",
			},
		},
	} {
		lang := Language{}
		if err := yaml.Unmarshal([]byte(tc.raw), &lang); err != nil {
			t.Error(err)
		}
		lang.Name = tc.unmarshaled.Name

		assert.Equal(t, tc.unmarshaled, lang)
	}
}
