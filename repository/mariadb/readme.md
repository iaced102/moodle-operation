# adapt db before restore from backup files
```
sed -i 's/utf8mb4_0900_ai_ci/utf8_unicode_ci/g' moodle.sql
sed -i 's/utf8mb4/utf8/g' moodle.sql
sed -i 's/utf8_unicode_520_ci/utf8_unicode_ci/g' moodle.sql
```
