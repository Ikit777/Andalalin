package utils

import (
	"github.com/google/uuid"
)

type Notification struct {
	IdUser  uuid.UUID
	Status  string
	Tanggal string
	Kode    string
	Token   []string
}

func sendPushNotifications(data *Notification) {

}
