package transport

import (
	"encoding/json"
	"log"
	"errors"
)

type BaseRequest struct {
	Action string `json:"action"`
	GameRequest json.RawMessage `json:"gameRequest"`
}

type Request interface {
	Validate() error
}

type CreateGameRequest struct {
	PlayerName string `json:"playerName"`
}

type JoinGameRequest struct {
	PlayerName string `json:"playerName"`
	GameId string `json::"gameId"`
}

func (request *CreateGameRequest) Validate() error {
	return nil
}

func (request *JoinGameRequest) Validate() error {
	return nil
}

func NewBaseRequest(apiGatewayRequestBody string) (*BaseRequest, error) {
	request := BaseRequest{}
	if err := json.Unmarshal([]byte(apiGatewayRequestBody), &request); err != nil {
		log.Println("Unable to decode request", err.Error())
		return nil, err
	}
	return &request, nil
}

func NewGameRequest(baseRequest *BaseRequest) (Request, error) {
	unmarshalGameRequest := func(gameRequest Request) (Request, error) {
		if err := json.Unmarshal([]byte(baseRequest.GameRequest), &gameRequest); err != nil {
			log.Println("Unable to unmarshal game request", err.Error())
			return nil, err
		}
		return gameRequest, nil
	}

	switch action := baseRequest.Action; action {
		case "create_game":
			return unmarshalGameRequest(&CreateGameRequest{})
		case "join_game":
			return unmarshalGameRequest(&JoinGameRequest{})
		default:
			return nil, errors.New("player action %v is not supported")
	}
}