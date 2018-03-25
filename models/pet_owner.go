package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type PetOwner struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	PersonID  uuid.UUID `json:"person_id" db:"person_id"`
	PetID     uuid.UUID `json:"pet_id" db:"pet_id"`
	Pet       Pet       `belongs_to:"pets"`
	Person    Person    `belongs_to:"people"`
}

// String is not required by pop and may be deleted
func (p PetOwner) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PetOwners is not required by pop and may be deleted
type PetOwners []PetOwner

// String is not required by pop and may be deleted
func (p PetOwners) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PetOwner) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PetOwner) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PetOwner) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
