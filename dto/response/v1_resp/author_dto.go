package v1resp

import "github.com/google/uuid"

type AuthorDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
