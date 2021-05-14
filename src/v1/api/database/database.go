package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)


// https://github.com/jmoiron/sqlx]'
var dBInstance *sqlx.DB
var DBInstancePublic *sqlx.DB

type ClassDatabase struct {
	Instance *sqlx.DB
}



type UserListDisplay struct {
	ID           uint   `json:"id" db:"id"`
	FirstName    string `json:"firstname" db:"firstname"`
	LastName     string `json:"lastname" db:"lastname"`
	Email        string `json:"email" db:"email"`
	MobileNumber string `json:"mobile_number" db:"mobile_number"`
	Role         string `json:"role" db:"role"`
	BankType     string `json:"bank_type" db:"bank_type"`
	BankNo       string `json:"bank_no" db:"bank_no"`
	Address      string `json:"address" db:"address"`
	Birthday     string `json:"birthday" db:"birthday"`
	Gender       string `json:"gender" db:"gender"`
	M88Account   string `json:"m88_account" db:"m88_account"`
}


func GetPaginationDetails(
	sql string,
	count int,
	page int,
	rowLimit int,
	pageLimit int,
) (string, []int, int, int, int, int, int) {

	var pages []int

	pageStart := 0
	pageEnd := pageLimit
	totalPages := int(math.Ceil(float64(count) / float64(rowLimit)))

	if page > totalPages || totalPages == 1 {
		page = 0
	}

	var rowsPerPage int

	if rowLimit > count {
		rowsPerPage = count
	} else {
		rowsPerPage = rowLimit
	}

	var offset int

	if count != 0 {
		if page == 0 {
			offset = 0
		} else if totalPages == 0 {
			offset = 0
		} else {
			if page >= totalPages {
				offset = 0
			} else {
				offset = page * rowLimit
			}
		}
	} else {
		offset = 0
	}

	for !(page >= 0 && page <= pageEnd) {
		pageStart = pageStart + pageLimit
		pageEnd = pageEnd + pageLimit
	}

	for i := pageStart; i <= pageEnd; i++ {
		if i <= totalPages - 1 {
			pages = append(pages, i)
		} else {
			if len(pages) > 0 {
				previousPage := pages[0] - 1

				if previousPage > 1 {
					pages = append([]int{previousPage}, pages...)
				}
			}
		}
	}

	sql = sql + " LIMIT replace_limit OFFSET replace_offset"

	sql = strings.Replace(sql, "replace_limit", strconv.Itoa(rowsPerPage), -1)
	sql = strings.Replace(sql, "replace_offset", strconv.Itoa(offset), -1)

	if len(pages) == 0 {
		pages = append(pages, 0)
	}

	// Handle result count
	resultCount := 0

	// Check if it has next page
	hasNextPage := false

	for _, ele := range pages {
		if ele > page {
			hasNextPage = true
			break
		}
	}

	if hasNextPage == true {
		resultCount = rowLimit
	} else {
		resultCount = count - offset
	}

	return sql, pages, rowsPerPage, offset, page, count, resultCount
}

func GetQueryCount(
	sql string,
	args ...interface{},
) (int, error) {

	var (
		count int
	)

	sql = "SELECT COUNT(*) FROM (" + sql + ") AS a"

	err := DBInstancePublic.Get(&count, sql, args...)

	if err != nil {
		return count, err
	}

	return count, nil
}

// Sets the global variable db instance
func EstablishConnection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbPort := os.Getenv("DB_PORT")

	fmt.Println(
		dbHost,
		dbUser,
		dbPassword,
		dbDatabase,
	)

	// Connect to MYSQL and execute queries
	connString := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":"+ dbPort +")/" + dbDatabase
	db, err := sqlx.Open("mysql", connString)

	if err != nil {
		fmt.Println("Error establishing database connection")
		panic(err.Error())
	}

	maxConnections, _ := strconv.Atoi(os.Getenv("DB_DATABASE"))

	db.SetMaxOpenConns(maxConnections)

	dBInstance = db
	DBInstancePublic = db

	testConnection()
}

// Performs a simple query to see if the connection succeeded
func testConnection() {
	_, err := dBInstance.Query("SELECT 5 AS test")

	if err != nil {
		fmt.Println("Error establishing database connection")
		panic(err.Error())
	}
}

func translate(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	translateRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func translateRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		translateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			translateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			translateRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		// translatedString := dict[original.Interface().(string)]
		// copy.SetString(translatedString)

	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

}

func GetLastInsertID(t *sqlx.Tx) (int, error) {
	var lastInsertId int

	sql := `SELECT LAST_INSERT_ID()`
	err :=  t.Get(&lastInsertId, sql)

	return lastInsertId, err
}