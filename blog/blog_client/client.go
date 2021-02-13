package main

import (
	"Github/gRPC-ProtocolBuffer/blog/blogpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Hello I am client")
	conn, err := grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v",err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	createBlog(c)
	//readBlog(c)
}

func readBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Reading a Blog")
	readBlog,err := c.ReadBlog(context.Background(),&blogpb.ReadBlogRequest{BlogId: "000000000000000000000000"})
	if err != nil {
		log.Fatalf("Unexpected Error in Reading: %v",err)
	}
	fmt.Printf("Blog has been created: %v",readBlog)
}

func createBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Creating a Blog")
	blog := &blogpb.Blog{
		AuthorId: "StephaneOne",
		Title:    "My First Blog 2",
		Content:  "Content of the s First Blog",
	}
	createBlog,err := c.CreateBlog(context.Background(),&blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected Error: %v",err)
	}
	fmt.Printf("Blog has been created: %v",createBlog)
}
