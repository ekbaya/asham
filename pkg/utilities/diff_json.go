package utilities

import (
	"encoding/json"

	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

func DiffJSON(oldJSON, newJSON string) (string, error) {
	var oldDoc map[string]any
	var newDoc map[string]any

	err := json.Unmarshal([]byte(oldJSON), &oldDoc)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(newJSON), &newDoc)
	if err != nil {
		return "", err
	}

	differ := gojsondiff.New()
	oldJSONBytes, err := json.Marshal(oldDoc)
	if err != nil {
		return "", err
	}

	newJSONBytes, err := json.Marshal(newDoc)
	if err != nil {
		return "", err
	}

	diff, err := differ.Compare(oldJSONBytes, newJSONBytes)
	if err != nil {
		return "", err
	}

	if !diff.Modified() {
		return "No changes", nil
	}

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       true,
	}
	asciiFormatter := formatter.NewAsciiFormatter(oldDoc, config)
	result, err := asciiFormatter.Format(diff)
	if err != nil {
		return "", err
	}

	return result, nil
}
