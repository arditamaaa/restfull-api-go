package seeds

import (
	"fmt"
	"simple-api-go/http/model"
	"simple-api-go/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserList struct {
	Id       uint64
	Name     string
	Email    string
	Password string
	Role     string
}

func UserSeeder(db *gorm.DB) {
	list := []UserList{
		{
			Id:       1,
			Name:     "admin",
			Email:    "admin@gmail.com",
			Password: "password123",
			Role:     "admin",
		},
		{
			Id:       2,
			Name:     "user",
			Email:    "user@gmail.com",
			Password: "password123",
			Role:     "user",
		},
	}

	fmt.Println("Start seeding UserSeeder")
	for _, v := range list {
		hashedPassword, _ := util.HashPassword(v.Password)
		row := model.User{
			Name:     v.Name,
			Email:    v.Email,
			Password: hashedPassword,
			Role:     v.Role,
		}
		row.ID = v.Id
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "email", "password"}),
		}).Create(&row)
	}
	fmt.Println("Completed seeding UserSeeder")
}
