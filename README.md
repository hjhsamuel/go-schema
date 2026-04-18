# go-schema

A Go library for marshaling and unmarshaling schema-driven form definitions and values.

## Install

```bash
go get github.com/hjhsamuel/go-schema
```

## Supported Field Kinds

| Kind | Description |
| --- | --- |
| `select` | single select |
| `selectarray` | multiple select |
| `string` | single-line text |
| `text` | multi-line text |
| `number` | number input |
| `password` | password input |
| `append` | string list input |
| `array` | nested fields |
| `nestselect` | select one nested field set |
| `nestselectarray` | select multiple nested field sets |
| `incrementarray` | repeatable nested field set |
| `httprequest` | remote options selector |

## Quick Start

```go
package main

import (
	"fmt"

	go_schema "github.com/hjhsamuel/go-schema"
)

func main() {
	template := map[string]any{
		"id":   "moduleA",
		"name": "Module A",
		"fields": []any{
			map[string]any{
				"id":       "title",
				"name":     "Title",
				"kind":     "string",
				"required": true,
			},
			map[string]any{
				"id":     "level",
				"name":   "Level",
				"kind":   "select",
				"select": []any{"low", "high"},
			},
		},
	}

	schema := &go_schema.Schema{}
	if err := schema.LoadTemplate(template); err != nil {
		panic(err)
	}

	values := map[string]any{
		"moduleA": map[string]any{
			"title": "demo",
			"level": "high",
		},
	}

	loaded, err := schema.LoadValue(values)
	if err != nil {
		panic(err)
	}

	out, err := loaded.DumpValue()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", out)
}
```

## Schema Template Format

```json
{
  "id": "module_id",
  "name": "Module Name",
  "fields": []
}
```

## Field Template Base Keys

```json
{
  "id": "field_id",
  "name": "Field Name",
  "kind": "string",
  "required": false,
  "desc": "description",
  "validator": "optional validator tag"
}
```

- `id`: unique field identifier
- `name`: display name
- `kind`: field type
- `required`: optional, default `false`
- `desc`: optional field description
- `validator`: optional validator metadata (used for `string`/`text`/`number`/`password`)

## Kind-Specific Template and Value Examples

### `select`

Template:

```json
{
  "kind": "select",
  "select": ["A", "B", "C"]
}
```

Value:

```json
{
  "{id}": "A"
}
```

### `selectarray`

Template:

```json
{
  "kind": "selectarray",
  "selectarray": ["A", "B", "C"]
}
```

Value:

```json
{
  "{id}": ["A", "B"]
}
```

### `string` / `text` / `password`

Template:

```json
{
  "kind": "string"
}
```

Value:

```json
{
  "{id}": "input value"
}
```

### `number`

Template:

```json
{
  "kind": "number"
}
```

Value:

```json
{
  "{id}": 1.0
}
```

### `append`

Template:

```json
{
  "kind": "append"
}
```

Value:

```json
{
  "{id}": ["A", "B"]
}
```

### `array`

Template:

```json
{
  "kind": "array",
  "array": [
    { "id": "childA", "name": "Child A", "kind": "string" },
    { "id": "childB", "name": "Child B", "kind": "number" }
  ]
}
```

Value:

```json
{
  "{id}": {
    "childA": "text",
    "childB": 1.0
  }
}
```

### `nestselect`

Template:

```json
{
  "kind": "nestselect",
  "nestselect": [
    { "id": "optA", "name": "Option A", "kind": "string" },
    { "id": "optB", "name": "Option B", "kind": "number" }
  ]
}
```

Value (select one):

```json
{
  "{id}": {
    "optA": "value"
  }
}
```

### `nestselectarray`

Template:

```json
{
  "kind": "nestselectarray",
  "nestselectarray": [
    { "id": "optA", "name": "Option A", "kind": "string" },
    { "id": "optB", "name": "Option B", "kind": "number" }
  ]
}
```

Value (select multiple):

```json
{
  "{id}": {
    "optA": "value",
    "optB": 2.0
  }
}
```

### `incrementarray`

Template:

```json
{
  "kind": "incrementarray",
  "incrementarray": [
    { "id": "child", "name": "Child", "kind": "string" }
  ]
}
```

Value (repeatable groups):

```json
{
  "{id}": [
    { "child": "v1" },
    { "child": "v2" }
  ]
}
```

### `httprequest`

Template:

```json
{
  "kind": "httprequest",
  "httprequest": {
    "url": "https://example.com/api/options?p1=a&p2=b",
    "method": "GET",
    "multi_select": true,
    "user_data": "custom data"
  }
}
```

Value:

```json
{
  "{id}": ["option_a", "option_b"]
}
```

## Notes

- For `number`, input should be `float64` when decoded from JSON into `map[string]any`.
- `httprequest.method` supports only `GET` and `POST` (case-insensitive in template input).
- `Schema.LoadValue` expects input format as `{schema_id: {field_id: value}}`.
