SET GLOBAL event_scheduler = ON;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS upvotes;

CREATE TABLE users(
    fprt VARCHAR(64) NOT NULL,
    name VARCHAR(40) NOT NULL,
    PRIMARY KEY(fprt)
);
CREATE TABLE chats(
    uuid VARCHAR(36) NOT NULL,
    parent_uuid VARCHAR(36) DEFAULT NULL,
    user_fprt VARCHAR(64) NOT NULL,
    msg VARCHAR(250) NOT NULL,
    upvotes BIGINT NOT NULL DEFAULT 0,
    created_at DATETIME,
    PRIMARY KEY(uuid),
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

--SEED, EXAMPLE DATA
INSERT INTO users VALUES('1234bot', 'TEST MONKEY 1');
INSERT INTO users VALUES('2345bot', 'TEST MONKEY 2');
INSERT INTO users VALUES('3456bot', 'TEST MONKEY 3');
INSERT INTO chats VALUES(UUID(), NULL '1234bot', 'Test message zxy', 0, NOW());
INSERT INTO chats VALUES(UUID(), NULL, '1234bot', 'Test message xyz', 0, NOW());
INSERT INTO chats VALUES('3e72cf70-afb9-403b-9a9d-acf86ad41c8d', NULL, '1234bot', 'Test message yzx', 2, NOW());
INSERT INTO chats VALUES(UUID(), '3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '2345bot', 'Test response yzx', 0, NOW());
INSERT INTO chats VALUES(UUID(), '3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '2345bot', 'Test response xyz', 0, NOW());
INSERT INTO upvotes VALUES('3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '2345bot');
INSERT INTO upvotes VALUES('3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '3456bot');