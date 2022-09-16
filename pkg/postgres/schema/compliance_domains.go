// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableComplianceDomainsStmt holds the create statement for table `compliance_domains`.
	CreateTableComplianceDomainsStmt = &postgres.CreateStmts{
		GormModel: (*ComplianceDomains)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ComplianceDomainsSchema is the go schema for table `compliance_domains`.
	ComplianceDomainsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("compliance_domains")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ComplianceDomain)(nil)), "compliance_domains")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_COMPLIANCE_DOMAIN, "compliancedomain", (*storage.ComplianceDomain)(nil)))
		RegisterTable(schema, CreateTableComplianceDomainsStmt)
		return schema
	}()
)

const (
	ComplianceDomainsTableName = "compliance_domains"
)

// ComplianceDomains holds the Gorm model for Postgres table `compliance_domains`.
type ComplianceDomains struct {
	Id            string            `gorm:"column:id;type:varchar;primaryKey"`
	ClusterId     string            `gorm:"column:cluster_id;type:varchar"`
	ClusterName   string            `gorm:"column:cluster_name;type:varchar"`
	ClusterLabels map[string]string `gorm:"column:cluster_labels;type:jsonb"`
	Serialized    []byte            `gorm:"column:serialized;type:bytea"`
}
