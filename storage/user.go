package storage

import (
	"time"

	"github.com/goombaio/namegenerator"
)

func UserExists(id string) {
	db := connect()
	var cnt int

	//TODO: Query user and check if exists
	err := db.QueryRow("SELECT count(1) as cnt FROM users WHERE fprt = ?", id).Scan(&cnt)
	if err != nil {
		panic(err.Error())
	}

	//Insert user if not exists
	if cnt == 0 {
		seed := time.Now().UTC().UnixNano()
		nameGenerator := namegenerator.NewNameGenerator(seed)
		name := nameGenerator.Generate()

		insert, err := db.Query("INSERT INTO users VALUES (?, ?)", id, name)
		if err != nil {
			panic(err.Error())
		}
		insert.Close()
	}
	exit(db)
}
