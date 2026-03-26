# go-schema
A library for Markshal/Unmarshal JSON Schema, designed to facilitate flexible form rendering

## Usage

```go
go get github.com/hjhsamuel/go-schema
```

Supported field kind:

- select  
  single select
- selectarray  
  multiple select
- string  
  line edit
- text  
  text edit
- number  
  number edit
- password  
  password edit
- append  
  array line edit
- array  
  nest fields
- nestselect  
  nest single field select
- nestselectarray  
  nest multiple fields select
- incrementarray  
  increment nest field
- httprequest  
  http request


## Schema

all kinds of fields have base json:

```json
{
  "id": "", // schema identifier
  "name": "", // display information
  "fields": []  // subfields
}
```

## Field

```json
{
  "id": "",
  "name": "",
  "kind": "",
  "required": false,
  "desc": ""
}
```

base key:

- id  
field identifier
- name  
displayed field name
- kind  
field kind

value limit:

- required  
if field required
- desc  
description



### select  

- template input

```json
{
  "select": [
    "A",
    "B",
    "C"
  ]
}
```

- template output

```json
{
  "select": [
    "A",
    "B",
    "C"
  ]
}
```

- value input

```json
{
  "select": "A"
}
```

- value output

```json
{
  "select": "A"
}
```

### selectarray  

- template input
- template out
- value input
- value output

### string  

- template input
- template out
- value input
- value output

### text  

- template input
- template out
- value input
- value output

### number  

- template input
- template out
- value input
- value output

### password  

- template input
- template out
- value input
- value output

### append  

- template input
- template out
- value input
- value output

### array  

- template input
- template out
- value input
- value output

### nestselect  

- template input
- template out
- value input
- value output

### nestselectarray  

- template input
- template out
- value input
- value output

### incrementarray  

- template input
- template out
- value input
- value output

### httprequest 

- template input
- template out
- value input
- value output


