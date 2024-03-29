package main

import (
	"context"
	"github.com/Ad3bay0c/gRPC/blog/blogpb"
	"google.golang.org/grpc"
	"io"
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

func UpdateBlog(c blogpb.BlogServiceClient) {
	req := &blogpb.UpdateBlogRequest{
		BlogId: "61b8cc381dd789abc7ca95c4",
		Blog: &blogpb.Blog{},
	}
	_, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		log.Printf("Error while updating blog: %v\n", err)
	}
	req = &blogpb.UpdateBlogRequest{
		BlogId: "61b8cc381dd789abc7ca95c3",
		Blog: &blogpb.Blog{
			AuthorId: "Changed Author 2",
		},
	}
	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		log.Printf("Error while updating blog: %v\n", err)
		return
	}
	log.Printf("Blog was updated: %v", res)
}
func DeleteBlog(c blogpb.BlogServiceClient) {
    deleteRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
        BlogId: "61b8cc381dd789abc7ca95c3",
    })
    if err != nil {
        log.Printf("Error while deleting blog: %v\n", err)
		return
    }
    log.Printf("Blog was deleted: %v", deleteRes)
}

func ListBlog(c blogpb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
        log.Printf("Error while calling ListBlog RPC: %v\n", err)
        return
    }
	for {
		res, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Printf("Error while reading stream: %v\n", err)
            return
        }
        log.Printf("Blog was read: %v\n", res)
	}
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
	createBlog(c)

	// Read Blog
	ReadBlog(c)

	//Update Blog
	UpdateBlog(c)

	//Delete Blog
	//DeleteBlog(c)

	// List Blogs
	ListBlog(c)
}
