package moodle

import (
	"database/sql"
	"fmt"
	"log"
	"moodle/internal/core/domain"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MariaDB struct {
	db *sql.DB
}

// new maria adapter
func NewMariaDB(db *sql.DB) *MariaDB {
	return &MariaDB{db: db}
}

// select data from mysql
func (m *MariaDB) Select(dbname, query string) map[string]string {
	queryString := fmt.Sprintf("%s", query)
	rows, err := m.db.Query(queryString)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Println(err)
		return nil
	}
	// map column and column data as key and value
	columnsMap := make(map[string]string)
	// make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// see http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println(err)
			return nil
		}
		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			columnsMap[columns[i]] = value
		}
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return columnsMap
}


// create new database from existing Database
func (m *MariaDB) CreateDB(dbname string) error {
	queryString := fmt.Sprintf("CREATE DATABASE %s", dbname)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}


// drop database 
func (m *MariaDB) DropDB(dbname string) error {
	queryString := fmt.Sprintf("DROP DATABASE %s", dbname)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("drop database success")
	return err
}

// restore database from sql filepath
func (m *MariaDB) RestoreDB(dbname, filepath string) error {
	comand :=  "mysql -u root -h 45.124.94.112 -p0YU8381WUlk1u9ysVbF4Qb5FigNW8z8uCvPI " + dbname + " < " + filepath
	// print current time
	log.Println("restoring databse", dbname, "at:", time.Now())
	output, err := exec.Command("bash", "-c", comand).Output()
	log.Println("restore database", dbname, "done at:", time.Now())
	log.Println("backing up mdl_course of ", dbname, "at:", time.Now())
	m.CloneCourseTable(dbname, "mdl_course_backup")
	log.Println("backup mdl_course of ", dbname, "done at:", time.Now())
	log.Println("backing up mdl_course_categories of ", dbname, "at:", time.Now())
	m.CloneCourseCategoriesTable(dbname, "mdl_course_categories_backup")
	log.Println("backup mdl_course_categories of ", dbname, "done at:", time.Now())

	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(output))
	return err
}


// update shortname, fullname on mdl_course table
func (m *MariaDB) UpdateDB(dbname, shortname, fullname string) error {
	queryString := fmt.Sprintf("UPDATE %s.mdl_course SET shortname = '%s', fullname = '%s' WHERE id = 1", dbname, shortname, fullname)

	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

// clone table from existing database
func (m *MariaDB) CloneCourseTable(dbname, tablename string) error {
	queryString := fmt.Sprintf("create table %s.%s as select * from %s.mdl_course;", dbname, tablename, dbname)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("clone course table success")
	return err
}

func (m *MariaDB) CloneCourseCategoriesTable(dbname, tablename string) error {
	queryString := fmt.Sprintf("create table %s.%s as select * from %s.mdl_course_categories;", dbname, tablename, dbname)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("clone course_categories table success")
	return err
}

// delete row from table
func (m *MariaDB) DeleteRow(dbname, tablename string, id int) error {
	queryString := fmt.Sprintf("DELETE FROM %s.%s WHERE id = %d", dbname, tablename, id)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("delete row success")
	return err
}

// update password from table mdl_user
func (m *MariaDB) UpdatePassword(dbname, password string) error {
	queryString := fmt.Sprintf("UPDATE %s.mdl_user SET password = '%s' WHERE id = 1357", dbname, password)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("update password success")
	return err
}

// delete course from table mdl_course
func (m *MariaDB) DeleteCourse(dbname string, id int) error {
	queryString := fmt.Sprintf("DELETE FROM %s.mdl_course WHERE id = %d", dbname, id)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	// log delete course id from database
	log.Println("delete course: ", id, "from: ", dbname, "success")
	return err
}

// delete category from table mdl_course_categories
func (m *MariaDB) DeleteCategory(dbname string, cateCourse domain.CateCourse) error {

	// delete all course in category
	for _, v := range cateCourse.Courses {
		queryString := fmt.Sprintf("DELETE FROM %s.mdl_course WHERE id = %d", dbname, v)
		_, err := m.db.Exec(queryString)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("delete course: ", v, "success")
	}

	// delete all sub category 
	for _, v := range cateCourse.SubCate{
		queryString := fmt.Sprintf("DELETE FROM %s.mdl_course_categories WHERE id = %d", dbname, v)
		_, err := m.db.Exec(queryString)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("delete sub category: ", v, "success")
	}
	
	// delete all course enroll in sub category
	for _, v := range cateCourse.SubCate{
		queryString := fmt.Sprintf("DELETE FROM %s.mdl_enrol WHERE courseid = %d", dbname, v)
		_, err := m.db.Exec(queryString)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("delete course enroll in sub category: ", v, "success")
	}

	// delete grade book
	for _, v := range cateCourse.SubCate{
		queryString := fmt.Sprintf("DELETE FROM %s.mdl_grade_items WHERE courseid = %d", dbname, v)
		_, err := m.db.Exec(queryString)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("delete grade book in sub category: ", v, "success")
	}

	// delete categories
	queryString := fmt.Sprintf("DELETE FROM %s.mdl_course_categories WHERE id = %d", dbname, cateCourse.CategoryId)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("delete category: ", cateCourse.CategoryId,  "success")
	// // chage all course in category to uncategorized
	// for _, v := range cateCourse.Courses {
	// 	queryString := fmt.Sprintf("UPDATE %s.mdl_course SET category = 0 WHERE id = %d", dbname, v)
	// 	_, err := m.db.Exec(queryString)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	log.Println("change course: ", v, "to uncategorized success")
	// }

	// delete all course in category
	// for _, v := range cateCourse.Courses {
	// 	m.DeleteCourse(dbname, v)
	// }
	// // delete from mdl_course_sections
	// for _, v := range cateCourse.Courses {
	// 	queryString := fmt.Sprintf("DELETE FROM %s.mdl_course_sections WHERE course = %d", dbname, v)
	// 	_, err := m.db.Exec(queryString)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	log.Println("delete course section: ", v, "success")
	// }
	// // delete from mdl_course_modules
	// for _, v := range cateCourse.Courses {
	// 	queryString := fmt.Sprintf("DELETE FROM %s.mdl_course_modules WHERE course = %d", dbname, v)
	// 	_, err := m.db.Exec(queryString)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	log.Println("delete course module: ", v, "success")
	// }
	// // delete from mdl_course_modules_completion
	// for _, v := range cateCourse.Courses {
	// 	queryString := fmt.Sprintf("DELETE FROM %s.mdl_course_modules_completion WHERE coursemoduleid IN (SELECT id FROM %s.mdl_course_modules WHERE course = %d)", dbname, dbname, v)
	// 	_, err := m.db.Exec(queryString)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	log.Println("delete course module completion: ", v, "success")
	// }
	// // delete from mdl_resource
	// for _, v := range cateCourse.Courses {
	// 	queryString := fmt.Sprintf("DELETE FROM %s.mdl_resource WHERE course = %d", dbname, v)
	// 	_, err := m.db.Exec(queryString)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	log.Println("delete resource: ", v, "success")
	// }
	// // delete from mdl_assign
	// for _, v := range cateCourse.Courses {
	// 	queryString := fmt.Sprintf("DELETE FROM %s.mdl_assign WHERE course = %d", dbname, v)
	// 	_, err := m.db.Exec(queryString)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	log.Println("delete assign: ", v, "success")
	// }
	return err
}

// copy row from table mdl_course_backup to mdl_course
func (m *MariaDB) CopyCourseRow(dbname string, id []int) error {
	for _, v := range id {
		queryString := fmt.Sprintf("INSERT INTO %s.%s SELECT * FROM %s.mdl_course_backup WHERE id = %d", dbname, "mdl_course", dbname, v)
		_, err := m.db.Exec(queryString)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	log.Println("copy course success")
	return nil
}

// copy row from table mdl_course_categories_backup to mdl_course_categories
func (m *MariaDB) CopyCategoryRow(dbname string, cateCourse domain.CateCourse) error {
	queryString := fmt.Sprintf("INSERT INTO %s.%s SELECT * FROM %s.mdl_course_categories_backup WHERE id = %d", dbname, "mdl_course_categories", dbname, cateCourse.CategoryId)
	_, err := m.db.Exec(queryString)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("copy cate success")
	// copy sub category
	for _, v := range cateCourse.SubCate {
		queryString := fmt.Sprintf("INSERT INTO %s.%s SELECT * FROM %s.mdl_course_categories_backup WHERE id = %d", dbname, "mdl_course_categories", dbname, v)
		_, err := m.db.Exec(queryString)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("copy sub course cate success")
	}

	// update category for mdl_course
	for _, v := range cateCourse.Courses {
		queryString := fmt.Sprintf("UPDATE %s.mdl_course SET category = %d WHERE id = %d", dbname, cateCourse.CategoryId, v)
		_, err := m.db.Exec(queryString)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("update course category success")
	}

	// // copy course
	// for _, v := range cateCourse.Courses {
	// 	queryString := fmt.Sprintf("INSERT INTO %s.%s SELECT * FROM %s.mdl_course_backup WHERE id = %d", dbname, "mdl_course", dbname, v)
	// 	_, err := m.db.Exec(queryString)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// 	log.Println("copy course success")
	// }
	return nil
}
