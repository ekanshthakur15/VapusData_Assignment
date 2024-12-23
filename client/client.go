package main

import (
	"log"
	"net/http"

	proto "github.com/ekanshthakur15/vapusdata/protoc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client proto.BookStoreClient

func main(){
	conn , err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client = proto.NewBookStoreClient(conn)

	r := gin.Default()
	r.GET("/getBook",getBookController)
	r.GET("/listBooks", listBooksController)
	r.POST("/createBook", createBookController)
	r.DELETE("/deleteBook", deleteBookController)
	

	r.Run(":8000")

}

func createBookController(c *gin.Context) {

	var newBook proto.Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		
	}

	req := &proto.CreateBookRequest{Book: &newBook}

	res, err := client.CreateBook(c, req)
	if err != nil{
		log.Fatalln("error creating book: ", err)
	}
	
	log.Println("Book created Successfully with id: \n", res.Id)

	c.JSON(http.StatusCreated, gin.H{"id" : res.Id})
}

func getBookController(c *gin.Context) {
	id := c.Query("id")

	req := &proto.GetBookRequest{Id: id}

	res , err := client.GetBook(c, req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": res.Book})
}

func deleteBookController(c *gin.Context) {
	id := c.Query("id")

	req := &proto.DeleteBookRequest{Id: id}

	res , err := client.DeleteBook(c, req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": res.Success})
}
func listBooksController(c *gin.Context) {

	req := &proto.ListBooksRequest{}

	res , err := client.ListBooks(c, req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": res})

}