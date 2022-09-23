// Code generated by pg-bindings generator. DO NOT EDIT.
package schema

import (
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

// ConvertTestGrandparentFromProto converts a `*storage.TestGrandparent` to Gorm model
func ConvertTestGrandparentFromProto(obj *storage.TestGrandparent) (*schema.TestGrandparents, error) {
	serialized, err := obj.Marshal()
	if err != nil {
		return nil, err
	}
	model := &schema.TestGrandparents{
		Id:         obj.GetId(),
		Val:        obj.GetVal(),
		Priority:   obj.GetPriority(),
		RiskScore:  obj.GetRiskScore(),
		Serialized: serialized,
	}
	return model, nil
}

// ConvertTestGrandparent_EmbeddedFromProto converts a `*storage.TestGrandparent_Embedded` to Gorm model
func ConvertTestGrandparent_EmbeddedFromProto(obj *storage.TestGrandparent_Embedded, idx int, test_grandparents_Id string) (*schema.TestGrandparentsEmbeddeds, error) {
	model := &schema.TestGrandparentsEmbeddeds{
		TestGrandparentsId: test_grandparents_Id,
		Idx:                idx,
		Val:                obj.GetVal(),
	}
	return model, nil
}

// ConvertTestGrandparent_Embedded_Embedded2FromProto converts a `*storage.TestGrandparent_Embedded_Embedded2` to Gorm model
func ConvertTestGrandparent_Embedded_Embedded2FromProto(obj *storage.TestGrandparent_Embedded_Embedded2, idx int, test_grandparents_Id string, test_grandparents_embeddeds_idx int) (*schema.TestGrandparentsEmbeddedsEmbedded2, error) {
	model := &schema.TestGrandparentsEmbeddedsEmbedded2{
		TestGrandparentsId:           test_grandparents_Id,
		TestGrandparentsEmbeddedsIdx: test_grandparents_embeddeds_idx,
		Idx:                          idx,
		Val:                          obj.GetVal(),
	}
	return model, nil
}

// ConvertTestGrandparentToProto converts Gorm model `TestGrandparents` to its protobuf type object
func ConvertTestGrandparentToProto(m *schema.TestGrandparents) (*storage.TestGrandparent, error) {
	var msg storage.TestGrandparent
	if err := msg.Unmarshal(m.Serialized); err != nil {
		return nil, err
	}
	return &msg, nil
}

func TestTestGrandparentSerialization(t *testing.T) {
	obj := &storage.TestGrandparent{}
	assert.NoError(t, testutils.FullInit(obj, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	m, err := ConvertTestGrandparentFromProto(obj)
	assert.NoError(t, err)
	conv, err := ConvertTestGrandparentToProto(m)
	assert.NoError(t, err)
	assert.Equal(t, obj, conv)
}