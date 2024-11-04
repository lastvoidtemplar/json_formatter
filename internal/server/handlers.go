package server

import (
	"net/http"

	"github.com/lastvoidtemplar/json_formatter/web/template"
)

func PageHandler(w http.ResponseWriter, r *http.Request) {
	template.Page().Render(r.Context(), w)
}
