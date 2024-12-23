package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"sync"

	proto "github.com/ekanshthakur15/vapusdata/protoc"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{
	proto.UnimplementedBookStoreServer
}

var (
	books = make(map[string]*proto.Book)
	mu sync.Mutex
)

func main(){
	listener ,  tcpErr := net.Listen("tcp", ":9000")
	if tcpErr != nil{
		panic(tcpErr)
	}

	srv := grpc.NewServer()
	proto.RegisterBookStoreServer(srv, &server{})

	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (s *server) CreateBook(c context.Context, req *proto.CreateBookRequest) (*proto.CreateBookResponse, error){
	fmt.Println("recieved the request to create a book")
	mu.Lock()

	defer mu.Unlock()

	newID := uuid.New().String()

	book := req.Book
	book.Id = newID

	books[newID] = book

	log.Printf("Book created: %v\n", newID)


	return &proto.CreateBookResponse{ Id: newID}, nil
}

func (s *server) GetBook(c context.Context, req *proto.GetBookRequest) (*proto.GetBookResponse, error) {
	fmt.Println("Recieved the request to get a book.")

	id := req.Id
	if(id == ""){
		return nil, errors.New("provide the id")
	}

	book, exists := books[id]

	if(!exists) {
		return nil, errors.New("book doesn't exist")
	}

	fmt.Println(book)

	return &proto.GetBookResponse{Book: book} , nil
}

func (s *server) DeleteBook(c context.Context, req *proto.DeleteBookRequest) (*proto.DeleteBookResponse, error) {

	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Received the request to DELETE a book")

	id := req.Id
	if id == "" {
		return &proto.DeleteBookResponse{Success: false}, errors.New("provide an id")
	}

	delete(books, id)

	return &proto.DeleteBookResponse{Success: true} , nil

}

func (s *server) ListBooks(c context.Context, req *proto.ListBooksRequest) (*proto.ListBooksResponse, error) {
	fmt.Println("Received the request to LIST all books.")

	var allBooks []*proto.Book

	for _ , book := range books {
		allBooks = append(allBooks, book)
	}

	return &proto.ListBooksResponse{Books: allBooks}, nil

}