package main

import (
	"context"
	"fmt"
	"github.com/zdnscloud/lvmd/client"
	pb "github.com/zdnscloud/lvmd/proto"
	"time"
)

func main() {
	addr := "202.173.9.5:1736"
	timeout := 3 * time.Second
	cli, err := client.New(addr, timeout)
	defer cli.Close()
	if err != nil {
		fmt.Println("connect failed!", err)
	}
	/*
		req := pb.CreatePVRequest{
			Block: "/dev/vdb",
		}
		out, err := cli.CreatePV(context.TODO(), &req)
		fmt.Println(out.CommandOutput, err)
	*/
	/*
		req := pb.ExtendVGRequest{
			Name:           "k8s",
			PhysicalVolume: "/dev/vdc",
		}
		out, err := cli.ExtendVG(context.TODO(), &req)
		fmt.Println(out.CommandOutput, err)
	*/
	/*
		req := pb.ListPVRequest{}
		out, err := cli.ListPV(context.TODO(), &req)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range out.Pvinfos {
			fmt.Println(v)
		}
	*/
	/*
		req := pb.RemovePVRequest{
			Block: "/dev/vdb",
		}
		out, err := cli.RemovePV(context.TODO(), &req)
		fmt.Println(out.CommandOutput, err)
	*/
	req := pb.ValidateRequest{
		Block: "/dev/vdb",
	}
	out, err := cli.Validate(context.TODO(), &req)
	fmt.Println(out.Validate, err)
	/*
		req := pb.DestoryRequest{
			Block: "/dev/vdb",
		}
		out, err := cli.Destory(context.TODO(), &req)
		fmt.Println(out.CommandOutput, err)
	*/
}
