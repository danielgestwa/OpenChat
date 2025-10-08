package storage

import (
	"database/sql"
	"errors"
)

type Chat struct {
	UUID         string `json:"uuid"`
	ParentUUID   string `json:"parent_uuid"`
	User         string `json:"user"`
	Message      string `json:"msg" form:"msg" binding:"required"`
	Upvotes      string `json:"upvotes"`
	Vanish       string `json:"vanish"`
	RepliesCount string `json:"replies_cnt"`
}

var vanishTimeMin string = "60"

func QueryChats() ([]Chat, error) {
	db := connect()
	query := `
		SELECT
			c.uuid,
			c.msg,
			u.name AS user,
			c.upvotes,
			(` + vanishTimeMin + ` - TIMESTAMPDIFF(MINUTE, c.created_at, NOW())) AS vanish,
			(SELECT COUNT(*) FROM chats cs WHERE cs.parent_uuid = c.uuid) AS replies_cnt
		FROM
			chats AS c 
		JOIN
			users AS u
			ON c.user_fprt = u.fprt
		WHERE c.parent_uuid IS NULL
		ORDER BY c.id DESC;
	`
	results, err := db.Query(query) //db.Query("SELECT c.uuid, c.msg, u.name as user, c.upvotes, (" + vanishTimeMin + " - TIMESTAMPDIFF(MINUTE, c.created_at, NOW())) as vanish FROM chats as c, users as u WHERE c.user_fprt = u.fprt AND parent_uuid IS NULL ORDER BY id DESC")
	chats, newErr := fillChats(results, err)
	exit(db)
	return chats, newErr
}

func QueryReplies(parentId string) ([]Chat, error) {
	db := connect()
	results, err := db.Query("SELECT c.uuid, c.msg, u.name as user, c.upvotes, ("+vanishTimeMin+" - TIMESTAMPDIFF(MINUTE, c.created_at, NOW())) as vanish, '0' as replies_cnt FROM chats as c, users as u WHERE c.user_fprt = u.fprt AND parent_uuid = ? ORDER BY id", parentId)
	replies, newErr := fillChats(results, err)
	exit(db)
	return replies, newErr
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
	insert, err := db.Query("INSERT INTO chats VALUES (NULL, UUID(), NULL, ?, ?, 0, NOW())", chat.User, chat.Message)

	if err != nil {
		panic(err.Error())
	}

	insert.Close()
	exit(db)
}

func AddReply(chat Chat) {
	db := connect()
	insert, err := db.Query("INSERT INTO chats VALUES (NULL, UUID(), ?, ?, ?, 0, NOW())", chat.ParentUUID, chat.User, chat.Message)

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
		err = results.Scan(&chat.UUID, &chat.Message, &chat.User, &chat.Upvotes, &chat.Vanish, &chat.RepliesCount)
		if err != nil {
			return nil, errors.New("ERROR DURING ROWS GATHERING")
		}

		chats = append(chats, chat)
	}
	results.Close()
	return chats, nil
}
