package main

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/assoc/models"
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
	// if err := loadPerson(); err != nil {
	// 	log.Fatal(err)
	// }
	if err := loadPersonAssoc(); err != nil {
		log.Fatal(err)
	}
	// if err := loadPet(); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := loadPetOwner(); err != nil {
	// 	log.Fatal(err)
	// }
}

func populateDB() error {
	return DB.Transaction(func(tx *pop.Connection) error {
		code := &models.CourseCode{
			Course: models.Course{},
		}
		if err := tx.Eager().Create(code); err != nil {
			return errors.WithStack(err)
		}
		if code.ID == uuid.Nil {
			return errors.New("code/course was never saved and an error was never thrown")
		}

		mark := &models.Person{Name: "Mark"}
		if err := tx.Eager().Create(mark); err != nil {
			return errors.WithStack(err)
		}

		rachel := &models.Person{Name: "Rachel"}
		if err := tx.Eager().Create(rachel); err != nil {
			return errors.WithStack(err)
		}

		pet := &models.Pet{Name: "Ringo"}
		if err := tx.Eager().Create(pet); err != nil {
			return errors.WithStack(err)
		}

		po := &models.PetOwner{
			PersonID: mark.ID,
			PetID:    pet.ID,
		}
		// using tx.Eager().Create(po) returns an error
		if err := tx.Eager().Create(po); err != nil {
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

func loadPerson() error {
	fmt.Println("\n\n### --- LOAD PERSON --- ###\n\n")
	p := &models.Person{}
	if err := DB.Eager().First(p); err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("### p ->", p)
	fmt.Println("### p.Pets ->", p.Pets)
	fmt.Println("### p.PetOwners ->", p.PetOwners)
	return nil
}

func loadPersonAssoc() error {
	fmt.Println("\n\n### --- LOAD PERSON --- ###\n\n")
	p := &models.Person{}
	if err := DB.First(p); err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("### p ->", p)
	fmt.Println("### p.Pets ->", p.Pets)
	fmt.Println("### p.PetOwners ->", p.PetOwners)
	if err := DB.Load(p); err != nil {
		return errors.WithStack(err)
	}
	if p.Pets == nil {
		return errors.New("p.Pets should not be nil")
	}
	fmt.Println("### p.Pets ->", p.Pets)
	if p.PetOwners == nil {
		return errors.New("p.Pets should not be nil")
	}
	fmt.Println("### p.PetOwners ->", p.PetOwners)
	if err := DB.Load(&p.PetOwners); err != nil {
		return errors.WithStack(err)
	}
	if p.PetOwners == nil {
		return errors.New("p.Pets should not be nil")
	}
	fmt.Println("### p.PetOwners ->", p.PetOwners)
	return nil
}

func loadPet() error {
	fmt.Println("\n\n### --- LOAD PET --- ###\n\n")
	p := &models.Pet{}
	if err := DB.Eager().First(p); err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("### p ->", p)
	fmt.Println("### p.Owners ->", p.Owners)
	return nil
}

func loadPetOwner() error {
	fmt.Println("\n\n### --- LOAD PET OWNER --- ###\n\n")
	po := &models.PetOwner{}
	if err := DB.Eager().First(po); err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("### po ->", po)
	fmt.Println("### po.Pet ->", po.Pet)
	fmt.Println("### po.Person ->", po.Person)
	return nil
}
