package domain

import (
	"github.com/james-mcallister/may/database"
)

type Domain struct {
	mayDB database.Database
}

func NewDomain() Domain {
	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}
	return Domain{
		mayDB: db,
	}
}

func (d Domain) Init() {
	db, err := d.mayDB.Connect()
	if err != nil {
		panic(err)
	}
	if err := database.InitDB(db); err != nil {
		panic(err)
	}
}
