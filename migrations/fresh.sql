USE local;

SET GLOBAL event_scheduler = ON;
DROP TABLE IF EXISTS upvotes;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;

CREATE TABLE users(
    fprt VARCHAR(64) NOT NULL,
    name VARCHAR(40) NOT NULL,
    PRIMARY KEY(fprt)
);
CREATE TABLE chats(
    id BIGINT NOT NULL AUTO_INCREMENT,
    uuid VARCHAR(36) NOT NULL,
    parent_uuid VARCHAR(36) DEFAULT NULL,
    user_fprt VARCHAR(64) NOT NULL,
    msg VARCHAR(12000) NOT NULL,
    upvotes BIGINT NOT NULL DEFAULT 0,
    created_at DATETIME,
    PRIMARY KEY(uuid),
    INDEX(id),
    INDEX(parent_uuid),
    FOREIGN KEY (user_fprt) REFERENCES users(fprt) ON DELETE CASCADE
);
CREATE TABLE upvotes(
    chat_uuid VARCHAR(36) NOT NULL,
    user_fprt VARCHAR(64) NOT NULL,
    PRIMARY KEY(chat_uuid, user_fprt),
    FOREIGN KEY (chat_uuid) REFERENCES chats(uuid) ON DELETE CASCADE,
    FOREIGN KEY (user_fprt) REFERENCES users(fprt) ON DELETE CASCADE
);

DROP EVENT IF EXISTS delete_old_chats;
CREATE EVENT delete_old_chats ON SCHEDULE EVERY 1 MINUTE DO DELETE FROM chats WHERE created_at < NOW() - INTERVAL 1 HOUR;