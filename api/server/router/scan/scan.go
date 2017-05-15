package scan

import "github.com/maliceio/malice/api/server/router"

// scanRouter is a router to talk with the scans controller
type scanRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new scan router
func NewRouter(b Backend) router.Router {
	r := &scanRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routes to the scans controller
func (r *scanRouter) Routes() []router.Route {
	return r.routes
}

func (r *scanRouter) initRoutes() {
	r.routes = []router.Route{
		// GET
		router.NewGetRoute("/scans", r.getscansList),
		router.NewGetRoute("/scans/{name:.*}", r.getscanByName),
		// POST
		router.NewPostRoute("/scans/create", r.postscansCreate),
		router.NewPostRoute("/scans/prune", r.postscansPrune, router.WithCancel),
		// DELETE
		router.NewDeleteRoute("/scans/{name:.*}", r.deletescans),
	}
}
