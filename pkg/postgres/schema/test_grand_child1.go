// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableTestGrandChild1Stmt holds the create statement for table `test_grand_child1`.
	CreateTableTestGrandChild1Stmt = &postgres.CreateStmts{
		GormModel: (*TestGrandChild1)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// TestGrandChild1Schema is the go schema for table `test_grand_child1`.
	TestGrandChild1Schema = func() *walker.Schema {
		schema := GetSchemaForTable("test_grand_child1")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestGrandChild1)(nil)), "test_grand_child1")
		referencedSchemas := map[string]*walker.Schema{
			"storage.TestChild1":       TestChild1Schema,
			"storage.TestGGrandChild1": TestGGrandChild1Schema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(64), "testgrandchild1", (*storage.TestGrandChild1)(nil)))
		RegisterTable(schema, CreateTableTestGrandChild1Stmt)
		return schema
	}()
)

const (
	TestGrandChild1TableName = "test_grand_child1"
)

// TestGrandChild1 holds the Gorm model for Postgres table `test_grand_child1`.
type TestGrandChild1 struct {
	Id            string     `gorm:"column:id;type:varchar;primaryKey"`
	ParentId      string     `gorm:"column:parentid;type:varchar"`
	ChildId       string     `gorm:"column:childid;type:varchar"`
	Val           string     `gorm:"column:val;type:varchar"`
	Serialized    []byte     `gorm:"column:serialized;type:bytea"`
	TestChild1Ref TestChild1 `gorm:"foreignKey:parentid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
