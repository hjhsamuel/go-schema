package go_schema

import (
	"reflect"
	"strings"
	"testing"
)

func mustLoadFieldTemplate(t *testing.T, tpl map[string]any) *Field {
	t.Helper()
	f := &Field{}
	if err := f.LoadTemplate(tpl); err != nil {
		t.Fatalf("LoadTemplate failed: %v", err)
	}
	return f
}

func assertErrContains(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", msg)
	}
	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("expected error containing %q, got %q", msg, err.Error())
	}
}

func TestFieldLoadTemplateErrors(t *testing.T) {
	tests := []struct {
		name string
		tpl  map[string]any
		err  string
	}{
		{
			name: "missing id",
			tpl:  map[string]any{"name": "n", "kind": "string"},
			err:  "id is required",
		},
		{
			name: "missing name",
			tpl:  map[string]any{"id": "id", "kind": "string"},
			err:  "name is required",
		},
		{
			name: "missing kind",
			tpl:  map[string]any{"id": "id", "name": "n"},
			err:  "kind is required",
		},
		{
			name: "unsupported kind",
			tpl: map[string]any{
				"id":   "id",
				"name": "n",
				"kind": "unknown",
			},
			err: "unsupported field kind",
		},
		{
			name: "array kind not list",
			tpl: map[string]any{
				"id":    "id",
				"name":  "n",
				"kind":  "array",
				"array": "not-list",
			},
			err: "field kind array is not list",
		},
		{
			name: "select element not string",
			tpl: map[string]any{
				"id":     "id",
				"name":   "n",
				"kind":   "select",
				"select": []any{1.0},
			},
			err: "select element type is not string",
		},
		{
			name: "http method unsupported",
			tpl: map[string]any{
				"id":   "id",
				"name": "n",
				"kind": "httprequest",
				"httprequest": map[string]any{
					"url":    "https://example.com",
					"method": "PUT",
				},
			},
			err: "not support method",
		},
		{
			name: "http missing url",
			tpl: map[string]any{
				"id":   "id",
				"name": "n",
				"kind": "httprequest",
				"httprequest": map[string]any{
					"method": "GET",
				},
			},
			err: "has not subfield `url`",
		},
		{
			name: "http invalid url",
			tpl: map[string]any{
				"id":   "id",
				"name": "n",
				"kind": "httprequest",
				"httprequest": map[string]any{
					"url":    "://bad",
					"method": "GET",
				},
			},
			err: "url parse error",
		},
		{
			name: "nestselect kind not list",
			tpl: map[string]any{
				"id":         "id",
				"name":       "n",
				"kind":       "nestselect",
				"nestselect": "bad",
			},
			err: "field kind nestselect is not list",
		},
		{
			name: "selectarray kind not list",
			tpl: map[string]any{
				"id":          "id",
				"name":        "n",
				"kind":        "selectarray",
				"selectarray": "bad",
			},
			err: "field kind selectarray is not list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Field{}
			err := f.LoadTemplate(tt.tpl)
			assertErrContains(t, err, tt.err)
		})
	}
}

func TestFieldTemplateAndValueRoundTrip(t *testing.T) {
	roundTripCases := []struct {
		name      string
		template  map[string]any
		input     any
		wantValue any
	}{
		{
			name: "select",
			template: map[string]any{
				"id":       "f1",
				"name":     "f1",
				"kind":     "select",
				"required": true,
				"select":   []any{"a", "b"},
			},
			input:     "b",
			wantValue: "b",
		},
		{
			name: "selectarray",
			template: map[string]any{
				"id":          "f2",
				"name":        "f2",
				"kind":        "selectarray",
				"selectarray": []any{"x", "y", "z"},
			},
			input:     []any{"x", "z"},
			wantValue: []string{"x", "z"},
		},
		{
			name: "string",
			template: map[string]any{
				"id":        "f3",
				"name":      "f3",
				"kind":      "string",
				"validator": "required",
			},
			input:     "hello",
			wantValue: "hello",
		},
		{
			name: "text",
			template: map[string]any{
				"id":   "f4",
				"name": "f4",
				"kind": "text",
			},
			input:     "long text",
			wantValue: "long text",
		},
		{
			name: "number",
			template: map[string]any{
				"id":   "f5",
				"name": "f5",
				"kind": "number",
			},
			input:     3.14,
			wantValue: 3.14,
		},
		{
			name: "password",
			template: map[string]any{
				"id":   "f6",
				"name": "f6",
				"kind": "password",
			},
			input:     "secret",
			wantValue: "secret",
		},
		{
			name: "append",
			template: map[string]any{
				"id":   "f7",
				"name": "f7",
				"kind": "append",
			},
			input:     []any{"a", 2},
			wantValue: []string{"a", "2"},
		},
	}

	for _, tt := range roundTripCases {
		t.Run(tt.name, func(t *testing.T) {
			base := mustLoadFieldTemplate(t, tt.template)
			loaded, err := base.LoadValue(tt.input)
			if err != nil {
				t.Fatalf("LoadValue failed: %v", err)
			}

			got, err := loaded.DumpValue()
			if err != nil {
				t.Fatalf("DumpValue failed: %v", err)
			}
			if !reflect.DeepEqual(got, tt.wantValue) {
				t.Fatalf("DumpValue mismatch, want=%#v got=%#v", tt.wantValue, got)
			}
		})
	}
}

func TestFieldTemplateAndValueRoundTripCompositeKinds(t *testing.T) {
	arrayTpl := map[string]any{
		"id":   "arr",
		"name": "arr",
		"kind": "array",
		"array": []any{
			map[string]any{"id": "c1", "name": "c1", "kind": "string", "required": true},
			map[string]any{"id": "c2", "name": "c2", "kind": "number"},
		},
	}
	arrayField := mustLoadFieldTemplate(t, arrayTpl)
	arrayLoaded, err := arrayField.LoadValue(map[string]any{"c1": "v", "c2": 1.5})
	if err != nil {
		t.Fatalf("array LoadValue failed: %v", err)
	}
	arrayOut, err := arrayLoaded.DumpValue()
	if err != nil {
		t.Fatalf("array DumpValue failed: %v", err)
	}
	if !reflect.DeepEqual(arrayOut, map[string]any{"c1": "v", "c2": 1.5}) {
		t.Fatalf("array mismatch: %#v", arrayOut)
	}

	nestSelectTpl := map[string]any{
		"id":   "ns",
		"name": "ns",
		"kind": "nestselect",
		"nestselect": []any{
			map[string]any{"id": "s1", "name": "s1", "kind": "string"},
			map[string]any{"id": "s2", "name": "s2", "kind": "number"},
		},
	}
	nestSelect := mustLoadFieldTemplate(t, nestSelectTpl)
	nsLoaded, err := nestSelect.LoadValue(map[string]any{"s2": 8.0})
	if err != nil {
		t.Fatalf("nestselect LoadValue failed: %v", err)
	}
	nsOut, err := nsLoaded.DumpValue()
	if err != nil {
		t.Fatalf("nestselect DumpValue failed: %v", err)
	}
	if !reflect.DeepEqual(nsOut, map[string]any{"s2": 8.0}) {
		t.Fatalf("nestselect mismatch: %#v", nsOut)
	}

	nestArrayTpl := map[string]any{
		"id":   "nsa",
		"name": "nsa",
		"kind": "nestselectarray",
		"nestselectarray": []any{
			map[string]any{"id": "m1", "name": "m1", "kind": "string"},
			map[string]any{"id": "m2", "name": "m2", "kind": "number"},
		},
	}
	nestArray := mustLoadFieldTemplate(t, nestArrayTpl)
	nsaLoaded, err := nestArray.LoadValue(map[string]any{"m1": "a", "m2": 2.0})
	if err != nil {
		t.Fatalf("nestselectarray LoadValue failed: %v", err)
	}
	nsaOut, err := nsaLoaded.DumpValue()
	if err != nil {
		t.Fatalf("nestselectarray DumpValue failed: %v", err)
	}
	if !reflect.DeepEqual(nsaOut, map[string]any{"m1": "a", "m2": 2.0}) {
		t.Fatalf("nestselectarray mismatch: %#v", nsaOut)
	}

	incTpl := map[string]any{
		"id":   "inc",
		"name": "inc",
		"kind": "incrementarray",
		"incrementarray": []any{
			map[string]any{"id": "i1", "name": "i1", "kind": "string"},
			map[string]any{"id": "i2", "name": "i2", "kind": "number"},
		},
	}
	inc := mustLoadFieldTemplate(t, incTpl)
	incLoaded, err := inc.LoadValue([]any{
		map[string]any{"i1": "x", "i2": 1.0},
		map[string]any{"i1": "y", "i2": 2.0},
	})
	if err != nil {
		t.Fatalf("incrementarray LoadValue failed: %v", err)
	}
	incOut, err := incLoaded.DumpValue()
	if err != nil {
		t.Fatalf("incrementarray DumpValue failed: %v", err)
	}
	wantInc := []any{
		map[string]any{"i1": "x", "i2": 1.0},
		map[string]any{"i1": "y", "i2": 2.0},
	}
	if !reflect.DeepEqual(incOut, wantInc) {
		t.Fatalf("incrementarray mismatch: %#v", incOut)
	}
}

func TestFieldHttpRequestTemplateAndValue(t *testing.T) {
	tpl := map[string]any{
		"id":       "h1",
		"name":     "h1",
		"kind":     "httprequest",
		"required": true,
		"httprequest": map[string]any{
			"url":          "https://example.com/data",
			"method":       "post",
			"multi_select": false,
			"user_data":    "meta",
		},
	}
	base := mustLoadFieldTemplate(t, tpl)
	dump := base.DumpTemplate()
	httpPart, ok := dump["httprequest"].(map[string]any)
	if !ok {
		t.Fatalf("httprequest template missing in dump")
	}
	if httpPart["method"] != "POST" {
		t.Fatalf("expected normalized method POST, got %v", httpPart["method"])
	}

	loaded, err := base.LoadValue([]any{"a", "b"})
	if err != nil {
		t.Fatalf("LoadValue failed: %v", err)
	}
	v, err := loaded.DumpValue()
	if err != nil {
		t.Fatalf("DumpValue failed: %v", err)
	}
	if !reflect.DeepEqual(v, []any{"a"}) {
		t.Fatalf("expected single selected value, got %#v", v)
	}

	_, err = base.LoadValue([]any{})
	assertErrContains(t, err, "field is required")
}

func TestFieldLoadValueAndDumpValueErrors(t *testing.T) {
	selectField := mustLoadFieldTemplate(t, map[string]any{
		"id":     "s",
		"name":   "s",
		"kind":   "select",
		"select": []any{"a"},
	})
	_, err := selectField.LoadValue(1)
	assertErrContains(t, err, "field value is not a string")
	_, err = selectField.LoadValue("x")
	assertErrContains(t, err, "select value 'x' invalid")

	nestSelect := mustLoadFieldTemplate(t, map[string]any{
		"id":   "ns",
		"name": "ns",
		"kind": "nestselect",
		"nestselect": []any{
			map[string]any{"id": "n1", "name": "n1", "kind": "string"},
		},
	})
	emptyNest := &Field{Kind: NestSelect}
	_, err = emptyNest.DumpValue()
	assertErrContains(t, err, "nestselect: field value is empty")
	_, err = nestSelect.LoadValue("bad")
	assertErrContains(t, err, "field value is not a map")

	emptyInc := &Field{Kind: IncrementArray}
	_, err = emptyInc.LoadValue([]any{})
	assertErrContains(t, err, "template field is empty")

	selectDumpErr := &Field{Kind: Select}
	_, err = selectDumpErr.DumpValue()
	assertErrContains(t, err, "field value is empty")

	badKind := &Field{Kind: FieldKind("bad-kind")}
	_, err = badKind.DumpValue()
	assertErrContains(t, err, "unsupported field kind")
}

func TestFieldDumpTemplateCompositeBranches(t *testing.T) {
	arrayField := &Field{
		ID:       "arr",
		Name:     "arr",
		Kind:     Array,
		required: true,
		desc:     "array-desc",
		array: []*Field{
			{ID: "c1", Name: "c1", Kind: String},
		},
	}
	arrayDump := arrayField.DumpTemplate()
	if arrayDump["required"] != true || arrayDump["desc"] != "array-desc" {
		t.Fatalf("array meta not dumped: %#v", arrayDump)
	}
	if _, ok := arrayDump["array"].([]any); !ok {
		t.Fatalf("array field list not dumped: %#v", arrayDump)
	}

	nestSelect := &Field{
		ID:   "ns",
		Name: "ns",
		Kind: NestSelect,
		nestSelectOne: []*Field{
			{ID: "n1", Name: "n1", Kind: String},
		},
	}
	if _, ok := nestSelect.DumpTemplate()["nestselect"].([]any); !ok {
		t.Fatalf("nestselect not dumped correctly")
	}

	nestArray := &Field{
		ID:   "na",
		Name: "na",
		Kind: NestSelectArray,
		nestSelectMultiple: []*Field{
			{ID: "n1", Name: "n1", Kind: String},
		},
	}
	if _, ok := nestArray.DumpTemplate()["nestselectarray"].([]any); !ok {
		t.Fatalf("nestselectarray not dumped correctly")
	}

	incEmpty := (&Field{ID: "inc", Name: "inc", Kind: IncrementArray}).DumpTemplate()
	if v, ok := incEmpty["incrementarray"].([]any); !ok || len(v) != 0 {
		t.Fatalf("empty incrementarray dump mismatch: %#v", incEmpty)
	}

	incNonEmpty := (&Field{
		ID:   "inc2",
		Name: "inc2",
		Kind: IncrementArray,
		incrementArray: [][]*Field{
			{
				{ID: "i1", Name: "i1", Kind: String},
			},
		},
	}).DumpTemplate()
	if v, ok := incNonEmpty["incrementarray"].([]any); !ok || len(v) != 1 {
		t.Fatalf("non-empty incrementarray dump mismatch: %#v", incNonEmpty)
	}

	strField := (&Field{
		ID:        "s",
		Name:      "s",
		Kind:      String,
		validator: "v",
	}).DumpTemplate()
	if strField["validator"] != "v" {
		t.Fatalf("validator not dumped: %#v", strField)
	}

	httpNil := (&Field{ID: "h1", Name: "h1", Kind: HttpRequest}).DumpTemplate()
	if _, ok := httpNil["httprequest"]; ok {
		t.Fatalf("nil httprequest should not be dumped: %#v", httpNil)
	}
}

func TestFieldLoadValueErrorPathsByKind(t *testing.T) {
	_, err := (&Field{Kind: Array}).LoadValue("bad")
	assertErrContains(t, err, "array: field value is not a map")

	_, err = (&Field{Kind: SelectArray, required: true, selectMultiple: []string{"a"}}).LoadValue([]any{})
	assertErrContains(t, err, "selectarray value invalid")
	_, err = (&Field{Kind: SelectArray, selectMultiple: []string{"a"}}).LoadValue("bad")
	assertErrContains(t, err, "field value is not list")

	_, err = (&Field{Kind: String}).LoadValue(1)
	assertErrContains(t, err, "string: field value is not a string")
	_, err = (&Field{Kind: Text}).LoadValue(1)
	assertErrContains(t, err, "text: field value is not a string")
	_, err = (&Field{Kind: Number}).LoadValue("1")
	assertErrContains(t, err, "number: field value is not a number")
	_, err = (&Field{Kind: Password}).LoadValue(1)
	assertErrContains(t, err, "password: field value is not a string")

	_, err = (&Field{Kind: NestSelectArray}).LoadValue("bad")
	assertErrContains(t, err, "nestselectarray: field value is not a map")

	_, err = (&Field{
		Kind: IncrementArray,
		incrementArray: [][]*Field{
			{{ID: "a", Name: "a", Kind: String}},
		},
	}).LoadValue("bad")
	assertErrContains(t, err, "incrementArray: field value is not a map")
	_, err = (&Field{
		Kind: IncrementArray,
		incrementArray: [][]*Field{
			{{ID: "a", Name: "a", Kind: String}},
		},
	}).LoadValue([]any{"bad"})
	assertErrContains(t, err, "elem value is not a map")
	_, err = (&Field{
		Kind: IncrementArray,
		incrementArray: [][]*Field{
			{{ID: "a", Name: "a", Kind: String}},
		},
	}).LoadValue([]any{map[string]any{}})
	assertErrContains(t, err, "value of field a not set")

	_, err = (&Field{Kind: Append}).LoadValue("bad")
	assertErrContains(t, err, "append: field value is not a list")

	_, err = (&Field{Kind: HttpRequest}).LoadValue([]any{})
	assertErrContains(t, err, "template config is nil")
	_, err = (&Field{Kind: HttpRequest, required: true, httprequest: &HttpField{}}).LoadValue(nil)
	assertErrContains(t, err, "field is required")
	_, err = (&Field{Kind: HttpRequest, httprequest: &HttpField{}}).LoadValue("bad")
	assertErrContains(t, err, "field value is not a list")
}

func TestFieldLoadValueHttpRequestBranches(t *testing.T) {
	f := &Field{
		ID:          "h",
		Name:        "h",
		Kind:        HttpRequest,
		httprequest: &HttpField{MultiSelect: true},
	}
	v, err := f.LoadValue([]any{"a", "b"})
	if err != nil {
		t.Fatalf("LoadValue failed: %v", err)
	}
	out, err := v.DumpValue()
	if err != nil {
		t.Fatalf("DumpValue failed: %v", err)
	}
	if !reflect.DeepEqual(out, []any{"a", "b"}) {
		t.Fatalf("multi-select should keep all values, got %#v", out)
	}

	optional := &Field{
		ID:          "h2",
		Name:        "h2",
		Kind:        HttpRequest,
		httprequest: &HttpField{MultiSelect: false},
	}
	v2, err := optional.LoadValue(nil)
	if err != nil {
		t.Fatalf("optional nil load failed: %v", err)
	}
	out2, err := v2.DumpValue()
	if err != nil {
		t.Fatalf("optional nil dump failed: %v", err)
	}
	if !reflect.DeepEqual(out2, []any{}) {
		t.Fatalf("optional nil should dump empty list, got %#v", out2)
	}
}

func TestFieldLoadListAndDumpMapErrorPaths(t *testing.T) {
	f := &Field{}
	_, err := f.loadList([]any{true})
	assertErrContains(t, err, "field type invalid")

	_, err = f.loadList([]any{map[string]any{"name": "bad", "kind": "string"}})
	assertErrContains(t, err, "id is required")

	_, err = f.dumpMap(func() {})
	assertErrContains(t, err, "unsupported type")
}

func TestFieldDumpValueHttpRequestNil(t *testing.T) {
	out, err := (&Field{Kind: HttpRequest}).DumpValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}
