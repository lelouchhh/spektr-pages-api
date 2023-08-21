package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	_cityHttp "spektr-pages-api/city/delivery/http"
	_cityRepo "spektr-pages-api/city/repository/postgres"
	_cityUsecase "spektr-pages-api/city/usecase"
	_tariffHttp "spektr-pages-api/tariff/delivery/http"
	_tariffRepo "spektr-pages-api/tariff/repository/postgres"
	_tariffUsecase "spektr-pages-api/tariff/usecase"
	"syscall"
	"time"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbPort := viper.GetString("database.port")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	// Construct the connection string
	connection := fmt.Sprintf("postgres://%v:%v@db:%v/%v?sslmode=disable",
		dbUser,
		dbPass,
		dbPort,
		dbName)
	fmt.Println(connection)
	// Open a connection to the database
	dbConn, err := sqlx.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify the connection
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(err)
	defer func() {
		// Close the database connection
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	g := gin.Default()
	g.Static("/assets", "./static")
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	tariffRepo := _tariffRepo.NewTariffRepository(dbConn)
	tariffUcase := _tariffUsecase.NewTariffUsecase(tariffRepo, timeoutContext)
	_tariffHttp.NewTariffHandler(g, tariffUcase)
	cityRepo := _cityRepo.NewCityRepository(dbConn)
	cityUcase := _cityUsecase.NewCityUsecase(cityRepo, timeoutContext)
	_cityHttp.NewCityHandler(g, cityUcase)
	server := &http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: g,
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for a termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
