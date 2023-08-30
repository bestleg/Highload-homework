package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"otus-homework/internal/database"
	"otus-homework/internal/env"
)

type PostRequestBody struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Sex        bool   `json:"sex"`
	Password   string `json:"password"`
}

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
				return
			}
			age, err := strconv.Atoi(record[1])
			if err != nil {
				return
			}
			data := PostRequestBody{
				FirstName:  firstName,
				SecondName: secondName,
				Birthdate:  time.Now().AddDate(-age, 0, 0).Format("2006-01-02"),
				Biography:  gofakeit.BeerName(),
				City:       record[2],
				Sex:        gofakeit.Bool(),
				Password:   "slozhnoooo!!!!",
			}
			body, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			resp, err := http.Post("http://localhost:4444/user/register", "application/json", bytes.NewReader(body))
			defer resp.Body.Close()

			if err != nil {
				fmt.Println(err)
			} else if resp != nil {
				fmt.Println(resp.StatusCode)
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("response Body:", string(body))
			}
		}
	}

}
