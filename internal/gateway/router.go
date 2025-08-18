package gateway

import (
	"context"
	"net/http"

	grpcInternal "github.com/travis-james/DBCache/internal/client"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
	"google.golang.org/protobuf/encoding/protojson"
)

type router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter(grpcClient *grpcInternal.Client) *router {
	retval := &router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
	retval.addRoute("GET", "/health", checkHealth(grpcClient))
	return retval
}

func (r *router) addRoute(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}
	r.routes[path][method] = handler
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handlers, ok := r.routes[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
		return
	}
	handler, methodExists := handlers[req.Method]
	if !methodExists {
		http.NotFound(w, req)
		return
	}
	handler(w, req)
}

func checkHealth(grpcClient *grpcInternal.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := grpcClient.CheckHealth(context.Background(), &pb.Empty{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := protojson.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
