package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

type server struct{}

type BlogItem struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	AuthorID string             `json:"author_id" bson:"author_id"`
	Content  string             `json:"content" bson:"content"`
	Title    string             `json:"title" bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := req.GetBlog()
	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Content: blog.GetContent(),
		Title: blog.GetTitle(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Printf("An Error occurred while inserting to the database: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "Internal Error: %v", err.Error())
	}
	oID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "Cannot convert to ObjectID")
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oID.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func main() {
	log.Println("Blog Server started")

	// if we crash the code, we get the file name and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	if client.Ping(ctx, readpref.Primary()) != nil {
		log.Fatal("Could not connect to MongoDB")
		return
	}
	collection = client.Database("myDB").Collection("blog")

	log.Println("Database Connected Successfully")
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("An Error occurred: %v", err.Error())
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		if serverErr := s.Serve(listen); serverErr != nil {
			log.Fatalf("Error Connecting to port: %v", serverErr.Error())
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	s.Stop()
	log.Printf("Stopping Listening...")
	listen.Close()
	client.Disconnect(ctx)
	log.Printf("End of Blog Program")
}
