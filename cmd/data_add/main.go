package main

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"time"

	"example.com/internal/database"
	"example.com/internal/env"
)

func main() {
	db, err := database.New(env.GetString("DB_DSN", "postgres:postgres@localhost:5432/postgres?sslmode=disable"), false)
	defer db.Close()
	if err != nil {
		return
	}

	file, err := os.Open("./cmd/data_add/people.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	for {
		record, e := reader.Read()
		if e == nil {
			secondName, firstName, found := strings.Cut(record[0], " ")
			if !found {
				continue
			}
			age, err := strconv.Atoi(record[1])
			if err != nil {
				continue
			}
			db.InsertToUserData(context.Background(), database.UserData{
				FirstName:  firstName,
				SecondName: secondName,
				Birthdate:  time.Now().AddDate(-age, 0, 0),
				Sex:        false,
				Biography:  "",
				City:       record[2],
			})
		}
	}

}
