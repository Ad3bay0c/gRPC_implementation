package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
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
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthorID string             `json:"author_id" bson:"author_id,omitempty"`
	Content  string             `json:"content" bson:"content,omitempty"`
	Title    string             `json:"title" bson:"title,omitempty"`
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

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogID := req.GetBlogId()

	oID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot convert to ObjectID: %v", err.Error())
	}
	blog := &BlogItem{}
	err = collection.FindOne(context.Background(), bson.M{"_id": oID}).Decode(blog)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "An internal Server Error: %v", err.Error())
	}
	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(blog),
	}, nil
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blogID := req.GetBlogId()
	blog := req.GetBlog()
	data := &BlogItem{}
	oID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot convert to ObjectID: %v", err.Error())
	}
	err = collection.FindOne(context.Background(), bson.M{"_id": oID}).Decode(data)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ID cannot be found: %v", err.Error())
	}
	blogItem := &BlogItem{
		AuthorID: blog.AuthorId,
		Content: blog.Content,
		Title: blog.Title,
	}

	//_, err = collection.ReplaceOne(context.Background(), bson.M{"_id": oID}, blogItem)
	err = collection.FindOneAndUpdate(context.Background(), bson.M{"_id": oID}, bson.D{{"$set", *blogItem}}).Decode(data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "An internal Server Error: %v", err.Error())
	}
	blogItem.ID = oID
	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(blogItem),
	}, nil
}
func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, "An internal Server Error: %v", err.Error())
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		data := &BlogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, "An internal Server Error: %v", err.Error())
		}
		stream.Send(&blogpb.ListBlogResponse{
			Blog: &blogpb.Blog{
				Id: data.ID.Hex(),
                AuthorId: data.AuthorID,
                Title: data.Title,
                Content: data.Content,
            },
		},
		)
		if err := cur.Err(); err != nil {
			return status.Errorf(codes.Internal, "An internal Server Error: %v", err.Error())
		}
	}
	return nil
}
func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogID := req.GetBlogId()
	oID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot convert to ObjectID: %v", err.Error())
	}
	res, err := collection.DeleteOne(context.Background(), bson.M{"_id": oID})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ID cannot be found: %v", err.Error())
	}
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "ID cannot be found: %v", err.Error())
	}
	return &blogpb.DeleteBlogResponse{
        BlogId: blogID,
    }, nil
}
func dataToBlogPb(data *BlogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}
}
func main() {
	log.Println("Blog Server started")

	// if we crash the code, we get the file name and line number
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Llongfile)

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

	reflection.Register(s)
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
