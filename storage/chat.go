package storage

import (
	"database/sql"
	"errors"
)

type Chat struct {
	UUID       string `json:"uuid"`
	ParentUUID string `json:"parent_uuid"`
	User       string `json:"user"`
	Message    string `json:"msg" form:"msg" binding:"required"`
	Upvotes    string `json:"upvotes"`
	Vanish     string `json:"vanish"`
}

var vanishTimeMin string = "60"

func QueryChats() ([]Chat, error) {
	db := connect()
	results, err := db.Query("SELECT c.uuid, c.msg, u.name as user, c.upvotes, (" + vanishTimeMin + " - TIMESTAMPDIFF(MINUTE, c.created_at, NOW())) as vanish FROM chats as c, users as u WHERE c.user_fprt = u.fprt ORDER BY created_at DESC")
	chats, newErr := fillChats(results, err)
	exit(db)
	return chats, newErr
}

func QueryReplies(parentId string) ([]Chat, error) {
	db := connect()
	results, err := db.Query("SELECT c.uuid, c.msg, u.name as user, c.upvotes, ("+vanishTimeMin+" - TIMESTAMPDIFF(MINUTE, c.created_at, NOW())) as vanish FROM chats as c, users as u WHERE c.user_fprt = u.fprt AND parent_uuid = ? ORDER BY created_at DESC", parentId)
	chats, newErr := fillChats(results, err)
	exit(db)
	return chats, newErr
}

func QueryChat(id string) (Chat, error) {
	db := connect()

	var chat Chat
	err := db.QueryRow("SELECT uuid, msg FROM chats WHERE uuid = ?", id).Scan(&chat.Upvotes)
	if err != nil {
		return Chat{}, errors.New("ROWS NOT FOUND")
	}

	return chat, nil
}

func AddChat(chat Chat) {
	db := connect()
	insert, err := db.Query("INSERT INTO chats VALUES (UUID(), NULL, ?, ?, 0, NOW())", chat.User, chat.Message)

	if err != nil {
		panic(err.Error())
	}

	insert.Close()
	exit(db)
}

func AddReply(chat Chat) {
	db := connect()
	insert, err := db.Query("INSERT INTO chats VALUES (UUID(), ?, ?, ?, 0, NOW())", chat.UUID, chat.User, chat.Message)

	if err != nil {
		panic(err.Error())
	}

	insert.Close()
	exit(db)
}

func fillChats(results *sql.Rows, err error) ([]Chat, error) {
	if err != nil {
		return nil, errors.New("ROWS NOT FOUND")
	}
	var chats []Chat
	for results.Next() {
		var chat Chat
		// for each row, scan the result into our tag composite object
		err = results.Scan(&chat.UUID, &chat.Message, &chat.User, &chat.Upvotes, &chat.Vanish)
		if err != nil {
			return nil, errors.New("ERROR DURING ROWS GATHERING")
		}

		chats = append(chats, chat)
	}
	results.Close()
	return chats, nil
}
