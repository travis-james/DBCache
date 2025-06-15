package server

import pb "github.com/travis-james/DBCache/internal/genproto"

type Server struct {
	pb.UnimplementedDBCacheServiceServer
}
