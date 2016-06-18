package query_builder

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestQueryBuilder(t *testing.T) {
	Convey("Can use query builder for one name", t, func() {
		setup()

		qb := NewSqlQueryBuilder()
		qb.WithName("Joe")

		res := []SqlQueryBuilderRecord{}
		qb.Apply(func(_e interface{}) {
			e := _e.(SqlQueryBuilderRecord)
			res = append(res, e)
		})

		So(res[0].Name, ShouldEqual, "Joe")
		So(len(res), ShouldEqual, 1)
	})

	Convey("No match query", t, func() {
		setup()

		qb := NewSqlQueryBuilder()
		qb.WithAge(33)

		res := []SqlQueryBuilderRecord{}
		qb.Apply(func(_e interface{}) {
			e := _e.(SqlQueryBuilderRecord)
			res = append(res, e)
		})

		So(len(res), ShouldEqual, 0)
	})
}

var _sql *gorm.DB

type SqlQueryBuilder struct {
	QueryBuilder

	NameAgeMixin
}

func NewSqlQueryBuilder() *SqlQueryBuilder {
	qb := &SqlQueryBuilder{}

	// Initialize our embedded `QueryBuilder` object.  The first parameter is the
	// runer.
	qb.Init(sqlQueryBuilderRunner, &qb.NameAgeMixin.QueryMixin)

	return qb
}

type SqlQueryBuilderRecord struct {
	Name string
	Age  int
}

// Define our 'runner' which will actually interpret and execute the query.
func standardQueryTransformer(queries []interface{}) *gorm.DB {
	sql := _sql

	// Because the hypothetical sql interface is also fluent
	// not much has to be done here.
	for _, _q := range queries {
		switch q := _q.(type) {
		case WithNameQuery:
			sql = sql.Where("name = ?", q.Name)
		case WithAgeQuery:
			sql = sql.Where("age = ?", q.Age)
		default:
			panic(fmt.Errorf("Unknown query type '%T'", q))
		}
	}

	return sql
}

var sqlQueryBuilderRunner = func(queries []interface{}) ([]interface{}, error) {
	sql := standardQueryTransformer(queries)

	_res := []SqlQueryBuilderRecord{}
	err := sql.Table("users").Find(&_res).Error
	if err != nil {
		return nil, err
	}

	var res []interface{}
	for _, e := range _res {
		res = append(res, e)
	}

	return res, nil
}

func setup() {
	// Setup SQL
	var err error
	_sql, err = gorm.Open("sqlite3", TempNonExistantFilePath())
	if err != nil {
		panic(err)
	}

	// Create our users table
	_sql.Exec(`
    CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    age INTEGER
    );

    INSERT INTO users (name, age) VALUES ("Joe", 21);
    INSERT INTO users (name, age) VALUES ("Bob", 34);
    INSERT INTO users (name, age) VALUES ("Ron", 45);
	`)
}

/*
	----------------------------------------------------------------------
	NameAgeMixin - Adds the 'WithName' and 'WithAge' query helpers
	----------------------------------------------------------------------
*/
type NameAgeMixin struct {
	QueryMixin
}

// The object that stores the query information
type WithNameQuery struct {
	Name string
}

type WithAgeQuery struct {
	Age int
}

// The query we've added
func (t *NameAgeMixin) WithName(name string) {
	t.Push(WithNameQuery{Name: name})
}

func (t *NameAgeMixin) WithAge(age int) {
	t.Push(WithAgeQuery{Age: age})
}
