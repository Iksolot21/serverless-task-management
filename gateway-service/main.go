package main

import (
  "fmt"
  "log"
  "net/http"
   "os"

   "github.com/gorilla/mux"
   "github.com/joho/godotenv"
     "your-repo/gateway-service/config"
    "your-repo/gateway-service/handlers"
      "your-repo/gateway-service/pb"
      "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
        "your-repo/gateway-service/middleware"
        "your-repo/gateway-service/internal/errors"
)
func main() {
    err := godotenv.Load()
    if err != nil {
      log.Fatalf("Error loading .env file: %v", err)
    }

     cfg, err := config.LoadConfig()
    if err != nil {
     log.Fatalf("Error loading config: %v", err)
   }

  authConn, err := grpc.Dial(cfg.AuthServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
      log.Fatalf("Did not connect: %v", err)
   }
    defer authConn.Close()
  taskConn, err := grpc.Dial(cfg.TaskServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
      log.Fatalf("Did not connect: %v", err)
  }
    defer taskConn.Close()

  userConn, err := grpc.Dial(cfg.UserServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
       log.Fatalf("Did not connect: %v", err)
    }
  defer userConn.Close()
    notificationConn, err := grpc.Dial(cfg.NotificationServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
      if err != nil {
      log.Fatalf("Did not connect: %v", err)
    }
  defer notificationConn.Close()
     r := mux.NewRouter()
       r.Use(middleware.CORSMiddleware(cfg.FrontendURL))
         // Auth Routes
        r.HandleFunc("/auth/register", handlers.RegisterUser(authConn)).Methods("POST")
        r.HandleFunc("/auth/login", handlers.LoginUser(authConn)).Methods("POST")
        r.HandleFunc("/me", middleware.AuthMiddleware(authConn, cfg, handlers.GetCurrentUser(userConn, cfg))).Methods("GET")
          // Task Routes
        r.HandleFunc("/tasks", middleware.AuthMiddleware(authConn, cfg, handlers.GetTasks(taskConn, cfg))).Methods("GET")
        r.HandleFunc("/tasks", middleware.AuthMiddleware(authConn, cfg, handlers.CreateTask(taskConn, cfg))).Methods("POST")
       r.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(authConn, cfg, handlers.GetTaskById(taskConn, cfg))).Methods("GET")
        r.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(authConn, cfg, handlers.PatchTaskById(taskConn, cfg))).Methods("PATCH")
       r.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(authConn, cfg, handlers.DeleteTaskById(taskConn, cfg))).Methods("DELETE")
           //User Routes
       r.HandleFunc("/users", handlers.GetUsers(userConn)).Methods("GET")
        r.HandleFunc("/users/{id}", handlers.GetUserById(userConn)).Methods("GET")
         // Notification Routes
        r.HandleFunc("/notification", handlers.SendNotification(notificationConn)).Methods("POST")
     log.Println(fmt.Sprintf("server starting at %s", cfg.Port ))
      err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
    if err != nil {
            errors.RespondWithError(w http.StatusInternalServerError, "Failed to start gateway")
            os.Exit(1)
       }
}