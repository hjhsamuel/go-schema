package go_schema

import (
	"reflect"
	"testing"
)

func mustLoadSchemaTemplate(t *testing.T, tpl map[string]any) *Schema {
	t.Helper()
	s := &Schema{}
	if err := s.LoadTemplate(tpl); err != nil {
		t.Fatalf("LoadTemplate failed: %v", err)
	}
	return s
}

func TestSchemaLoadTemplateAndDumpTemplate(t *testing.T) {
	tpl := map[string]any{
		"id":   "moduleA",
		"name": "Module A",
		"fields": []any{
			map[string]any{"id": "f1", "name": "f1", "kind": "string"},
			map[string]any{
				"id":     "f2",
				"name":   "f2",
				"kind":   "select",
				"select": []any{"a", "b"},
			},
		},
	}

	s := mustLoadSchemaTemplate(t, tpl)
	if s.ID != "moduleA" || s.Name != "Module A" || len(s.Fields) != 2 {
		t.Fatalf("unexpected schema loaded: %+v", s)
	}

	out := s.DumpTemplate()
	if out["id"] != "moduleA" || out["name"] != "Module A" {
		t.Fatalf("dump mismatch: %#v", out)
	}
	fields, ok := out["fields"].([]map[string]any)
	if !ok || len(fields) != 2 {
		t.Fatalf("dumped fields mismatch: %#v", out["fields"])
	}
}

func TestSchemaLoadTemplateErrors(t *testing.T) {
	tests := []struct {
		name string
		tpl  map[string]any
		err  string
	}{
		{
			name: "missing id",
			tpl:  map[string]any{"name": "n"},
			err:  "id is required",
		},
		{
			name: "missing name",
			tpl:  map[string]any{"id": "m"},
			err:  "name is required",
		},
		{
			name: "field elem is not map",
			tpl: map[string]any{
				"id":     "m",
				"name":   "n",
				"fields": []any{"bad"},
			},
			err: "field elem is not a map",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Schema{}
			err := s.LoadTemplate(tt.tpl)
			assertErrContains(t, err, tt.err)
		})
	}
}

func TestSchemaLoadValueAndDumpValue(t *testing.T) {
	tpl := map[string]any{
		"id":   "moduleA",
		"name": "Module A",
		"fields": []any{
			map[string]any{
				"id":       "f1",
				"name":     "f1",
				"kind":     "string",
				"required": true,
			},
			map[string]any{
				"id":     "f2",
				"name":   "f2",
				"kind":   "number",
				"number": nil,
			},
		},
	}
	base := mustLoadSchemaTemplate(t, tpl)

	loaded, err := base.LoadValue(map[string]any{
		"moduleA": map[string]any{
			"f1": "abc",
			"f2": 12.5,
		},
	})
	if err != nil {
		t.Fatalf("LoadValue failed: %v", err)
	}

	out, err := loaded.DumpValue()
	if err != nil {
		t.Fatalf("DumpValue failed: %v", err)
	}
	want := map[string]any{
		"moduleA": map[string]any{
			"f1": "abc",
			"f2": 12.5,
		},
	}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("dump mismatch, want=%#v got=%#v", want, out)
	}
}

func TestSchemaLoadValueErrors(t *testing.T) {
	base := mustLoadSchemaTemplate(t, map[string]any{
		"id":   "moduleA",
		"name": "Module A",
		"fields": []any{
			map[string]any{
				"id":       "f1",
				"name":     "f1",
				"kind":     "string",
				"required": true,
			},
		},
	})

	_, err := base.LoadValue(map[string]any{})
	assertErrContains(t, err, "module moduleA not found")

	_, err = base.LoadValue(map[string]any{"moduleA": "bad"})
	assertErrContains(t, err, "field value is not a map")

	_, err = base.LoadValue(map[string]any{"moduleA": map[string]any{}})
	assertErrContains(t, err, "field f1 is required")
}

func TestSchemaDumpValueError(t *testing.T) {
	s := &Schema{
		ID:   "m",
		Name: "m",
		Fields: []*Field{
			{ID: "s1", Name: "s1", Kind: Select},
		},
	}
	_, err := s.DumpValue()
	assertErrContains(t, err, "select: field value is empty")
}
