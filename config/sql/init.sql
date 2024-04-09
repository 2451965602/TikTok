CREATE TABLE west.`user`
(
    `user_id`    bigint       NOT NULL AUTO_INCREMENT,
    `username`   varchar(255) NOT NULL,
    `password`   varchar(255) NOT NULL,
    `avatar_url` varchar(255) NOT NULL,
    `opt_secret` varchar(255) NOT NULL DEFAULT '',
    `created_at` timestamp    NOT NULL DEFAULT current_timestamp,
    `updated_at` timestamp    NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp,
    `deleted_at` timestamp NULL DEFAULT NULL,
    CONSTRAINT `user_id` PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4;

CREATE TABLE west.`video`
(
    `video_id`      bigint       NOT NULL AUTO_INCREMENT,
    `user_id`       bigint       NOT NULL,
    `video_url`     varchar(255) NOT NULL,
    `cover_url`     varchar(255) NOT NULL,
    `title`         varchar(255) NOT NULL,
    `description`   varchar(255) NOT NULL,
    `visit_count`   bigint       NOT NULL,
    `like_count`    bigint       NOT NULL,
    `comment_count` bigint       NOT NULL,
    `created_at`    timestamp    NOT NULL DEFAULT current_timestamp,
    `updated_at`    timestamp    NOT NULL ON UPDATE current_timestamp DEFAULT current_timestamp,
    `deleted_at`    timestamp NULL DEFAULT NULL,
    CONSTRAINT `video_id` PRIMARY KEY (`video_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100000000000000000 DEFAULT CHARSET=utf8mb4;

CREATE TABLE west.`comment`
(
    `video_id` varchar(255) NULL DEFAULT NULL,
    `root_id`  varchar(255) NULL DEFAULT NULL,
    `user_id`  varchar(255) NOT NULL,
    `comment_id` bigint       NOT NULL AUTO_INCREMENT,
    `content`    varchar(255) NOT NULL,
    `created_at` timestamp    NOT NULL DEFAULT current_timestamp,
    `updated_at` timestamp    NOT NULL ON UPDATE current_timestamp DEFAULT current_timestamp,
    `deleted_at` timestamp NULL DEFAULT NULL,
    CONSTRAINT `comment_id` PRIMARY KEY (`comment_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000000000000000000 DEFAULT CHARSET=utf8mb4;

CREATE TABLE west.`like`
(
    `user_id`  bigint NOT NULL,
    `video_id` bigint NULL DEFAULT NULL,
    `root_id`  bigint NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES user (user_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE west.`social`
(
    `user_id`    bigint NOT NULL,
    `to_user_id` bigint NOT NULL,
    `status`     bigint NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(user_id)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4

CREATE TABLE west.`messages`
(
    `msg_id`       bigint       NOT NULL AUTO_INCREMENT,
    `from_user_id` bigint       NOT NULL,
    `to_user_id`   bigint       NOT NULL,
    `content`      varchar(255) NOT NULL,
    `created_at`   timestamp    NOT NULL DEFAULT current_timestamp,
    `status`       bigint       NOT NULL,
    CONSTRAINT `msg_id` PRIMARY KEY (`msg_id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;