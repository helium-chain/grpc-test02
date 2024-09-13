package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	pb "example.com/learn-grpc-02/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	// 方法一
	// creds, err1 := credentials.NewServerTLSFromFile(
	//		"/root/workspace/learn-grpc/key/test.pem",
	//		"/root/workspace/learn-grpc/key/test.key",
	//	)
	//
	//	if err1 != nil {
	//		fmt.Printf("证书错误：%v", err1)
	//		return
	//	}

	// 方法二
	cert, err := tls.LoadX509KeyPair(
		"/root/workspace/learn-grpc-02/key/test.pem",
		"/root/workspace/learn-grpc-02/key/test.key")
	if err != nil {
		fmt.Printf("私钥错误：%v", err)
		return
	}
	creds := credentials.NewServerTLSFromCert(&cert)

	listen, _ := net.Listen("tcp", ":9090")
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterOrderManagementServer(grpcServer, &service{})

	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var _ pb.OrderManagementServer = (*service)(nil)

var orders = make(map[string]pb.Order, 8)

func init() {
	// 测试数据
	orders["1"] = pb.Order{Id: "1", Items: []string{"1", "2", "3", "4", "5", "7"}, Destination: "101"}
	orders["2"] = pb.Order{Id: "2", Items: []string{"6", "5", "4", "3", "2", "0"}, Destination: "102"}
}

type service struct {
	pb.UnimplementedOrderManagementServer
}

// SearchOrders 搜索订单
func (s *service) SearchOrders(query *wrapperspb.StringValue, stream pb.OrderManagement_SearchOrdersServer) error {
	for _, order := range orders {
		for _, str := range order.Items {
			if strings.Contains(str, query.GetValue()) {
				if err := stream.Send(&order); err != nil {
					return status.Error(codes.NotFound, "not fond")
				}
			}
		}
	}

	return status.New(codes.OK, "").Err()
}
