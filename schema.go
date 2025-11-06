package go_schema

import (
	"github.com/pkg/errors"
)

type Schema struct {
	ID     string
	Name   string
	Fields []*Field
}

func (s *Schema) LoadTemplate(val map[string]any) error {
	if v, ok := val["id"].(string); ok {
		s.ID = v
	} else {
		return errors.New("id is required or value is not string")
	}

	if v, ok := val["name"].(string); ok {
		s.Name = v
	} else {
		return errors.New("name is required or value is not string")
	}

	if v, ok := val["fields"]; ok {
		if fields, fok := v.([]any); fok {
			s.Fields = make([]*Field, 0)
			for _, field := range fields {
				if fieldVal, vok := field.(map[string]any); vok {
					child := &Field{}
					if cErr := child.LoadTemplate(fieldVal); cErr != nil {
						return cErr
					}
					s.Fields = append(s.Fields, child)
				} else {
					return errors.New("field elem is not a map")
				}
			}
		}
	}
	return nil
}

func (s *Schema) DumpTemplate() map[string]any {
	res := make(map[string]any)
	res["id"] = s.ID
	res["name"] = s.Name

	fields := make([]map[string]any, 0)
	for _, field := range s.Fields {
		val := field.DumpTemplate()
		fields = append(fields, val)
	}
	res["fields"] = fields
	return res
}

func (s *Schema) LoadValue(val map[string]any) (*Schema, error) {
	v, ok := val[s.ID]
	if !ok {
		return nil, errors.Errorf("module %s not found", s.ID)
	}

	res := &Schema{
		ID:     s.ID,
		Name:   s.Name,
		Fields: make([]*Field, 0),
	}
	fieldVals, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("field value is not a map")
	}

	for _, field := range s.Fields {
		fieldVal, dok := fieldVals[field.ID]
		if field.required && !dok {
			return nil, errors.Errorf("field %s is required", field.ID)
		}

		f, err := field.LoadValue(fieldVal)
		if err != nil {
			return nil, err
		}

		res.Fields = append(res.Fields, f)
	}

	return res, nil
}

func (s *Schema) DumpValue() (map[string]any, error) {
	fieldMap := make(map[string]any)
	for _, field := range s.Fields {
		v, err := field.DumpValue()
		if err != nil {
			return nil, err
		}
		fieldMap[field.ID] = v
	}
	res := map[string]any{
		s.ID: fieldMap,
	}
	return res, nil
}
