package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/blog/blogpb"
	"google.golang.org/grpc"
	"log"
)
func createBlog (c blogpb.BlogServiceClient) {
	blog := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "123456789",
			Title: "The stroy changer",
			Content: "Content of the blog",
		},
	}
	res, err := c.CreateBlog(context.Background(), blog)
	if err != nil {
		log.Printf("Error while creating blog: %v", err)
		return
	}
	log.Printf("Blog has been created: %v", res)
}

func ReadBlog(c blogpb.BlogServiceClient) {
	_, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "61b8cc381dd789abc7ca95c4",
	})
	if err != nil {
		log.Printf("Error while reading blog: %v\n", err)
	}
	blogID := "61b8cc381dd789abc7ca95c3"
	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})
	if err != nil {
		log.Printf("Error while reading blog: %v\n", err)
        return
	}
	log.Printf("Blog was read: %v", res)
}
func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Llongfile)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := blogpb.NewBlogServiceClient(conn)

	// Create Blog
	//createBlog(c)

	// Read Blog
	ReadBlog(c)
}
