package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	errors "github.com/rotisserie/eris"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/RBAC-server/api"
	"github.com/skyrocketOoO/RBAC-server/config"
	"github.com/skyrocketOoO/RBAC-server/internal/infra/graph"
	"github.com/skyrocketOoO/RBAC-server/internal/infra/mongo"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339
	// human-friendly logging without efficiency
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Logger initialized")

	if err := config.ReadConfig(); err != nil {
		log.Fatal().Msg(errors.ToString(err, true))
	}

	mongoClient, disconnectDb, err := mongo.InitDb()
	if err != nil {
		log.Fatal().Msg(errors.ToString(err, true))
	}
	defer disconnectDb()

	var graphInfra domain.GraphInfra
	graphInfra = graph.NewGraphInfra(dbRepo)
	usecase := usecase.NewUsecase(dbRepo, graphInfra)
	delivery := rest.NewDelivery(usecase)

	router := gin.Default()
	router.Use(middleware.CORS())
	api.Binding(router, delivery)

	router.Run(":8081")
}
