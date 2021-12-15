package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

func checkError(err error) error {
	log.Fatal(err)
	return err
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Home(w http.ResponseWriter, request *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		checkError(err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		checkError(err)
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		checkError(err)
	}

	return nil
}
