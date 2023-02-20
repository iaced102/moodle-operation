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
### path : /2b/83/2b831e652219e350659e4a71af9c4b9c4c99411c
### url: http://hello1.lms.bizflycloud.vn/pluginfile.php/1/theme_edumy/headerlogo1/1676573021/Logo mới Bizfly Cloud-01.png

## favicon: 
### path : /35/1d/351d302a3d094fdde4638e7dbdde86af3199be7e
### url: http://hello.lms.bizflycloud.vn/pluginfile.php/1/theme_edumy/favicon/1676573021/z3665638475480_d6dab64f97b26f1c73cd539411ae9990.jpg



# adapt db before restore from backup files
```
mysqldump -u root -p moodle > moodle01112022.sql
mysql -u duy -p5Yk7741J2JVWPTQkT9eKkcbAaTUs5XzTvIFL moodle < moodle01112022
sed -i 's/utf8mb4_0900_ai_ci/utf8_unicode_ci/g' moodle.sql
sed -i 's/utf8mb4/utf8/g' moodle.sql
sed -i 's/utf8_unicode_520_ci/utf8_unicode_ci/g' moodle.sql
```

