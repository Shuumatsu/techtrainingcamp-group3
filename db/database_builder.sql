DROP DATABASE IF EXISTS bytedance_envelope;

CREATE DATABASE bytedance_envelope;

USE bytedance_envelope;

CREATE TABLE user(
    uid BIGINT NOT NULL,
    envelope_list VARCHAR(1000),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(uid)
);

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