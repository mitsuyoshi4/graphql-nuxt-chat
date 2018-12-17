package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/websocket"
	"github.com/mitsuyoshi4/graphqlchat"
	"github.com/rs/cors"
)

const defaultPort = "8080"
const defaultRedis = "127.0.0.1"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = defaultPort
	}

	c := cors.New(cors.Options{
		//AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	//http.Handle("/query", handler.GraphQL(graphqlchat.NewExecutableSchema(graphqlchat.Config{Resolvers: &graphqlchat.Resolver{}})))
	http.Handle("/query", c.Handler(handler.GraphQL(graphqlchat.NewExecutableSchema(graphqlchat.New(redisURL)),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}))),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
