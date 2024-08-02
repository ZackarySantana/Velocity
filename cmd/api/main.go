package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/service/domain"
	mongodomain "github.com/zackarysantana/velocity/internal/service/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := fmt.Sprintf(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	mux := api.New(domain.NewService(mongodomain.NewMongoRepository(client)))
	slog.Info("Starting server", "addr", "0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", mux)
}
