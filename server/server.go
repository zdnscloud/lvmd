/*

Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package server

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/zdnscloud/lvmd/commands"
	pb "github.com/zdnscloud/lvmd/proto"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) ListLV(ctx context.Context, in *pb.ListLVRequest) (*pb.ListLVReply, error) {
	lvs, err := commands.ListLV(ctx, in.VolumeGroup)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to list LV: %v\nCommandOutput: %v", err, lvs)
	}

	pblvs := make([]*pb.LogicalVolume, len(lvs))
	for i, v := range lvs {
		pblvs[i] = v.ToProto()
	}
	return &pb.ListLVReply{Volumes: pblvs}, nil
}

func (s Server) CreateLV(ctx context.Context, in *pb.CreateLVRequest) (*pb.CreateLVReply, error) {
	log, err := commands.CreateLV(ctx, in.VolumeGroup, in.Name, in.Size, in.Mirrors, in.Tags)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to create lv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.CreateLVReply{CommandOutput: log}, nil
}

func (s Server) CreateThinPool(ctx context.Context, in *pb.CreateThinPoolRequest) (*pb.CreateThinPoolReply, error) {
	log, err := commands.CreateThinPoolUseAllSize(ctx, in.VolumeGroup, in.Pool)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to create thin pool: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.CreateThinPoolReply{CommandOutput: log}, nil
}

func (s Server) ChangeLV(ctx context.Context, in *pb.ChangeLVRequest) (*pb.ChangeLVReply, error) {
	log, err := commands.ChangeLV(ctx, in.VolumeGroup, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to change lv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.ChangeLVReply{CommandOutput: log}, nil
}

func (s Server) CreateThinLV(ctx context.Context, in *pb.CreateThinLVRequest) (*pb.CreateThinLVReply, error) {
	vg := fmt.Sprintf("%s/%s", in.VolumeGroup, in.Pool)
	log, err := commands.CreateThinLV(ctx, vg, in.Name, in.Size, in.Mirrors, in.Tags)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to create thin lv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.CreateThinLVReply{CommandOutput: log}, nil
}

func (s Server) RemoveLV(ctx context.Context, in *pb.RemoveLVRequest) (*pb.RemoveLVReply, error) {
	log, err := commands.RemoveLV(ctx, in.VolumeGroup, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to remove lv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.RemoveLVReply{CommandOutput: log}, nil
}

func (s Server) CloneLV(ctx context.Context, in *pb.CloneLVRequest) (*pb.CloneLVReply, error) {
	log, err := commands.CloneLV(ctx, in.SourceName, in.DestName)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to clone lv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.CloneLVReply{CommandOutput: log}, nil
}

func (s Server) ResizeLV(ctx context.Context, in *pb.ResizeLVRequest) (*pb.ResizeLVReply, error) {
	log1, err := commands.ResizeLV(ctx, in.VolumeGroup, in.Name, in.Size)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to resize lv: %v\nCommandOutput: %v", err, streamline(log1))
	}
	log2, err := commands.ResizeLVe2fsck(ctx, in.VolumeGroup, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to e2fsck lv: %v\nCommandOutput: %v", err, streamline(log2))
	}
	log3, err := commands.ResizeLV2fs(ctx, in.VolumeGroup, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to resize2fs lv: %v\nCommandOutput: %v", err, streamline(log3))
	}
	return &pb.ResizeLVReply{CommandOutput: log1 + "|" + log2 + "|" + log3}, nil
}

func (s Server) ListVG(ctx context.Context, in *pb.ListVGRequest) (*pb.ListVGReply, error) {
	vgs, err := commands.ListVG(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to list vg: %v\nCommandOutput: %v", err, vgs)
	}

	pbvgs := make([]*pb.VolumeGroup, len(vgs))
	for i, v := range vgs {
		pbvgs[i] = v.ToProto()
	}
	return &pb.ListVGReply{VolumeGroups: pbvgs}, nil
}

func (s Server) CreateVG(ctx context.Context, in *pb.CreateVGRequest) (*pb.CreateVGReply, error) {
	log, err := commands.CreateVG(ctx, in.Name, in.PhysicalVolume, in.Tags)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to create vg: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.CreateVGReply{CommandOutput: log}, nil
}

func (s Server) ExtendVG(ctx context.Context, in *pb.ExtendVGRequest) (*pb.ExtendVGReply, error) {
	log, err := commands.ExtendVG(ctx, in.Name, in.PhysicalVolume)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to extend vg: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.ExtendVGReply{CommandOutput: log}, nil
}

func (s Server) ReduceVG(ctx context.Context, in *pb.ExtendVGRequest) (*pb.ExtendVGReply, error) {
	log, err := commands.ReduceVG(ctx, in.Name, in.PhysicalVolume)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to reduce vg: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.ExtendVGReply{CommandOutput: log}, nil
}

func (s Server) RemoveVG(ctx context.Context, in *pb.CreateVGRequest) (*pb.RemoveVGReply, error) {
	log, err := commands.RemoveVG(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to remove vg: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.RemoveVGReply{CommandOutput: log}, nil
}

func (s Server) AddTagLV(ctx context.Context, in *pb.AddTagLVRequest) (*pb.AddTagLVReply, error) {
	log, err := commands.AddTagLV(ctx, in.VolumeGroup, in.Name, in.Tags)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to add tags to lv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.AddTagLVReply{CommandOutput: log}, nil
}

func (s Server) RemoveTagLV(ctx context.Context, in *pb.RemoveTagLVRequest) (*pb.RemoveTagLVReply, error) {
	log, err := commands.RemoveTagLV(ctx, in.VolumeGroup, in.Name, in.Tags)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to remove tags from lv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.RemoveTagLVReply{CommandOutput: log}, nil
}

func (s Server) CreatePV(ctx context.Context, in *pb.CreatePVRequest) (*pb.CreatePVReply, error) {
	log, err := commands.CreatePV(ctx, in.Block)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to create pv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.CreatePVReply{CommandOutput: log}, nil
}

func (s Server) RemovePV(ctx context.Context, in *pb.RemovePVRequest) (*pb.RemovePVReply, error) {
	log, err := commands.RemovePV(ctx, in.Block)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to remove pv: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.RemovePVReply{CommandOutput: log}, nil
}

func (s Server) ListPV(ctx context.Context, in *pb.ListPVRequest) (*pb.ListPVReply, error) {
	pvs, err := commands.ListPV(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to list pv: %v\nCommandOutput: %v", err, pvs)
	}
	pbpvs := make([]*pb.PVInfo, len(pvs))
	for i, v := range pvs {
		pbpvs[i] = v.ToProto()
	}
	return &pb.ListPVReply{Pvinfos: pbpvs}, nil
}

func (s Server) Validate(ctx context.Context, in *pb.ValidateRequest) (*pb.ValidateReply, error) {
	v, err := commands.Validate(ctx, in.Block)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to validate block: %v\nCommandOutput: %v", err, v)
	}
	return &pb.ValidateReply{Validate: v}, nil
}

func (s Server) Destory(ctx context.Context, in *pb.DestoryRequest) (*pb.DestoryReply, error) {
	log, err := commands.Destory(ctx, in.Block)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to destory block: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.DestoryReply{CommandOutput: log}, nil
}

func (s Server) Match(ctx context.Context, in *pb.MatchRequest) (*pb.MatchReply, error) {
	log := commands.Match(ctx, in.Block)
	return &pb.MatchReply{CommandOutput: log}, nil
}

func (s Server) GetPVNum(ctx context.Context, in *pb.CreateVGRequest) (*pb.GetPVNumReply, error) {
	log, err := commands.GetPVNum(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to get vg's pv num: %v\nCommandOutput: %v", err, streamline(log))
	}
	return &pb.GetPVNumReply{CommandOutput: log}, nil
}

func streamline(out string) string {
	var res string
	for _, l := range strings.Split(out, "\n") {
		if len(l) == 0 {
			continue
		}
		if !strings.Contains(l, "/etc/lvm/cache/.cache") {
			res = res + l + "\n"
		}
	}
	return strings.TrimSpace(res)
}
