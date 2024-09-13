package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "example.com/learn-grpc-02/ecommerce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	creds, _ := credentials.NewClientTLSFromFile(
		"/root/workspace/learn-grpc-02/key/test.pem",
		"*.heliu.site",
	)

	var opts []grpc.DialOption
	// 不带TLS这里是grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// 连接server端，使用ssl加密通信
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	fmt.Printf("now-Time: %s\n", time.Now().Format(time.DateTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.SearchOrders(ctx, &wrapperspb.StringValue{Value: "5"})

	if err != nil {
		log.Fatalf("error when calling SearchOrders: %v", err)
	}

	for {
		order, err := stream.Recv()
		if err == io.EOF {
			break
		}

		log.Println("SearchOrders:", order)
	}

	// Output:
	// 2024/09/13 12:24:17 SearchOrders: id:"1" items:"1" items:"2" items:"3" items:"4" items:"5" items:"7" destination:"101"
	// 2024/09/13 12:24:17 SearchOrders: id:"2" items:"6" items:"5" items:"4" items:"3" items:"2" items:"0" destination:"102"
}
