package main

import (
	"log"

	"github.com/kr/pretty"
	"github.com/markbates/assoc/models"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

var DB *pop.Connection

func init() {
	var err error
	DB, err = pop.Connect("development")
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = true
}

func main() {
	if err := DB.TruncateAll(); err != nil {
		log.Fatal(err)
	}
	if err := populateDB(); err != nil {
		log.Fatal(err)
	}
	if err := loadPerson(); err != nil {
		log.Fatal(err)
	}
}

func loadPerson() error {
	p := &models.Person{}
	if err := DB.Eager().First(p); err != nil {
		return errors.WithStack(err)
	}
	pretty.Println("### p ->", p)
	pretty.Println("### p.Pets ->", p.Pets)
	return nil
}

func populateDB() error {
	return DB.Transaction(func(tx *pop.Connection) error {
		mark := &models.Person{Name: "Mark"}
		if err := tx.Create(mark); err != nil {
			return errors.WithStack(err)
		}

		rachel := &models.Person{Name: "Rachel"}
		if err := tx.Create(rachel); err != nil {
			return errors.WithStack(err)
		}

		pet := &models.Pet{Name: "Ringo"}
		if err := tx.Create(pet); err != nil {
			return errors.WithStack(err)
		}

		po := &models.PetOwner{
			PersonID: mark.ID,
			PetID:    pet.ID,
		}
		if err := tx.Create(po); err != nil {
			return errors.WithStack(err)
		}

		po = &models.PetOwner{
			PersonID: rachel.ID,
			PetID:    pet.ID,
		}
		if err := tx.Create(po); err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}
