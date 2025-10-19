-- Create the meet table for a meet with races
CREATE TABLE meet (
  id SERIAL PRIMARY KEY,
  name varchar(255) NOT NULL
);
-- Create a unique index on the meet name
create unique index idx_meet_name on meet(name);
--
-- Create the race table for a race within a meet
create table race (
  id SERIAL PRIMARY KEY,
  meet_id INTEGER NOT NULL,
  name varchar(255) NOT NULL,
  FOREIGN KEY (meet_id) REFERENCES meet(id)
);
-- Create a unique index on the race name
create unique index idx_meet_race_name on race(meet_id, name);