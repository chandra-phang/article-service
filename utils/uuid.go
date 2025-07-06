package utils

import (
	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}

func GenerateUuidV4() uuid.UUID {
	return uuid.New()
}
