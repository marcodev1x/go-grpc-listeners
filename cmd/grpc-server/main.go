package main

import (
	"database/sql"
	"fmt"
	"net"

	"github.com/marcodev1x/grpc-tests/internal"
	"github.com/marcodev1x/grpc-tests/internal/pb"
	"github.com/marcodev1x/grpc-tests/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	schema := `
		CREATE TABLE IF NOT EXISTS categories (
    	id TEXT PRIMARY KEY,
    	name TEXT NOT NULL,
    	description TEXT
	);
	`

	if _, err := db.Exec(schema); err != nil {
		panic(fmt.Errorf("failed to create categories table: %w", err))
	}

	category := database.NewCategory(db)
	categoryService := service.NewCategoryService(category)

	// Registra o serviço no server grpc
	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":50051")

	fmt.Println("Started server")

	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
