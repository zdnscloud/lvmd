package main

import (
	"flag"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zdnscloud/cement/log"
	pb "github.com/zdnscloud/lvmd/proto"
	"github.com/zdnscloud/lvmd/server"
)

func main() {
	var addr string
	flag.StringVar(&addr, "listen", ":1736", "server listen address")
	flag.Parse()

	log.InitLogger(log.Debug)
	defer log.CloseLogger()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed:%s", err.Error())
	}

	svr := server.NewServer()
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterLVMServer(grpcServer, &svr)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("run grpc server failed:%s", err.Error())
	}
}
