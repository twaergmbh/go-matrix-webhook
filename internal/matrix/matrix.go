package matrix

import (
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Matrix struct {
	client *mautrix.Client
}

func NewMatrix(serverUrl string, userID string, token string) (*Matrix, error) {
	c, err := mautrix.NewClient(serverUrl, id.NewUserID(userID, serverUrl), token)
	if err != nil {
		return &Matrix{}, err
	}

	return &Matrix{c}, nil
}

func (m Matrix) JoinRoom(roomID string) error {
	_, err := m.client.JoinRoom(roomID, "", nil)
	return err
}

func (m Matrix) SendMessage(roomID, message string) error {
	_, err := m.client.SendMessageEvent(id.RoomID(roomID), event.EventMessage, &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    message,
	})
	return err
}

func (m *Matrix) StartPrivateChat(userID string) (string, error) {
	room, err := m.client.CreateRoom(&mautrix.ReqCreateRoom{
		Visibility: "private",
		Invite:     []id.UserID{id.NewUserID(userID, m.client.HomeserverURL.Host)},
		IsDirect:   true,
		Preset:     "private_chat",
	})
	if err != nil {
		return "", err
	}
	return room.RoomID.String(), nil
}

func (m *Matrix) CreateOrFindPrivateChat(userID string) (string, error) {
	userIDObj := id.NewUserID(userID, m.client.HomeserverURL.Host)

	// Retrieve the direct chats mapping from account data
	var directMap map[string][]string

	err := m.client.GetAccountData("m.direct", &directMap)

	if err != nil {
		log.Error().Str("user_id", userIDObj.String()).Err(err).Msg("Failed to fetch m.direct data")
		// Inspect the error more closely
		if httpErr, ok := err.(*mautrix.HTTPError); ok {
			log.Error().Str("status_code", httpErr.RespError.ErrCode).Str("status", httpErr.Message).Msg("Detailed HTTP error")
		}
		return "", err
	}

	// Check if there is an existing room with the specified user
	//if roomIDs, found := directMap[userIDObj.String()]; found && len(roomIDs) > 0 {
	//	// Return the first room ID found
	//	return roomIDs[0], nil
	//}

	// No existing room found, create a new private chat
	return m.StartPrivateChat(userID)
}
