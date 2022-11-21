package commands

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/bregydoc/gtranslate"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

var ToVi = &cobra.Command{
	Use:   "tovi",
	Short: "convert en to vi",
	RunE:  toViF,
}

func init() {
	ToVi.Flags().String("translation-dir", "../server/i18n", "Path to folder with the sPhoton Chat server translate")
	ToVi.Flags().String("replace", "../server/i18n/replace.json", "Path to file of list string replace")
	I18nCmd.AddCommand(ToVi)
}

type convertToViFromEn struct {
	replace map[string]string
	ids     []string
	srcEn   map[string]interface{}
	srcVi   map[string]interface{}
}

func toViF(command *cobra.Command, args []string) error {
	translationDir, err := command.Flags().GetString("translation-dir")
	if err != nil {
		return errors.New("invalid translation-dir parameter")
	}
	replaceFile, err := command.Flags().GetString("replace")
	if err != nil {
		return errors.New("invalid replace parameter")
	}
	replaceF, err := os.ReadFile(replaceFile)
	if err != nil {
		return err
	}
	converter := &convertToViFromEn{}
	err = json.Unmarshal(replaceF, &converter.replace)
	if err != nil {
		return err
	}

	enFile, err := os.ReadFile(path.Join(translationDir, "en.json"))
	if err != nil {
		return err
	}
	en := []Translation{}
	err = json.Unmarshal(enFile, &en)
	if err != nil {
		return err
	}
	src := make(map[string]interface{})
	ids := []string{}
	for _, v := range en {
		src[v.Id] = v.Translation
		ids = append(ids, v.Id)
	}
	converter.ids = ids
	converter.srcEn = src

	viFile, err := os.ReadFile(path.Join(translationDir, "vi.json"))
	if err != nil {
		return err
	}
	vi := []Translation{}
	err = json.Unmarshal(viFile, &vi)
	if err != nil {
		return err
	}
	srcVi := make(map[string]interface{})
	for _, v := range vi {
		srcVi[v.Id] = v.Translation
	}
	converter.srcVi = srcVi
	return converter.Write(path.Join(translationDir, "vi.json"))
}

func (c *convertToViFromEn) Write(file string) error {
	result := make([]Translation, 0)
	for _, id := range c.ids {
		value, ok := c.srcVi[id]
		if !ok {
			value, ok = c.srcEn[id]
			if !ok {
				panic("missing key: " + id)
			}
		}
		v, err := c.replaceString(value)
		if err != nil {
			return err
		}
		v, err = c.replaceMap(v)
		if err != nil {
			return err
		}
		result = append(result, Translation{
			Id:          id,
			Translation: v,
		})
	}
	content, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, content, 0666)
}

func (c *convertToViFromEn) replaceString(value interface{}) (interface{}, error) {
	if v, ok := value.(string); ok {
		if len(strings.TrimSpace(v)) == 0 {
			return v, nil
		}
		text, err := gtranslate.Translate(v, language.English, language.Vietnamese)
		if err != nil {
			return nil, err
		}
		for search, replace := range c.replace {
			text = strings.ReplaceAll(text, search, replace)
		}
		return text, nil
	}
	return value, nil
}

func (c *convertToViFromEn) replaceMap(input interface{}) (interface{}, error) {
	if rows, ok := input.([]map[string]string); ok {
		result := make([]map[string]string, 0)
		for _, row := range rows {
			r := row
			for id, value := range row {
				text, err := gtranslate.Translate(value, language.English, language.Vietnamese)
				if err != nil {
					return nil, err
				}
				for search, replace := range c.replace {
					text = strings.ReplaceAll(text, search, replace)
				}
				r[id] = text
			}
			result = append(result, r)
		}
		return result, nil
	}
	return input, nil
}
