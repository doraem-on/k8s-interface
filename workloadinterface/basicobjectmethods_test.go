package workloadinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// IsBaseObject only checks that kind, apiVersion and metadata.name are
// non-empty, never that they are strings, so a base object carrying non-string
// values is accepted the same way a workload is.
func TestBaseObjectAccessorsWithNonStringFields(t *testing.T) {
	b, err := NewBaseObjBytes([]byte(`{
		"apiVersion": ["oops"],
		"kind": 1,
		"metadata": {"name": 42, "namespace": ["oops"]}
	}`))
	assert.NoError(t, err)

	for _, tc := range []struct {
		name string
		got  func(*BaseObject) string
	}{
		{"GetNamespace", (*BaseObject).GetNamespace},
		{"GetName", (*BaseObject).GetName},
		{"GetApiVersion", (*BaseObject).GetApiVersion},
		{"GetKind", (*BaseObject).GetKind},
		{"GetGroup", (*BaseObject).GetGroup},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, "", tc.got(b))
		})
	}

	assert.NotPanics(t, func() { _ = b.GetID() })
}

func TestListWorkloadsAccessorsWithNonStringFields(t *testing.T) {
	lw, err := NewListWorkloads([]byte(`{"apiVersion": ["oops"], "kind": 1, "metadata": {"name": 42}}`))
	assert.NoError(t, err)

	assert.Equal(t, "", lw.GetName())
	assert.Equal(t, "", lw.GetApiVersion())
	assert.Equal(t, "", lw.GetKind())
	assert.NotPanics(t, func() { _ = lw.GetID() })
}
