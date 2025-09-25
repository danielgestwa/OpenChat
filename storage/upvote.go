package storage

import "errors"

type Upvote struct {
	Count int `json:"cnt"`
}

func QueryUpvote(id string) (Chat, error) {
	db := connect()

	var chat Chat
	err := db.QueryRow("SELECT upvotes FROM chats WHERE uuid = ?", id).Scan(&chat.Upvotes)
	if err != nil {
		return Chat{}, errors.New("ROWS NOT FOUND")
	}

	return chat, nil
}

func AddUpvote(uuid string, userFprt string) {
	db := connect()

	var upvote Upvote
	err := db.QueryRow("SELECT count(1) as cnt FROM upvotes WHERE chat_uuid = ? AND user_fprt = ?", uuid, userFprt).Scan(&upvote.Count)
	if err != nil {
		panic(err.Error())
	}

	if upvote.Count == 0 {
		//Increase
		update, err := db.Query("UPDATE chats SET upvotes = upvotes + 1 WHERE uuid = ?", uuid)
		if err != nil {
			panic(err.Error())
		}
		update.Close()

		//Insert to upvotes
		insert, err := db.Query("INSERT INTO upvotes VALUES (?, ?)", uuid, userFprt)
		if err != nil {
			panic(err.Error())
		}
		insert.Close()
	}
	exit(db)
}
