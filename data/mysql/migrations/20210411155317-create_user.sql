-- +migrate Up
CREATE TABLE `user`
(
    `user_id` binary(16) NOT NULL,
    `email` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `role` smallint(2) NOT NULL,
    PRIMARY KEY (`user_id`),
    INDEX `user_id_index` (`user_id`),
    UNIQUE INDEX `user_email_index` (`email`)
);

-- +migrate Down
DROP TABLE `user`;