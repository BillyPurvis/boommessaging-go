package uuid

import uuid "github.com/satori/go.uuid"

// CreateUUID Generates a UUID
func CreateUUID() uuid.UUID {

	uuid := uuid.Must(uuid.NewV4())

	return uuid
}
