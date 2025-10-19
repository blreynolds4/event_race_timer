delete from athlete_race;
delete from athlete;
delete from race;
delete from meet;
commit;
select *
from meet;
select *
from race;
select count(*)
from athlete_race;
select *
from athlete_race;
select count(*)
from athlete;
select *
from athlete;
-- for overall results
SELECT ar.bib,
  ar.place,
  ar.finish_time,
  ar.xc_place,
  a.id,
  a.da_id,
  a.first_name,
  a.last_name,
  a.team,
  a.grade,
  a.gender,
  r.name
FROM athlete a
  JOIN athlete_race ar ON a.id = ar.athlete_id
  inner join race r on ar.race_id = r.id
  inner join meet m on r.meet_id = m.id
WHERE m.id = 104 -- and r.id = 156
  and bib in (78, 88)
ORDER BY ar.finish_time;
UPDATE athlete_race
set finish_time = null,
  place = null,
  xc_place = null,
  finish_source = null,
  place_source = null;