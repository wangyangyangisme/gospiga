//go:generate protoc -I ../../proto --go_out=plugins=grpc:../../proto ../../proto/finder.proto

package grpc

import (
	"context"
	"fmt"

	pb "gospiga/proto"
)

type finderServer struct {
	app App
}

func NewFinderServer(app App) *finderServer {
	return &finderServer{app}
}

func (s *finderServer) AllRecipeTags(ctx context.Context, req *pb.AllRecipeTagsRequest) (*pb.AllRecipeTagsResponse, error) {
	tags, err := s.app.AllRecipeTags()
	if err != nil {
		return nil, fmt.Errorf("error retrieving recipe tags: %w", err)
	}

	return &pb.AllRecipeTagsResponse{Tags: tags}, nil
}
