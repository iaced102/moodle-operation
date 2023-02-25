delete from mdl_course where id=32;
delete from mdl_course where id=33;
delete from mdl_course where id=34;
delete from mdl_course where id=35;
delete from mdl_course where id=36;
delete from mdl_course where id=37;
delete from mdl_course where id=38;
delete from mdl_course where id=39;
delete from mdl_course where id=40;
delete from mdl_course where id=41;
delete from mdl_course where id=42;
delete from mdl_course where id=43;
delete from mdl_course where id=44;
delete from mdl_course where id=45;
delete from mdl_course where id=46;
delete from mdl_course where id=47;
delete from mdl_course where id=48;
delete from mdl_course where id=49;
delete from mdl_course where id=50;
delete from mdl_course where id=51;
delete from mdl_course where id=52;
delete from mdl_course where id=53;
delete from mdl_course where id=54;
delete from mdl_course where id=55;
delete from mdl_course where id=56;
delete from mdl_course where id=57;
delete from mdl_course where id=58;
delete from mdl_course where id=59;
delete from mdl_course where id=60;
delete from mdl_course where id=61;
delete from mdl_course where id=62;
delete from mdl_course where id=63;
delete from mdl_course where id=64;
delete from mdl_course where id=65;
delete from mdl_course where id=66;
delete from mdl_course where id=67;
delete from mdl_course where id=68;
delete from mdl_course where id=69;
delete from mdl_course where id=70;
delete from mdl_course where id=71;
delete from mdl_course where id=72;


delete from mdl_course_categories where id=27;
delete from mdl_course_categories where id=28;
delete from mdl_course_categories where id=29;
delete from mdl_course_categories where id=30;
delete from mdl_course_categories where id=31;
delete from mdl_course_categories where id=32;
delete from mdl_course_categories where id=33;
delete from mdl_course_categories where id=34;
delete from mdl_course_categories where id=35;
delete from mdl_course_categories where id=36;
delete from mdl_course_categories where id=37;
delete from mdl_course_categories where id=38;
delete from mdl_course_categories where id=39;

# clone table
create table 8cb46d10_574b_4165_be54_2e9205a75b2e.mdl_course_backup as select * from 8cb46d10_574b_4165_be54_2e9205a75b2e.mdl_course;
create table 8cb46d10_574b_4165_be54_2e9205a75b2e.mdl_course_categories_backup as select * from 8cb46d10_574b_4165_be54_2e9205a75b2e.mdl_course_categories;
create table mdl_course as select * from mdl_course_deleted;
create table mdl_course_categories as select * from mdl_course_categories_deleted;

drop table mdl_course;
drop table mdl_course_categories;


-- update visible field in mdl_course_categories where id=29;
UPDATE mdl_course_categories SET visible = 0 WHERE id = 29;

-- check last 5 minutes online users
SELECT * FROM mdl_user WHERE lastaccess > UNIX_TIMESTAMP() - 300;

