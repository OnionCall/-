package services

import (
	"math/rand"
	"strconv"

	"github.com/google/uuid"
)

func GenerateUuid() uuid.UUID {
    uuid := uuid.New()
    return uuid
}

func GenerateDefaultName() string {
    number := rand.Intn(9999)
	username := "User" + strconv.Itoa(number)

	return username
}