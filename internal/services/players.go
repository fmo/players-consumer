package services

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/fmo/players-consumer/internal/database"
	"github.com/fmo/players-consumer/internal/models"
	"github.com/sirupsen/logrus"
)

const tableName = "fmo-players"

type PlayersService struct {
	DB     *database.Database
	Logger *logrus.Logger
}

func NewPlayers(db *database.Database, l *logrus.Logger) PlayersService {
	return PlayersService{
		DB:     db,
		Logger: l,
	}
}

func (ps PlayersService) CreateOrUpdate(p models.Player) (response *dynamodb.PutItemOutput, err error) {
	playerParsed, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      playerParsed,
		TableName: aws.String(tableName),
	}

	return ps.DB.Connection.PutItem(input)
}

func (ps PlayersService) FindPlayersByTeamId(teamId int) (players []models.Player, err error) {
	filter := expression.Name("teamId").Equal(expression.Value(teamId))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return players, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	result, err := ps.DB.Connection.Scan(input)
	if err != nil {
		return players, err
	}

	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &players)
		if err != nil {
			return players, err
		} else {
			return players, nil
		}
	}

	return nil, errors.New("no result")
}

func (ps PlayersService) FindPlayersByRapidApiId(apiFootballId, transfermarktId string) (player *models.Player, err error) {
	var filter expression.ConditionBuilder
	if apiFootballId != "" {
		ps.Logger.Debugf("filter apiFootballId")
		filter = expression.Name("apiFootballId").Equal(expression.Value(apiFootballId))
	}
	if transfermarktId != "" {
		ps.Logger.Debugf("filter transfermarktId with transfermarktId %s", transfermarktId)
		filter = expression.Name("transfermarktId").Equal(expression.Value(transfermarktId))
	}

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	result, err := ps.DB.Connection.Scan(input)
	if err != nil {
		return nil, err
	}

	ps.Logger.Debugf("Found number of players: %d", len(result.Items))

	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &player)
		if err != nil {
			return nil, err
		} else {
			return player, nil
		}
	}

	return nil, errors.New("no result")
}

func (ps PlayersService) FindPlayerById(playerId string) (player *models.Player, err error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(playerId),
			},
		},
	}

	result, err := ps.DB.Connection.GetItem(input)
	if err != nil {
		return player, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("no player found with playerId: %s", playerId)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &player)
	if err != nil {
		return player, err
	}

	return player, nil
}
