package entity

import "time"

// CreateTimestamper can apply the time of creation to an entity.
type CreateTimestamper interface {

	// SetCreatedAt applies the time of creation.
	SetCreatedAt(time.Time)
}

// UpdateTimestamper can apply the time of last update to an entity.
type UpdateTimestamper interface {

	// SetUpdatedAt applies the time of the last update.
	SetUpdatedAt(time.Time)
}
