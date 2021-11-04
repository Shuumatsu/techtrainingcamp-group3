DROP DATABASE IF EXISTS bytedance_envelope;

CREATE DATABASE bytedance_envelope;

USE bytedance_envelope;

CREATE TABLE user(
  uid BIGINT NOT NULL,
  amount INT NOT NULL,
  envelope_list VARCHAR(1000),
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(uid)
);

CREATE TABLE envelope(
  envelope_id BIGINT NOT NULL,
  uid BIGINT NOT NULL,
  opened BOOLEAN,
  value INT NOT NULL,
  snatch_time BIGINT NOT NULL,
  PRIMARY KEY(envelope_id)
);