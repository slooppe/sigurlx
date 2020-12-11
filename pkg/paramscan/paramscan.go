package paramscan

import (
	"strings"

	"github.com/drsigned/gos"
)

// Params is a
type Params struct {
	Param string   `json:"param,omitempty"`
	Risks []string `json:"risks,omitempty"`
}

// Run is a
func Run(URL string, params []Params) ([]string, []Params, error) {
	paramsList := make([]string, 0)
	riskyParams := make([]Params, 0)

	parsedURL, err := gos.ParseURL(URL)
	if err != nil {
		return paramsList, riskyParams, err
	}

	for parameter := range parsedURL.Query() {
		if strings.HasSuffix(parameter, "\\") {
			parameter = strings.TrimSuffix(parameter, "\\")
		}

		paramsList = append(paramsList, parameter)

		for param := range params {
			if parameter == params[param].Param {
				riskyParams = append(riskyParams, params[param])
				break
			}
		}
	}

	return paramsList, riskyParams, nil
}
