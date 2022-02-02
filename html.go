package reco

import (
	"bytes"
	"context"
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

//go:embed error.html
var errorHtml string

func HTML(ctx context.Context, rec *Recovery) ([]byte, error) {
	tmpl, err := template.New("error.html").Parse(errorHtml)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, rec); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func RenderHTML(ctx context.Context, w http.ResponseWriter, rec *Recovery) error {
	data, err := HTML(ctx, rec)
	if err != nil {
		return err
	}
	hdr := w.Header()
	hdr.Set("Content-Type", "text/html; charset=utf-8")
	hdr.Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(data)
	return err
}

func renderHTML(ctx context.Context, w http.ResponseWriter, rec *Recovery) {
	if err := RenderHTML(ctx, w, rec); err != nil {
		log.Printf("reco: error rendering HTML: %v", err)
	}
}

func HTTPRenderer(w http.ResponseWriter) Renderer {
	return func(ctx context.Context, rec *Recovery) error {
		return RenderHTML(ctx, w, rec)
	}
}
