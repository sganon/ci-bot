package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (api API) handleMR(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `{ "error": "unable to parse form" }`)
		return
	}

	username := r.Form.Get("text")
	if username == "" {

	}
	fmt.Println(username)
}
