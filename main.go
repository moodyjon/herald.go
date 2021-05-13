package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "github.com/lbryio/hub/protobuf/go"
	"github.com/lbryio/hub/server"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "serve" {
		l, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterHubServer(s, &server.Server{})

		log.Printf("listening on %s\n", l.Addr().String())
		if err := s.Serve(l); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		return
	}

	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewHubClient(conn)

	var query string
	if len(os.Args) > 1 {
		query = strings.Join(os.Args[1:], " ")
	} else {
		log.Printf("error: no search query provided\n")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//searchRequest := &pb.SearchRequest{Query: query}
	//searchRequest := &pb.SearchRequest{XId: [][]byte{[]byte(query)}}
	searchRequest := &pb.SearchRequest{ClaimType: []string{query}}
	r, err := c.Search(ctx, searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found %d results\n", r.GetTotal())

	for _, t := range r.Txos {
		fmt.Printf("%s:%d\n", server.FromHash(t.TxHash), t.Nout)
	}
}
