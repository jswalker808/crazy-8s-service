package game

type Connection struct {
	id string
	gameId string
}

func NewConnection() Connection {
	return Connection{}
}

func (connection Connection) GetId() string {
	return connection.id
}

func (connection Connection) GetGameId() string {	
	return connection.gameId
}