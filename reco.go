package reco

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"time"
)

type Renderer func(ctx context.Context, rec *Recovery) error

type Recovery struct {
	File          string
	Line          int
	Title         string
	Err           error
	Stack         string
	Request       *http.Request
	DumpedRequest string
	Started       time.Time
	Elapsed       time.Duration
}

func HTTPHandler(handler http.Handler, renderers ...Renderer) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer HTTPReco(w, r, renderers...)
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func Reco(ctx context.Context, renderers ...Renderer) {
	if rec := recover(); rec != nil {
		reco(ctx, nil, nil, rec, renderers...)
	}
}

func HTTPReco(w http.ResponseWriter, r *http.Request, renderers ...Renderer) {
	if rec := recover(); rec != nil {
		reco(r.Context(), w, r, rec, renderers...)
	}
}

func reco(ctx context.Context, w http.ResponseWriter, r *http.Request, recovered interface{}, renderers ...Renderer) {
	err, ok := recovered.(error)
	if !ok || err == nil {
		err = fmt.Errorf("%v", recovered)
	}
	var file string
	var line int
	var stack string
	if pa := uppermostPanic(); pa != nil {
		file = pa.file
		line = pa.line
		stack = formatStack(pa.pcSkip, false)
	} else {
		stack = string(debug.Stack())
	}
	var dumpedRequest string
	var started time.Time
	var elapsed time.Duration
	if r != nil {
		if requestData, err := httputil.DumpRequest(r, false); err == nil {
			dumpedRequest = string(requestData)
		}

		if st := Started(ctx); !st.IsZero() {
			started = st
			elapsed = time.Now().Sub(started)
		}
	}
	var title string
	if r != nil {
		title = fmt.Sprintf("Panic on %s", r.Host)
	} else {
		title = "Panic"
	}
	rec := &Recovery{
		File:          file,
		Line:          line,
		Title:         title,
		Err:           err,
		Stack:         stack,
		Request:       r,
		DumpedRequest: dumpedRequest,
		Started:       started,
		Elapsed:       elapsed,
	}
	if w != nil {
		debug := true
		if r != nil {
			if debugValue, ok := r.Context().Value(debugContextKey).(bool); ok {
				debug = debugValue
			}
		}
		if debug {
			renderHTML(ctx, w, rec)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Internal Server Error")
		}
	}
	for _, renderer := range renderers {
		if err := renderer(ctx, rec); err != nil {
			log.Printf("error logging error with %v: %v", renderer, err)
		}
	}
}
