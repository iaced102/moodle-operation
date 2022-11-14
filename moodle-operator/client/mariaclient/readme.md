# adapt db before restore from backup files
```
mysqldump -u root -p moodle > moodle01112022.sql
mysql -u duy -p5Yk7741J2JVWPTQkT9eKkcbAaTUs5XzTvIFL moodle < moodle01112022
sed -i 's/utf8mb4_0900_ai_ci/utf8_unicode_ci/g' moodle.sql
sed -i 's/utf8mb4/utf8/g' moodle.sql
sed -i 's/utf8_unicode_520_ci/utf8_unicode_ci/g' moodle.sql
```
