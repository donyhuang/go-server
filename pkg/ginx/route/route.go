package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Route struct {
	Path       string            `json:"-"`
	Method     string            `json:"-"`
	Parent     []string          `json:"-"`
	Handle     []gin.HandlerFunc `json:"-"`
	Desc       string            `json:"desc"`
	Permission string            `json:"permission"`
	AllPath    string            `json:"route"`
}
type Routes struct {
	r     []Route
	group map[string]gin.IRouter
}

func NewRoutes() *Routes {
	return &Routes{
		r:     make([]Route, 0),
		group: make(map[string]gin.IRouter),
	}
}
func (r *Routes) Add(routes ...Route) {
	for _, v := range routes {
		nowPath := strings.TrimLeft(v.Path, "/")
		v.AllPath = strings.Join(v.Parent, "") + "/" + nowPath
		v.Permission = strings.ReplaceAll(v.AllPath, "/", "_")
		r.r = append(r.r, v)
	}
}
func (r *Routes) GetAllRoutes() []Route {
	return r.r
}
func (r *Routes) BuildRoute(engine *gin.Engine) {
	for _, route := range r.r {

		var iRoute gin.IRouter = engine
		var path string
		if len(route.Parent) > 0 {
			for _, p := range route.Parent {
				path += p
				if _, ok := r.group[path]; !ok {
					r.group[path] = iRoute.Group(p)
				}
				iRoute = r.group[path]
			}
		}
		iRoute.Handle(route.Method, route.Path, route.Handle...)
	}
}
func (r *Routes) GetAuthHandle(lazy bool) gin.HandlerFunc {
	var cRoutesPoint *[]Route
	if !lazy {
		cRoutes := make([]Route, len(r.r), len(r.r))
		copy(cRoutes, r.r)
		cRoutesPoint = &cRoutes
	} else {
		cRoutesPoint = &r.r
	}
	return func(context *gin.Context) {
		rsp := make(map[string][]Route)
		for _, rr := range *cRoutesPoint {
			key := strings.Join(rr.Parent, "")
			if _, ok := rsp[key]; !ok {
				rsp[key] = make([]Route, 0)
			}
			rsp[key] = append(rsp[key], rr)
		}
		context.JSON(http.StatusOK, rsp)
	}
}
func (r *Routes) GetGroup(path string) gin.IRouter {
	return r.group[path]
}
