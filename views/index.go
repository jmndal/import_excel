package views

import (
	"net/http"

	"github.com/uadmin/uadmin"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	c := map[string]interface{}{}

	uadmin.RenderHTML(w, r, "templates/index.html", c)
}
