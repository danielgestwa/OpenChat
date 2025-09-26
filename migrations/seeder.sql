USE local;

INSERT INTO users VALUES('1234bot', 'TEST MONKEY 1');
INSERT INTO users VALUES('2345bot', 'TEST MONKEY 2');
INSERT INTO users VALUES('3456bot', 'TEST MONKEY 3');
INSERT INTO chats VALUES(NULL, UUID(), NULL, '1234bot', 'Test message zxy', 0, NOW());
INSERT INTO chats VALUES(NULL, UUID(), NULL, '1234bot', 'Test message xyz', 0, NOW());
INSERT INTO chats VALUES(NULL, '3e72cf70-afb9-403b-9a9d-acf86ad41c8d', NULL, '1234bot', 'Test message yzx', 2, NOW());
INSERT INTO chats VALUES(NULL, UUID(), '3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '2345bot', 'Test response yzx', 0, NOW());
INSERT INTO chats VALUES(NULL, UUID(), '3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '2345bot', 'Test response xyz', 0, NOW());
INSERT INTO upvotes VALUES('3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '2345bot');
INSERT INTO upvotes VALUES('3e72cf70-afb9-403b-9a9d-acf86ad41c8d', '3456bot');