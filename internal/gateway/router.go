package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	grpcInternal "github.com/travis-james/DBCache/internal/client"
	pb "github.com/travis-james/DBCache/pkg/protobuf"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewRouter(grpcClient *grpcInternal.Client) *gin.Engine {
	router := gin.Default()
	router.GET("/health", checkHealth(grpcClient))
	router.GET("/data/:id", getData(grpcClient))

	router.Run("localhost:8080")
	return router
}

func checkHealth(grpcClient *grpcInternal.Client) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		resp, err := grpcClient.CheckHealth(context.Background(), &pb.Empty{})
		if err != nil {
			gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data, err := protojson.Marshal(resp)
		if err != nil {
			gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		gctx.Data(http.StatusOK, "application/json", data)
	}
}

func getData(grpcClient *grpcInternal.Client) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		resp, err := grpcClient.GetData(context.Background(), &pb.GetRequest{QueryId: gctx.Param("id")})
		if err != nil {
			gctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var respData map[string]any
		if err := json.Unmarshal([]byte(resp.Data), &respData); err != nil {
			gctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
			return
		}

		gctx.JSON(http.StatusOK, gin.H{
			"fromCache": resp.FromCache,
			"data":      respData,
		})
	}
}
