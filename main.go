package main

import (
	"embed"
	"github.com/labstack/echo/v4"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

//go:embed dist
var staticPublic embed.FS

func main() {
	res, _ := fs.Sub(staticPublic, "dist")
	var handler http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		upath := request.URL.Path
		name := path.Clean(upath)
		name = strings.TrimPrefix(name, "/")
		ext := filepath.Ext(name)
		var f fs.File
		var err error
		if ext == "" {
			ext = "html"
			f, err = res.Open("index.html")
		} else {
			f, err = res.Open(name)
		}
		if err != nil {
			// todo error handle
		}

		writer.Header().Set(echo.HeaderContentType, mime.TypeByExtension(ext))
		writer.WriteHeader(200)
		io.Copy(writer, f) // todo error handle
	}

	e := echo.New()
	e.Any("/*", echo.WrapHandler(handler))
	e.Start(":1234")
}
