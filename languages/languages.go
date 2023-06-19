package languages

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"strings"

	"ghx/utils"

	"github.com/goccy/go-yaml"
)

const (
	LinguistVersion  string = "v7.25.0"
	LanguagesYAMLURL string = "https://api.github.com/repositories/1725199/contents/lib/linguist/languages.yml"
	HexColor         string = "#%02x%02x%02x"
)

var (
	errInvalidHexLength    = fmt.Errorf("invalid length, must be 7 or 4")
	structFieldNameMatcher = func(fieldName string) func(s string) bool {
		return func(s string) bool { return strings.EqualFold(fieldName, s) }
	}
)

type LanguageType string

const (
	Data        LanguageType = "data"
	Markup      LanguageType = "markup"
	Programming LanguageType = "programming"
	Prose       LanguageType = "prose"
	None        LanguageType = "none"
)

type Language struct {
	ID                 uint64       `yaml:"language_id"`
	Name               string       `yaml:"name"`
	Group              string       `yaml:"group,omitempty"`
	FSName             string       `yaml:"fs_name,omitempty"`
	Aliases            []string     `yaml:"aliases,omitempty"`
	Filenames          []string     `yaml:"filenames,omitempty"`
	Interpreters       []string     `yaml:"interpreters,omitempty"`
	Type               LanguageType `yaml:"type"`
	Color              color.RGBA64 `yaml:"color"`
	Extensions         []string     `yaml:"extensions"`
	TextmateScope      string       `yaml:"tm_scope,omitempty"`
	AceMode            string       `yaml:"ace_mode,omitempty"`
	CodeMirrorMode     string       `yaml:"codemirror_mode,omitempty"`
	CodeMirrorMIMEType string       `yaml:"codemirror_mime_type,omitempty"`
}

func HexStringToRGBA64(s string) (c color.RGBA64, err error) {
	c.A = 0xff

	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, HexColor, &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = errInvalidHexLength

	}

	return
}

func RGBA64ToHexString(c color.RGBA64) string {
	return fmt.Sprintf(HexColor, c.R, c.G, c.B)
}

func (l *Language) UnmarshalYAML(unmarshalFn func(interface{}) error) error {
	raw := map[string]interface{}{}
	if err := unmarshalFn(&raw); err != nil {
		return err
	}

	l.ID = raw["language_id"].(uint64)
	l.Type = LanguageType(raw["type"].(string))

	if colorRaw, ok := raw["color"]; ok {
		colorParsed, err := HexStringToRGBA64(colorRaw.(string))
		if err != nil {
			return err
		}

		l.Color = colorParsed
	}

	if extsRaw, ok := raw["extensions"]; ok {
		for _, iext := range extsRaw.([]interface{}) {
			l.Extensions = append(l.Extensions, iext.(string))
		}
	}

	l.TextmateScope = raw["tm_scope"].(string)
	l.AceMode = raw["ace_mode"].(string)

	return nil
}

type GetLanguagesOpts struct{ SortField string }

func GetLanguages(opts GetLanguagesOpts) ([]Language, error) {
	sortField := ""
	if opts.SortField != "" {
		sortField = opts.SortField
	}

	languagesYAMLURL, err := url.Parse(LanguagesYAMLURL)
	if err != nil {
		return nil, err
	}

	q := languagesYAMLURL.Query()
	q.Add("ref", LinguistVersion)
	languagesYAMLURL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    languagesYAMLURL,
	})
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resJSON := map[string]interface{}{}
	if err := json.Unmarshal(resBody, &resJSON); err != nil {
		return nil, err
	}

	yamlBytes := []byte{}
	if yamlBytes, err = base64.StdEncoding.DecodeString(
		resJSON["content"].(string),
	); err != nil {
		return nil, err
	}

	languagesMap := map[string]Language{}
	if err := yaml.Unmarshal(yamlBytes, &languagesMap); err != nil {
		return nil, err
	}

	languages := []Language{}
	for languageName, language := range languagesMap {
		language.Name = languageName
		languages = append(languages, language)
	}

	utils.SortSlice(&utils.SortSliceOpts[Language]{
		Slice: languages, SortField: sortField,
	})

	return languages, nil
}
