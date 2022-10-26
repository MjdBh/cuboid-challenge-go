package testutils

import (
	"fmt"

	"cuboid-challenge/app/db"
)

func ConnectDB() {
	db.Connect()
}

func ClearDB() {
	if err := db.CONN.Exec("DELETE FROM bags; DELETE FROM cuboids;").Error; err != nil {
		panic(fmt.Errorf("failed to ClearDB. %w", err))
	}
}

func AddRecords(records ...interface{}) {
	for _, record := range records {
		if r := db.CONN.Create(record); r.Error != nil {
			panic(fmt.Errorf("failed to AddRecords %w", r.Error))
		}
	}
}

func UpdateRecords(records ...interface{}) {
	for _, record := range records {
		if r := db.CONN.Save(record); r.Error != nil {
			panic(fmt.Errorf("failed to UpdateRecords. %w", r.Error))
		}
	}
}

func FindRecord(record interface{}, id interface{}) bool {
	r := db.CONN.First(record, id)

	return r.Error == nil
}
