DROP DATABASE IF EXISTS bytedance_envelope;

CREATE DATABASE bytedance_envelope;

USE bytedance_envelope;

CREATE TABLE user(
    uid BIGINT NOT NULL,
    envelope_list VARCHAR(1000),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(uid)
);
INSERT INTO user(`uid`, `envelope_list`) VALUES (1, "1,2,3");
INSERT INTO user(`uid`, `envelope_list`) VALUES (2, "4");
INSERT INTO user(`uid`, `envelope_list`) VALUES (3, "5,6");
CREATE TABLE user_envelope(
    uid BIGINT NOT NULL,
    envelope_id BIGINT NOT NULL
);

CREATE TABLE envelope(
    envelope_id BIGINT NOT NULL,
    open_stat BOOLEAN,
    value INT NOT NULL,
    snatched_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(envelope_id)
);
INSERT INTO envelope(`envelope_id`, `open_stat`, `value`) values (1, 1, 50);
INSERT INTO envelope(`envelope_id`, `open_stat`, `value`) values (2, 0, 515130);
INSERT INTO envelope(`envelope_id`, `open_stat`, `value`) values (3, 0, 5460);
INSERT INTO envelope(`envelope_id`, `open_stat`, `value`) values (4, 1, 5044566);
INSERT INTO envelope(`envelope_id`, `open_stat`, `value`) values (5, 0, 54560);
INSERT INTO envelope(`envelope_id`, `open_stat`, `value`) values (6, 1, 7897950);
/*Future table below*/
/*
CREATE TABLE pool_list(
    pool_id INT NOT NULL,
    total_amount INT NOT NULL,
    remain_amount INT NOT NULL,
    total_money INT NOT NULL,
    remain_money INT NOT NULL,
    hold BOOLEAN
);
*/