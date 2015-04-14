package api

import "github.com/unrolled/render"

var (
	renderer      = render.New(render.Options{IndentJSON: true})
	emptyResponse = make(map[string]string)
)
