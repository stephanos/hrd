package entity

import "time"

// CreateTimestamper manages an entity's creation time.
type CreateTimestamper interface {

	// CreatedAt returns the creation time.
	CreatedAt() time.Time

	// SetCreatedAt sets the creation time.
	SetCreatedAt(time.Time)
}

// CreatedTime implements the CreateTimestamper.
// It adds and manages an indexed creation time field.
type CreatedTime struct {
	EntityCreatedAt time.Time `datastore:"created_at,index"`
}

// CreatedAt returns the entity's creation time.
func (mdl *CreatedTime) CreatedAt() time.Time {
	return mdl.EntityCreatedAt
}

// SetCreatedAt sets the entity's creation time.
func (mdl *CreatedTime) SetCreatedAt(t time.Time) {
	if mdl.EntityCreatedAt.IsZero() {
		mdl.EntityCreatedAt = t
	}
}

// UpdateTimestamper manages an entity's last update time.
type UpdateTimestamper interface {

	// UpdatedAt returns the last update time.
	UpdatedAt() time.Time

	// SetUpdatedAt sets the last update time.
	SetUpdatedAt(time.Time)
}

// UpdatedTime implements the UpdateTimestamper.
// It adds and manages an indexed last updated time field.
type UpdatedTime struct {
	EntityUpdatedAt time.Time `datastore:"updated_at,index"`
}

// UpdatedAt returns the entity's last update time.
func (mdl *UpdatedTime) UpdatedAt() time.Time {
	return mdl.EntityUpdatedAt
}

// SetUpdatedAt sets the entity's last update time.
func (mdl *UpdatedTime) SetUpdatedAt(t time.Time) {
	mdl.EntityUpdatedAt = t
}
