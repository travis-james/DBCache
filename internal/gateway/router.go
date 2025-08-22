package gateway

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	grpcInternal "github.com/travis-james/DBCache/internal/client"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewRouter(grpcClient *grpcInternal.Client) *gin.Engine {
	router := gin.Default()
	router.GET("/health", checkHealth(grpcClient))

	router.Run("localhost:8080")
	return router
}

func checkHealth(grpcClient *grpcInternal.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := grpcClient.CheckHealth(context.Background(), &pb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data, err := protojson.Marshal(resp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	}
}
