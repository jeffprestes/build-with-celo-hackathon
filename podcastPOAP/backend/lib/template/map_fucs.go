package template

import (
	"html/template"

	"github.com/jeffprestes/build-with-celo-hackathon/podcastPOAP/backend/lib/contx"
)

// FuncMaps to view
func FuncMaps() []template.FuncMap {
	return []template.FuncMap{
		map[string]interface{}{
			"Tr": contx.I18n,
		}}
}
