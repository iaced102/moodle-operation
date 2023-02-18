# API
```
go run main.go
```

# CLI
```
go run main.go COMMAND
```

# filedir 
## logo: 
### path : 2b/83/2b831e652219e350659e4a71af9c4b9c4c99411c
### url: http://123.31.39.253/pluginfile.php/1/theme_edumy/headerlogo1/1676349798/logo.png

## favicon: 
### path : /9b/f7/9bf7dca13c73a844c338cd16d1253c82b4bebfe8
### url: http://123.31.39.253/pluginfile.php/1/theme_edumy/favicon/1676349393/favicon.ico



# adapt db before restore from backup files
```
mysqldump -u root -p moodle > moodle01112022.sql
mysql -u duy -p5Yk7741J2JVWPTQkT9eKkcbAaTUs5XzTvIFL moodle < moodle01112022
sed -i 's/utf8mb4_0900_ai_ci/utf8_unicode_ci/g' moodle.sql
sed -i 's/utf8mb4/utf8/g' moodle.sql
sed -i 's/utf8_unicode_520_ci/utf8_unicode_ci/g' moodle.sql
```

