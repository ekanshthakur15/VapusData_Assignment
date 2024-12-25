package main

import (
	"context"
	"net/http"

	proto "github.com/ekanshthakur15/vapusdata/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var client proto.BookStoreClient

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client = proto.NewBookStoreClient(conn)

	r := gin.Default()
	r.Use(MetadataMiddleware()) // Attach middleware before routes

	// Routes
	r.GET("/getBook", getBookController)
	r.GET("/listBooks", listBooksController)
	r.POST("/createBook", createBookController)
	r.DELETE("/deleteBook", deleteBookController)
	r.PUT("/updateBook", updateBookController)

	r.POST("/signup", signUpController)
	r.POST("/login", loginController)

	r.Run(":8000")
}

// Controllers
func createBookController(c *gin.Context) {
	var newBook proto.Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, exists := getGRPCContext(c)
	if !exists {
		return
	}

	req := &proto.CreateBookRequest{Book: &newBook}
	res, err := client.CreateBook(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": res.Id})
}

func getBookController(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing book ID"})
		return
	}

	ctx, exists := getGRPCContext(c)
	if !exists {
		return
	}

	req := &proto.GetBookRequest{Id: id}
	res, err := client.GetBook(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": res.Book})
}

func deleteBookController(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing book ID"})
		return
	}

	ctx, exists := getGRPCContext(c)
	if !exists {
		return
	}

	req := &proto.DeleteBookRequest{Id: id}
	res, err := client.DeleteBook(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": res.Success})
}

func listBooksController(c *gin.Context) {
	ctx, exists := getGRPCContext(c)
	if !exists {
		return
	}

	req := &proto.ListBooksRequest{}
	res, err := client.ListBooks(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": res.Books})
}

func updateBookController(c *gin.Context) {
	var updatedBook proto.Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, exists := getGRPCContext(c)
	if !exists {
		return
	}

	req := &proto.UpdateBookRequest{Book: &updatedBook}
	res, err := client.UpdateBook(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": res.Success})
}

func signUpController(c *gin.Context) {
	var signupRequest proto.CreateUserRequest
	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := client.CreateUser(context.Background(), &signupRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signup successful", "user_id": res.UserId})
}

func loginController(c *gin.Context) {
	var loginRequest proto.AuthenticationUserRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := client.Authentication(context.Background(), &loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": res.Token})
}

// Middleware
func MetadataMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip middleware for signup and login routes
		if c.Request.URL.Path == "/signup" || c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		md := metadata.New(map[string]string{
			"authorization": token,
		})
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		c.Set("grpcContext", ctx)

		c.Next()
	}
}

func getGRPCContext(c *gin.Context) (context.Context, bool) {
	ctxValue, exists := c.Get("grpcContext")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC context not found"})
		return nil, false
	}
	return ctxValue.(context.Context), true
}
