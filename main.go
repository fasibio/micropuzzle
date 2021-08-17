package main

import (
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/fasibio/micropuzzle/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {
	logger.Initialize("info")
	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "gzip"))

	ChiFileServer(r, "/", http.Dir("./public"))

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})
	defer server.Close()
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})
	r.Handle("/stream", server)

	logger.Get().Infow("Start Server on Port :3000")
	logger.Get().Fatal(http.ListenAndServe(":3000", r))
}

func mimeTypeForFile(file string) string {
	// We use a built in table of the common types since the system
	// TypeByExtension might be unreliable. But if we don't know, we delegate
	// to the system.
	ext := filepath.Ext(file)
	switch ext {
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"

		// ...

	default:
		return mime.TypeByExtension(ext)
	}
}

func ChiFileServer(r chi.Router, path string, root http.FileSystem) {

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		f, err := root.Open(path)
		mimetype := mimeTypeForFile(path)
		w.Header().Set("Content-Type", mimetype)
		if err == nil {
			if mimetype == "application/javascript" {
				io.Copy(w, f)
				return
			}
			err := handleTemplate(f, w, r)
			if err != nil {
				logger.Get().Warnw("Error handle template", "error", err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			// logger.Get().Info("Will return Fallback ", path)
			// f, err = root.Open("/index.html")
			// handleTemplate(f, w, r)
			// if err != nil {
			// 	logger.Get().Error("Error by return fallback ", err)
			// }
		}
	})
}

func handleTemplate(f http.File, dst io.Writer, r *http.Request) error {

	var maxLoadingTime time.Duration = 45 * time.Millisecond

	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	cache := NewInMemoryHandler()
	handler, err := NewTemplateHandler(r, maxLoadingTime, &cache)
	if err != nil {
		return err
	}
	tmpl, err := template.New("httptemplate").Parse(string(text))
	if err != nil {
		return err
	}

	return tmpl.Execute(dst, handler)
}
