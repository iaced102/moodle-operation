delete from mdl_course where id= 9;
delete from mdl_course where id=35;

delete from mdl_course where id=32;
delete from mdl_course where id=33;
delete from mdl_course where id=34;
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

delete from mdl_course_categories where id=30;
delete from mdl_course_categories where id=31;
delete from mdl_course_categories where id=32;
delete from mdl_course_categories where id=34;
delete from mdl_course_categories where id=35;
delete from mdl_course_categories where id=36;

delete from mdl_course_categories where id=5;
delete from mdl_course_categories where id=27;
delete from mdl_course_categories where id=28;
delete from mdl_course_categories where id=29;
delete from mdl_course_categories where id=33;
delete from mdl_course_categories where id=37;
delete from mdl_course_categories where id=38;
delete from mdl_course_categories where id=39;

# clone table
create table c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup as select * from c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course;
create table c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup as select * from c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories;
create table mdl_course as select * from mdl_course_backup;
create table mdl_course_categories as select * from mdl_course_categories_backup;

drop table mdl_course;
drop table mdl_course_categories;


-- update visible field in mdl_course_categories where id=29;
UPDATE mdl_course_categories SET visible = 0 WHERE id = 29;

-- check last 5 minutes online users
SELECT * FROM mdl_user WHERE lastaccess > UNIX_TIMESTAMP() - 300;

-- copy row
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 9;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 35;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 36;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 37;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 38;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 39;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 40;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 41;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 42;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 43;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 44;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 45;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 46;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 47;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 48;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 49;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 50;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 51;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 52;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 53;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 54;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 55;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 56;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 57;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 58;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 59;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 60;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 61;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 62;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 63;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 64;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 65;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 66;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 67;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 68;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 67;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 68;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 69;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 70;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 71;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_backup WHERE id = 72;

INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 5;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 27;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 28;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 29;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 33;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 37;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 38;
INSERT INTO c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories SELECT * FROM c0278b53_aacc_480f_a31f_d4943d639c4a.mdl_course_categories_backup WHERE id = 39;



create table mdl_course as select * from mdl_course_backup;
create table mdl_course_categories as select * from mdl_course_categories_backup;

INSERT INTO mdl_course SELECT * FROM mdl_course_backup WHERE id = > 35;
INSERT INTO mdl_course_categories SELECT * FROM mdl_course_categories_backup WHERE id > 27;
