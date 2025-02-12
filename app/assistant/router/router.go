package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/west2-online/DomTok/app/assistant/handler"
)

// GeneratedRegister registers routers.
func GeneratedRegister(r *server.Hertz) {
	root := r.Group("/")
	root.Any("/", handler.Entrypoint)
}
