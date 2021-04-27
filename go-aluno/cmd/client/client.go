package main

import (
	"fmt"
	"log"
	"context"
	"io"
	"time"

	"github.com/guiferrari/fullcycle-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Cound not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User {
		Id:		"0",
		Name: 	"Gordin",
		Email: 	"gordo@foda.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Cound not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User {
		Id:		"0",
		Name: 	"Gordin",
		Email: 	"gordo@foda.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Cound not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cound not receave the msg: %v", err)
		}
		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id: 	"g1",
			Name: 	"Guilherme 1",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g2",
			Name: 	"Guilherme 2",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g3",
			Name: 	"Guilherme 3",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g3",
			Name: 	"Guilherme 3",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g4",
			Name: 	"Guilherme 4",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g5",
			Name: 	"Guilherme 5",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g6",
			Name: 	"Guilherme 6",
			Email: 	"gui@gui",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id: 	"g1",
			Name: 	"Guilherme 1",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g2",
			Name: 	"Guilherme 2",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g3",
			Name: 	"Guilherme 3",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g3",
			Name: 	"Guilherme 3",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g4",
			Name: 	"Guilherme 4",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g5",
			Name: 	"Guilherme 5",
			Email: 	"gui@gui",
		},
		&pb.User{
			Id: 	"g6",
			Name: 	"Guilherme 6",
			Email: 	"gui@gui",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	go func() {

		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break;
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break;
			}
			fmt.Printf("Recebendo user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}

		close(wait)

	}()

	<-wait

}