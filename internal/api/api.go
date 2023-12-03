package api

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var UpgradeToWebsocket = upgrader.Upgrade

func GenerateAPIKey() (string, error) {
	keyBytes := make([]byte, 32)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", err
	}

	apiKey := base64.StdEncoding.EncodeToString(keyBytes)
	return apiKey, nil
}
