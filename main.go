package main

import (
	"go-chat/app"
	"go-chat/controller"
	_ "go-chat/docs"
	"go-chat/pkg/websocket"
	"go-chat/repository"
	"go-chat/service"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Swagger Chat-App API
// @version 1.0
// @description This is a Chat-App server
// @termsOfService http://swagger.io/terms

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongo, err := repository.CreateMongoClient()
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	authRepository := repository.NewAuthRepository(mongo)
	authService := service.NewAuthService(authRepository)
	authController := controller.NewAuthController(authService)

	chatRepository := repository.NewChatRepository(mongo)
	chatService := service.NewChatService(chatRepository)
	chatController := controller.NewChatController(chatService)

	userRepository := repository.NewUserRepository(mongo)
	userService := service.NewUserService(authRepository, userRepository)
	userController := controller.NewUserController(userService)

	router := app.SetupRoutes(authController, chatController, userController)

	hub := websocket.NewHub()

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	handler := c.Handler(router)

	http.Handle("/", handler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r, chatRepository)
	})

	go hub.Run()

	server := http.Server{
		Addr:    "localhost:8000",
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
