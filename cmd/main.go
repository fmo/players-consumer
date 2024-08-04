package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/fmo/football-proto/golang/player"
	"github.com/fmo/players-consumer/internal/database"
	"github.com/fmo/players-consumer/internal/kafka"
	"github.com/fmo/players-consumer/internal/models"
	"github.com/fmo/players-consumer/internal/services"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var logger = log.New()

func init() {
	logger.Out = os.Stdout

	logger.Level = log.DebugLevel
}

func main() {
	fmt.Println(fmt.Sprintf("Player consumer is up and running"))

	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		err := godotenv.Load()
		if err != nil {
			logger.Fatal("Error loading .env file")
		}
	}

	k := kafka.NewKafka()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*600)
	defer cancel()

	playersService := services.NewPlayers(
		database.NewDbAdapter(),
		logger,
	)

	msgNumber := 0

	for {
		message, err := k.Reader.ReadMessage(ctx)
		if err != nil {
			logger.Errorf("Error reading message: %v\n", err)
			continue
		}

		msgNumber++

		logger.Debugf("received the %d. message payload", msgNumber)

		if len(message.Value) == 0 {
			logger.Debugf("received an empty message")
			continue
		}

		var players []*pb.Player
		err = json.Unmarshal(message.Value, &players)
		if err != nil {
			logger.Errorf("Error unmarshalling message: %v", err)
			continue
		}

		debugPrint := ""
		for _, player := range players {
			debugPrint = fmt.Sprintf("%s, %s", debugPrint, player.Name)
		}

		logger.WithFields(log.Fields{
			"receivedMessage": fmt.Sprintf("%s...", debugPrint[2:100]),
		}).Debugf("unmarshalled the %d. message payload", msgNumber)

		for _, player := range players {
			p := models.Player{
				Team:                player.Team,
				TeamId:              player.TeamId,
				Name:                player.Name,
				Firstname:           player.Firstname,
				Lastname:            player.Lastname,
				Age:                 player.Age,
				Nationality:         player.Nationality,
				Photo:               player.Photo,
				ApiFootballId:       player.ApiFootballId,
				Appearances:         player.Appearances,
				Position:            player.Position,
				TransfermarktId:     player.TransfermarktId,
				ShirtNumber:         player.ShirtNumber,
				MarketValue:         player.MarketValue,
				MarketValueCurrency: player.MarketValueCurrency,
			}

			playerInDb, err := playersService.FindPlayersByRapidApiId(player.ApiFootballId, player.TransfermarktId)
			if err != nil {
				logger.Debugf("Finding player in db failed %v", err)
			}
			if playerInDb != nil {
				logger.WithFields(log.Fields{
					"playerInDbId": playerInDb.Id,
				}).Debug("Player already in the database")
				p.Id = playerInDb.Id
			} else {
				p.Id = uuid.New().String()
				logger.WithFields(log.Fields{
					"newPlayerId": p.Id,
				}).Info("New player will be created")
			}

			_, err = playersService.CreateOrUpdate(p)
			if err != nil {
				logger.Fatalf("Got error calling PutItem: %s", err)
			}

			logger.WithFields(log.Fields{
				"playerId": p.ApiFootballId,
				"teamName": p.Team,
			}).Infof("inserted or updated %s %s to the database", p.Firstname, p.Lastname)
		}
	}
	k.Reader.Close()
}
