package seeds

import (
	"fmt"
	"simple-api-go/http/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductList struct {
	Id    uint64
	Name  string
	Price float64
}

func ProductSeeder(db *gorm.DB) {
	list := []ProductList{
		{
			Id:    1,
			Name:  "Buku",
			Price: 10000,
		},
		{
			Id:    2,
			Name:  "Pulpen",
			Price: 5000,
		},
		{
			Id:    3,
			Name:  "Laptop",
			Price: 10550000,
		},
		{
			Id:    4,
			Name:  "Gitar",
			Price: 550000,
		},
		{
			Id:    5,
			Name:  "Tas",
			Price: 230000,
		},
		{
			Id:    6,
			Name:  "Rokok",
			Price: 50000,
		},
		{
			Id:    7,
			Name:  "Lap",
			Price: 4000,
		},
		{
			Id:    8,
			Name:  "Iphone",
			Price: 12199999,
		},
		{
			Id:    9,
			Name:  "Monitor",
			Price: 3500000,
		},
		{
			Id:    10,
			Name:  "Mouse",
			Price: 220000,
		},
	}

	fmt.Println("Start seeding ProductSeeder")
	for _, v := range list {
		row := model.Product{
			Name:  v.Name,
			Price: v.Price,
		}
		row.ID = v.Id
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "price"}),
		}).Create(&row)
	}
	fmt.Println("Completed seeding ProductSeeder")
}
