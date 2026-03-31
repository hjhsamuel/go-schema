# go-schema
A library for Markshal/Unmarshal JSON Schema, designed to facilitate flexible form rendering

## Usage

```shell
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

- template

    ```json
    {
      "select": [
        "A",
        "B",
        "C"
      ]
    }
    ```

- value

    ```json
    {
      "{id}": "A"
    }
    ```

### selectarray  

- template

    ```json
    {
      "selectarray": [
        "A",
        "B",
        "C"
      ]
    }
    ```

- value

    ```json
    {
      "{id}": [
        "A",
        "B"
      ]
    }
    ```

### string  

- template
- value

    ```json
    {
      "{id}": "the value of string field"
    }
    ```

### text  

- template
- value

    ```json
    {
      "{id}": "the value of text field"
    }
    ```

### number  

- template
- value

    ```json
    {
      "{id}": 1.0
    }
    ```

### password  

- template
- value

    ```json
    {
      "{password}": "the value of password field"
    }
    ```

### append  

- template
- value

    ```json
    {
      "{id}": [
        "A",
        "B"
      ]
    }
    ```

### array  

- template

    ```json
    {
      "array": [
        {
          "id": "{child A id}",
          "name": "{child A name}",
          "kind": "{child A kind}",
          ...
        },
        ...
      ]
    }
    ```

- value

    ```json
    {
      "{id}": {
        "{child A id}": ...,
        "{child B id}": ...
      }
    }
    ```

### nestselect  

- template

    ```json
    {
      "nestselect": [
        {
          "id": "{child A id}",
          "name": "{child A name}",
          "kind": "{child A kind}",
          ...
        },
        ...
      ]
    }
    ```

- value

    ```json
    {
      "{id}": {
        "{child A id}": ...
      }
    }
    ```

### nestselectarray  

- template

    ```json
    {
      "nestselectarray": [
        {
          "id": "{child A id}",
          ...
        },
        {
          "id": "{child B id}",
          ...
        },
        ...
      ]
    }
    ```

- value

    ```json
    {
      "{id}": {
        "{child A id}": ...,
        "{child B id}": ...,
        ...
      }
    }
    ```

### incrementarray  

- template

    ```json
    {
      "incrementarray": [
        {
          "id": "{child id}",
          ...
        }
      ]
    }
    ```

- value

    ```json
    {
      "{id}": [
        {
          "{child id}": ...
        },
        {
          "{child id}": ...
        },
        ...
      ]
    }
    ```

### httprequest 

- template

    ```json
    {
      "httprequest": {
        "url": "http://xxx/xx/xx?p1=a&p2=b",
        "method": "get",
        "multi_select": true,
        "user_data": "{user custom data}"
      }
    }
    ```

- value

    ```json
    {
      "{id}": [
        ...,
        ...
      ]
    }
    ```


## Example

```json

```