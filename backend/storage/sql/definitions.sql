-- users table.
DROP TABLE IF EXISTS users;

CREATE TABLE users
  (user_id bigserial primary key,
   user_name varchar(256) default 'user');

-- sms table.
DROP TABLE IF EXISTS sms;

CREATE TABLE sms
  (sms_id bigserial primary key, 
   user_id bigint references users(user_id) on delete cascade,
   seen boolean default 'false',
   msg_time timestamp with time zone,
   origin varchar(128) default '',
   body varchar(256) default '');

-- Create an index on user_id for sms table.
CREATE INDEX IF NOT EXISTS sms_user_index ON sms (user_id);