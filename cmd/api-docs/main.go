package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/realworld-aws-api/internal/app"
	"github.com/wolfeidau/realworld-aws-api/internal/logger"
)

var cfg struct {
	Version  kong.VersionFlag
	SpecFile string `help:"Path to the openapi spec yaml file." kong:"required"`
}

//go:embed public/*.html
var embededFiles embed.FS

func main() {
	kong.Parse(&cfg,
		kong.Vars{"version": fmt.Sprintf("%s_%s", app.Commit, app.BuildDate)}, // bind a var for version
	)

	log.Logger = logger.NewLogger()

	e := echo.New()

	e.HideBanner = true
	e.Logger.SetOutput(io.Discard)

	t := &Template{
		templates: template.Must(template.ParseFS(embededFiles, "public/*.html")),
	}

	e.Renderer = t
	e.GET("/", Index)
	e.GET("/openapi.yaml", FromFile(cfg.SpecFile))

	log.Info().Msg("listening on http://localhost:4040")
	e.Start(":4040")
}

func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]string{
		"OpenAPISpecURL": "/openapi.yaml",
	})
}

func FromFile(f string) echo.HandlerFunc {

	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal().Err(err).Msg("load of spec failed")
	}

	return func(c echo.Context) error {
		b := bytes.NewBuffer(data)

		return c.Stream(200, "application/x-yaml", b)
	}

}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Error().Err(err).Msg("render failed")
		return err
	}

	return nil
}
