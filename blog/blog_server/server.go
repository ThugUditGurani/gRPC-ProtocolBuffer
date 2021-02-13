package main

import (
	"Github/gRPC-ProtocolBuffer/blog/blogpb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var collection *mongo.Collection
type server struct {

}

func (s server) ReadBlog(ctx context.Context, request *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := request.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,fmt.Sprintf("Cannot parse ID"))
	}
	//Create an empty struct
	data := &blogItem{}
	res :=collection.FindOne(context.Background(),oid)
	if err := res.Decode(data);err != nil {
		return nil, status.Errorf(codes.NotFound,fmt.Sprintf("Cannot find blogs with specifies ID: %v",err))
	}
	return &blogpb.ReadBlogResponse{Blog: &blogpb.Blog{
		Id:      data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}},nil
}

type blogItem struct {
	ID primitive.ObjectID `bson:"_id",omitempty`
	AuthorID string  `bson:"author_id"`
	Content string	`bson:"content"`
	Title string	`bson:"title"`
}

func (s server) CreateBlog(ctx context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := request.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	insertOne, error := collection.InsertOne(context.Background(),data)
	if error != nil {
		fmt.Printf("Error While Inserting")
		return nil, status.Errorf(codes.Internal,fmt.Sprintf("Internal Error : %v",error))
	}
	objectId , ok := insertOne.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal,fmt.Sprintf("Cannot Convert Error : %v",error))
	}
	return &blogpb.CreateBlogResponse{Blog: &blogpb.Blog{Id: objectId.Hex(),AuthorId: blog.AuthorId,Title: blog.Title,Content: blog.Content}}, error
}



func main() {
	fmt.Println("Hello World")
	//if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//MongoDB Connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {

	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil { log.Fatal(err) }

	collection = client.Database("mydb").Collection("blog")


	//Listen Server
	lis , err := net.Listen("tcp","0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v",err)
	}
	s := grpc.NewServer()
	blogpb.RegisterBlogServiceServer(s,&server{})

	go func(){
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v",err)
		}
	}()

	//Wait for Control C to exit
	ch := make(chan os.Signal,1)
	signal.Notify(ch,os.Interrupt)

	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("MongoDB Connection Closing")
	client.Disconnect(context.Background())
	fmt.Println("End of the Program")
}
