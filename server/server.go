package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"sync"
	"time"

	proto "github.com/ekanshthakur15/vapusdata/protoc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)
var jwtSecret []byte
type server struct {
    proto.UnimplementedBookStoreServer
    books  map[string]*proto.Book   // map[bookId]*Book
}

// var (
// 	books = make(map[string]*proto.Book)
// 	mu sync.Mutex
// )

type User struct {
	ID string
	Username string
	Password string
}


var (users map[string]*User
	tokens map[string]string
	mu sync.Mutex
)

func init() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Get JWT_SECRET from the environment
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatalf("JWT_SECRET is not set in the .env file")
    }
    jwtSecret = []byte(secret)
}

func main() {

	// Initialize the global users map
    users = make(map[string]*User)
    tokens = make(map[string]string)


	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor),
	)

	bookStoreServer := &server{
		books:  make(map[string]*proto.Book),
	}

	proto.RegisterBookStoreServer(srv, bookStoreServer)
	reflection.Register(srv)

	log.Println("gRPC server running on :9000")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) CreateBook(c context.Context, req *proto.CreateBookRequest) (*proto.CreateBookResponse, error){
	fmt.Println("recieved the request to create a book")


	newID := uuid.New().String()

	book := req.Book
	book.Id = newID

	s.books[newID] = book

	log.Printf("Book created: %v\n", newID)


	return &proto.CreateBookResponse{ Id: newID}, nil
}

func (s *server) GetBook(c context.Context, req *proto.GetBookRequest) (*proto.GetBookResponse, error) {
	fmt.Println("Recieved the request to get a book.")

	id := req.Id
	if(id == ""){
		return nil, errors.New("provide the id")
	}

	book, exists := s.books[id]

	if(!exists) {
		return nil, errors.New("book doesn't exist")
	}

	fmt.Println(book)

	return &proto.GetBookResponse{Book: book} , nil
}

func (s *server) DeleteBook(c context.Context, req *proto.DeleteBookRequest) (*proto.DeleteBookResponse, error) {

	fmt.Println("Received the request to DELETE a book")

	id := req.Id
	if id == "" {
		return &proto.DeleteBookResponse{Success: false}, errors.New("provide an id")
	}

	delete(s.books, id)

	return &proto.DeleteBookResponse{Success: true} , nil

}

func (s *server) ListBooks(c context.Context, req *proto.ListBooksRequest) (*proto.ListBooksResponse, error) {
	fmt.Println("Received the request to LIST all books.")

	var allBooks []*proto.Book

	for _ , book := range s.books {
		allBooks = append(allBooks, book)
	}

	return &proto.ListBooksResponse{Books: allBooks}, nil

}

func (s *server) UpdateBook(c context.Context, req * proto.UpdateBookRequest) (*proto.UpdateBookResponse, error) {
	fmt.Println("Received the request to update a books.")
	// mu.Lock()

	// defer mu.Unlock()

	book := req.Book

	if(book.Id == ""){
		return &proto.UpdateBookResponse{Success: false}, errors.New("provide the ID")
	}

	_ , exists := s.books[book.Id]
	if(!exists){
		return &proto.UpdateBookResponse{Success: false}, errors.New("book not found");
	}

	s.books[book.Id] = book

	return &proto.UpdateBookResponse{Success: true}, nil
}

func (s *server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
    log.Println("CreateBook: Locking mutex")
	// mu.Lock()
	// log.Println("CreateBook: Mutex locked")
	// defer func() {
	// 	mu.Unlock()
	// 	log.Println("CreateBook: Mutex unlocked")
	// }()


    for _, user := range users {
        if user.Username == req.Username {
            return nil, status.Error(codes.AlreadyExists, "username already exists")
        }
    }

    hashedPassword, err := hashPassword(req.Password)
    if err != nil {
        return nil, status.Error(codes.Internal, "error hashing password")
    }

    userID := uuid.New().String()
    newUser := &User{
        ID:       userID,
        Username: req.Username,
        Password: hashedPassword,
    }
    users[req.Username] = newUser


    log.Printf("User created: %s (ID: %s)", req.Username, userID)
    log.Println("Users: ", users[req.Username])
    return &proto.CreateUserResponse{UserId: userID}, nil
}


func (s *server) Authentication(c context.Context, req *proto.AuthenticationUserRequest) (*proto.AuthenticationUserResponse, error) {
	log.Println("CreateBook: Locking mutex")
	// mu.Lock()
	// log.Println("CreateBook: Mutex locked")
	// defer func() {
	// 	mu.Unlock()
	// 	log.Println("CreateBook: Mutex unlocked")
	// }()


	user, exists := users[req.Username]
	if !exists {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	if !validatePassword(req.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "error generating token")
	}

	tokens[user.ID] = token 

	return &proto.AuthenticationUserResponse{Token: token}, nil
}

func authInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    // Methods that don't require authentication

	// Log the full method name to debug
    // fmt.Println("gRPC Method called:", info.FullMethod)

    unauthenticatedMethods := map[string]bool{
        "/BookStore/CreateUser":       true,
        "/BookStore/Authentication": true,
    }

    // Check if the method is in the unauthenticated list
    if unauthenticatedMethods[info.FullMethod] {
        return handler(ctx, req)
    }

    // Perform authentication for other methods
    token := grpcHeader(ctx, "authorization")
	fmt.Println("received token: ", token)
    if token == "" {
        return nil, status.Error(codes.Unauthenticated, "missing authorization token")
    }

    userID, err := validateJWTToken(token)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, err.Error())
    }

	// Check if the user exists in the users map
    // mu.Lock()
    // defer mu.Unlock()

	userPresent := false
	for _, user := range users {
		if user.ID == userID {
			userPresent = true
			break
		} else {
			userPresent = false
		}
	}
	if !userPresent {
		return nil , status.Error(codes.Unauthenticated, "user does not exist")
	}

	if tokens[userID] != token {
		return nil , status.Error(codes.Unauthenticated, "invalid token")
	}

    newCtx := context.WithValue(ctx, "userID", userID)
    return handler(newCtx, req)
}



func hashPassword(password string) (string, error){
	bytes ,  err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}


func validatePassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func validateJWTToken(tokenString string) (string, error) {
	token , err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	fmt.Println("The token: ",tokenString)
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims , ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return claims.Subject, nil
}

func generateJWTToken(userId string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject: userId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(jwtSecret)
}

func grpcHeader(ctx context.Context, key string) string {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return ""
    }

    values := md[key]
	// log.Println("authorisation: ", values)
    if len(values) > 0 {
		// fmt.Println("Length is: ", len(values))
        tokenParts := strings.Split(values[0], " ")
		// log.Println("parts: ", tokenParts)
		// log.Println("len of parts: ", len(tokenParts))

        if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
			// log.Println(tokenParts[0])
			// log.Println(tokenParts[1])
            return tokenParts[1]
        }
    }
    return ""
}
