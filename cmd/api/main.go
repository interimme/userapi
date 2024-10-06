package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/interimme/userapi/internal/controller"
	"github.com/interimme/userapi/internal/grpcserver"
	"github.com/interimme/userapi/internal/infrastructure"
	"github.com/interimme/userapi/internal/infrastructure/db"
	"github.com/interimme/userapi/internal/infrastructure/persistence"
	"github.com/interimme/userapi/internal/usecase"
	userapi "github.com/interimme/userapi/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Database connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		"db",       // DB_HOST
		"postgres", // DB_USER
		"postgres", // DB_PASSWORD
		"usersdb",  // DB_NAME
		"5432",     // DB_PORT
	)

	// Connect to the database
	dbConn, err := db.Connect(dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Retrieve the underlying *sql.DB to manage the connection pool
	sqlDB, err := dbConn.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Migrate the database schema
	err = dbConn.AutoMigrate(&persistence.UserGorm{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Initialize repository, use case, and controller
	userRepo := persistence.NewUserRepository(dbConn)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userController := controller.NewUserController(userUseCase)

	// Initialize Gin router
	router := infrastructure.NewRouter(userController)

	// Set up gRPC server
	grpcAddr := ":9090" // Define your gRPC port
	grpcServer := grpc.NewServer()
	grpcSrv := grpcserver.NewServer(userUseCase)
	userapi.RegisterUserServiceServer(grpcServer, grpcSrv)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Set up gRPC-Gateway
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()} // Insecure for simplicity; use TLS in production
	err = userapi.RegisterUserServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register gRPC-Gateway: %w", err)
	}

	// Listen on gRPC port
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcAddr, err)
	}

	// Set up HTTP server for gRPC-Gateway
	httpAddr := ":8080" // HTTP gateway port
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: gwMux,
	}

	// Create WaitGroup and channels for error handling
	var wg sync.WaitGroup
	errc := make(chan error, 3) // Buffer size 3 for gRPC, HTTP-Gateway, and Gin

	// Start gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("gRPC server listening on %s", grpcAddr)
		if err := grpcServer.Serve(grpcListener); err != nil {
			errc <- fmt.Errorf("gRPC server failed: %w", err)
		}
	}()

	// Start HTTP server for gRPC-Gateway
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("gRPC-Gateway HTTP server listening on %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- fmt.Errorf("HTTP gateway failed: %w", err)
		}
	}()

	// Start Gin HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		ginAddr := ":8000"
		log.Printf("Gin HTTP server listening on %s", ginAddr)
		if err := router.Run(ginAddr); err != nil {
			errc <- fmt.Errorf("Gin HTTP server failed: %w", err)
		}
	}()

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("Received signal: %v. Shutting down...", sig)

		// Gracefully stop gRPC server
		grpcServer.GracefulStop()

		// Shutdown HTTP server for gRPC-Gateway
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP gateway shutdown failed: %v", err)
		}

		// TODO: Graceful shutdown for Gin HTTP server
	}()

	// Wait for servers or apperrors
	select {
	case err := <-errc:
		return err
	case <-time.After(time.Hour): // Placeholder to prevent blocking; replace as needed
	}

	wg.Wait()
	return nil
}
