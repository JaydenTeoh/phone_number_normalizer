package main

import (
	"bytes"
	"database/sql"
	"fmt"

	phoneDB "github.com/JaydenTeoh/phone-number-normalizer/pkg/db"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "jaydenteoh"
	password = "password"
	dbname   = "phone_normalizer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(phoneDB.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phoneDB.Migrate("postgres", psqlInfo))

	db, err := phoneDB.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	if err := db.Seed(); err != nil {
		panic(err)
	}

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(p.Id))
			} else {
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
func normalize(phone string) string {
	var buf bytes.Buffer
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			buf.WriteRune(char)
		}
	}
	return buf.String()
}
