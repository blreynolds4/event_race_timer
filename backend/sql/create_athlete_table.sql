CREATE TABLE athlete (
  id SERIAL PRIMARY KEY,
  da_id varchar(255) NOT NULL,
  first_name varchar(255) NOT NULL,
  last_name varchar(255) NOT NULL,
  gender varchar(1) NOT NULL,
  grade integer not null,
  team varchar(255) NOT NULL
);
create unique index idx_athlete_da_id on athlete(da_id);
-- Create athlete_race table to associate athletes with races
create table athlete_race (
  athlete_id INTEGER NOT NULL,
  race_id INTEGER NOT NULL,
  bib integer NOT NULL,
  finish_time integer DEFAULT null,
  place integer DEFAULT null,
  xc_place integer DEFAULT null,
  finish_source varchar(50) DEFAULT null,
  place_source varchar(50) DEFAULT null,
  FOREIGN KEY (athlete_id) REFERENCES athlete(id),
  FOREIGN KEY (race_id) REFERENCES race(id)
);
-- one athlete in one race at a time
create unique index idx_athlete_race on athlete_race(athlete_id, race_id, bib);