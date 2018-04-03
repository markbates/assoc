package main

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/assoc/models"
	"github.com/stretchr/testify/require"
)

var DB *pop.Connection

func init() {
	var err error
	DB, err = pop.Connect("development")
	if err != nil {
		log.Fatal(err)
	}
	// pop.Debug = true
}

// saving a valid model with a valid
// nested assocation should save correctly.
func Test_Create_NestedModel(t *testing.T) {
	r := assert{require.New(t)}
	r.rollback(func() {

		code := &models.CourseCode{
			Course: models.Course{},
		}

		r.NoError(DB.Eager().Create(code))
		r.NotEqual(uuid.Nil, code.ID)
	})
}

// saving a valid model with a valid
// nested assocation should save correctly.
func Test_Create_NestedModel_Variation(t *testing.T) {
	r := assert{require.New(t)}
	r.rollback(func() {
		mark := &models.Person{
			Name: "Mark",
			Pets: models.Pets{
				models.Pet{Name: "Ringo"},
			},
		}

		r.NoError(DB.Eager().Create(mark))
		r.NotEqual(uuid.Nil, mark.ID)

		r.NoError(DB.Eager().First(mark))
		r.Len(mark.Pets, 1)
	})
}

// if an assocation is not valid then nothing should
// be saved
func Test_Create_ParentModel_Validation(t *testing.T) {
	r := assert{require.New(t)}
	r.rollback(func() {

		mark := &models.Person{
			Name: "",
			Pets: models.Pets{
				models.Pet{Name: "Ringo"},
			},
		}

		// person isn't valid, so the pet shouldn't be created
		verrs, err := DB.Eager().ValidateAndCreate(mark)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Zero(mark.ID)
	})
}

// if an assocation is not valid then nothing should
// be saved
func Test_Create_NestedModel_Validation(t *testing.T) {
	r := assert{require.New(t)}
	r.rollback(func() {
		r.delta(0, "people", func() {
			mark := &models.Person{
				Name: "Mark",
				Pets: models.Pets{
					models.Pet{Name: ""},
				},
			}

			// pets isn't valid, so the person shouldn't be created
			verrs, err := DB.Eager().ValidateAndCreate(mark)
			r.NoError(err)
			r.True(verrs.HasAny())
			r.Zero(mark.ID)
		})
	})
}

// it should save a simple model
func Test_Create_SimpleModel(t *testing.T) {
	r := assert{require.New(t)}
	r.rollback(func() {
		r.delta(1, "people", func() {
			mark := &models.Person{
				Name: "Mark",
			}

			r.NoError(DB.Eager().Create(mark))
			r.NotEqual(uuid.Nil, mark.ID)
		})
	})
}

func Test_Create_From_Join_Assoc(t *testing.T) {
	r := assert{require.New(t)}
	r.rollback(func() {
		r.delta(1, "pets", func() {
			r.delta(1, "people", func() {
				r.delta(1, "pet_owners", func() {
					mark := models.Person{
						Name: "Mark",
					}
					pet := models.Pet{Name: "Ringo"}

					po := &models.PetOwner{
						Person: mark,
						Pet:    pet,
					}
					r.NoError(DB.Eager().Create(po))

					r.NoError(DB.Eager().First(&mark))
					r.Len(mark.Pets, 1)
					r.NotEqual(uuid.Nil, mark.Pets[0].ID)
				})
			})
		})
	})
}

type assert struct {
	*require.Assertions
}

func (a *assert) rollback(fn func()) {
	defer func() {
		a.Nil(recover())
	}()
	DB.TruncateAll()
	defer DB.TruncateAll()
	fn()
}

func (r *assert) delta(d int, name string, fn func()) {
	sc, err := DB.Count(name)
	r.NoError(err)
	fn()
	ec, err := DB.Count(name)
	r.NoError(err)
	r.Equal(ec, sc+d)
}
