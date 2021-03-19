package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"path"
)

//go:embed static
var statics embed.FS

func main() {
	router := gin.Default()

	router.Use(ServeEmbedDir("static", statics))
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404")
	})

	router.Run(":8080")
}

type staticFS struct {
	fs  embed.FS
	dir string
}

func (s staticFS) Open(name string) (fs.File, error) {
	return s.fs.Open(path.Join(s.dir, name))
}
func (s staticFS) Exist(name string) bool {
	_, err := s.Open(name)
	return err == nil
}

func newStaticFS(basedir string, fs embed.FS) staticFS {
	return staticFS{
		fs:  fs,
		dir: basedir,
	}
}

func ServeEmbedDir(basedir string, efs embed.FS) gin.HandlerFunc {
	sfs := newStaticFS(basedir, efs)
	fileServer := http.FileServer(http.FS(sfs))
	return func(c *gin.Context) {
		if sfs.Exist(c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
