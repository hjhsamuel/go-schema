package go_schema

import (
	"fmt"

	"github.com/pkg/errors"
)

type Field struct {
	ID   string
	Name string
	Kind FieldKind

	// value limit
	required bool
	desc     string

	// validate
	validator string

	// optional kind
	array              []*Field
	selectOne          []string
	selectMultiple     []string
	str                string
	text               string
	number             float64
	password           string
	nestSelectOne      []*Field
	nestSelectMultiple []*Field
	incrementArray     []*Field
	append             []string
}

func (f *Field) LoadTemplate(val map[string]any) error {
	if v, ok := val["id"].(string); ok {
		f.ID = v
	} else {
		return errors.New("id is required or value is not string")
	}

	if v, ok := val["name"].(string); ok {
		f.Name = v
	} else {
		return errors.New("name is required or value is not string")
	}

	kind, ok := val["kind"].(string)
	if !ok {
		return errors.New("kind is required or value is not string")
	} else {
		f.Kind = FieldKind(kind)
	}

	kindVal, ok := val[kind]
	if !ok {
		kindVal = nil
	}

	switch f.Kind {
	case Array:
		if v, vok := kindVal.([]any); vok {
			children, err := f.loadList(v)
			if err != nil {
				return err
			}
			f.array = make([]*Field, 0)
			for _, child := range children {
				if childField, cok := child.(*Field); cok {
					f.array = append(f.array, childField)
				} else {
					return errors.Errorf("field %s array element type is not `Field`", f.ID)
				}
			}
		} else {
			return errors.Errorf("field kind %s is not list", kind)
		}
	case Select:
		if v, vok := kindVal.([]any); vok {
			children, err := f.loadList(v)
			if err != nil {
				return err
			}
			f.selectOne = make([]string, 0)
			for _, child := range children {
				if childVal, cok := child.(string); cok {
					f.selectOne = append(f.selectOne, childVal)
				} else {
					return errors.Errorf("field %s select element type is not string", f.ID)
				}
			}
		} else {
			return errors.Errorf("field kind %s is not list", kind)
		}
	case SelectArray:
		if v, vok := kindVal.([]any); vok {
			children, err := f.loadList(v)
			if err != nil {
				return err
			}
			f.selectMultiple = make([]string, 0)
			for _, child := range children {
				if childVal, cok := child.(string); cok {
					f.selectMultiple = append(f.selectMultiple, childVal)
				} else {
					return errors.Errorf("field %s selectarray element type is not string", f.ID)
				}
			}
		} else {
			return errors.Errorf("field kind %s is not list", kind)
		}
	case NestSelect:
		if v, vok := kindVal.([]any); vok {
			children, err := f.loadList(v)
			if err != nil {
				return err
			}
			f.nestSelectOne = make([]*Field, 0)
			for _, child := range children {
				if childField, cok := child.(*Field); cok {
					f.nestSelectOne = append(f.nestSelectOne, childField)
				} else {
					return errors.Errorf("field %s nestselect element type is not `Field`", f.ID)
				}
			}
		} else {
			return errors.Errorf("field kind %s is not list", kind)
		}
	case NestSelectArray:
		if v, vok := kindVal.([]any); vok {
			children, err := f.loadList(v)
			if err != nil {
				return err
			}
			f.nestSelectMultiple = make([]*Field, 0)
			for _, child := range children {
				if childField, cok := child.(*Field); cok {
					f.nestSelectMultiple = append(f.nestSelectMultiple, childField)
				} else {
					return errors.Errorf("field %s nestselectarray element type is not `Field`", f.ID)
				}
			}
		} else {
			return errors.Errorf("field kind %s is not list", kind)
		}
	case IncrementArray:
		if v, vok := kindVal.([]any); vok {
			children, err := f.loadList(v)
			if err != nil {
				return err
			}
			f.incrementArray = make([]*Field, 0)
			for _, child := range children {
				if childField, cok := child.(*Field); cok {
					f.incrementArray = append(f.incrementArray, childField)
				} else {
					return errors.Errorf("field %s incrementarray element type is not `Field`", f.ID)
				}
			}
		} else {
			return errors.Errorf("field kind %s is not list", kind)
		}
	case String, Text, Number, Password, Append:
	default:
		return errors.Errorf("unsupported field kind %s", f.Kind)
	}

	if v, vok := val["required"].(bool); vok {
		f.required = v
	}
	if v, vok := val["desc"].(string); vok {
		f.desc = v
	}
	if v, vok := val["validator"].(string); vok {
		f.validator = v
	}

	return nil
}

func (f *Field) loadList(val []any) ([]any, error) {
	res := make([]any, 0)
	for _, field := range val {
		switch field.(type) {
		case string, float64:
			res = append(res, field)
		default:
			if childVal, ok := field.(map[string]any); ok {
				child := &Field{}
				if err := child.LoadTemplate(childVal); err != nil {
					return nil, err
				}
				res = append(res, child)
			} else {
				return nil, errors.New("field type invalid")
			}
		}
	}
	return res, nil
}

func (f *Field) DumpTemplate() map[string]any {
	res := make(map[string]any)
	res["id"] = f.ID
	res["name"] = f.Name
	res["kind"] = f.Kind

	if f.required {
		res["required"] = true
	}
	if f.desc != "" {
		res["desc"] = f.desc
	}

	switch f.Kind {
	case Array:
		v := f.dumpList(f.array)
		res[string(f.Kind)] = v
	case Select:
		res[string(f.Kind)] = f.selectOne
	case SelectArray:
		res[string(f.Kind)] = f.selectMultiple
	case NestSelect:
		v := f.dumpList(f.nestSelectOne)
		res[string(f.Kind)] = v
	case NestSelectArray:
		v := f.dumpList(f.nestSelectMultiple)
		res[string(f.Kind)] = v
	case IncrementArray:
		v := f.dumpList(f.incrementArray)
		res[string(f.Kind)] = v
	case String, Text, Number, Password:
		if f.validator != "" {
			res["validator"] = f.validator
		}
	case Append:

	}

	return res
}

func (f *Field) dumpList(fields []*Field) []any {
	res := make([]any, 0)
	for _, field := range fields {
		res = append(res, field.DumpTemplate())
	}
	return res
}

func (f *Field) LoadValue(val any) (*Field, error) {
	field := &Field{
		ID:   f.ID,
		Name: f.Name,
		Kind: f.Kind,
	}

	switch f.Kind {
	case Array:
		if v, ok := val.(map[string]any); ok {
			field.array = make([]*Field, 0)
			for _, child := range f.array {
				if childVal, cok := v[child.ID]; cok {
					childField, err := child.LoadValue(childVal)
					if err != nil {
						return nil, err
					}
					field.array = append(field.array, childField)
				} else if child.required {
					return nil, errors.Errorf("field %s is required", f.ID)
				}
			}
		} else {
			return nil, errors.New("array: field value is not a map")
		}
	case Select:
		if v, ok := val.(string); ok {
			field.selectOne = make([]string, 0)
			for _, selectVal := range f.selectOne {
				if selectVal == v {
					field.selectOne = append(field.selectOne, selectVal)
					break
				}
			}
			if len(field.selectOne) == 0 {
				return nil, errors.New("select value invalid")
			}
		} else {
			return nil, errors.New("select: field value is not a string")
		}
	case SelectArray:
		if v, ok := val.([]any); ok {
			field.selectMultiple = make([]string, 0)
			for _, elem := range v {
				elemStr := fmt.Sprintf("%v", elem)
				matched := false
				for _, arrayVal := range f.selectMultiple {
					if arrayVal == elemStr {
						matched = true
						break
					}
				}
				if !matched {
					return nil, errors.Errorf("select value '%s' invalid", elemStr)
				}
				field.selectMultiple = append(field.selectMultiple, elemStr)
			}
		} else {
			return nil, errors.New("selectarray: field value is not list")
		}
	case String:
		if v, ok := val.(string); ok {
			field.str = v
		} else {
			return nil, errors.New("string: field value is not a string")
		}
	case Text:
		if v, ok := val.(string); ok {
			field.text = v
		} else {
			return nil, errors.New("text: field value is not a string")
		}
	case Number:
		if v, ok := val.(float64); ok {
			field.number = v
		} else {
			return nil, errors.New("number: field value is not a number")
		}
	case Password:
		if v, ok := val.(string); ok {
			field.password = v
		} else {
			return nil, errors.New("password: field value is not a string")
		}
	case NestSelect:
		if v, ok := val.(map[string]any); ok {
			for _, childField := range f.nestSelectOne {
				if childVal, cok := v[childField.ID]; cok {
					out, err := childField.LoadValue(childVal)
					if err != nil {
						return nil, err
					}
					field.nestSelectOne = append(field.nestSelectOne, out)
					break
				} else if childField.required {
					return nil, errors.Errorf("field %s is required", f.ID)
				}
			}
		} else {
			return nil, errors.New("nestselect: field value is not a map")
		}
	case NestSelectArray:
		if v, ok := val.(map[string]any); ok {
			for _, childField := range f.nestSelectMultiple {
				if childVal, cok := v[childField.ID]; cok {
					out, err := childField.LoadValue(childVal)
					if err != nil {
						return nil, err
					}
					field.nestSelectMultiple = append(field.nestSelectMultiple, out)
				} else if childField.required {
					return nil, errors.Errorf("field %s is required", f.ID)
				}
			}
		} else {
			return nil, errors.New("nestselectarray: field value is not a map")
		}
	case IncrementArray:
		if v, ok := val.(map[string]any); ok {
			for _, childField := range f.incrementArray {
				if childVal, cok := v[childField.ID]; cok {
					if listVal, lok := childVal.([]any); lok {
						for _, elem := range listVal {
							out, err := childField.LoadValue(elem)
							if err != nil {
								return nil, err
							}
							field.incrementArray = append(field.incrementArray, out)
						}
					} else {
						return nil, errors.New("incrementarray: field value is not a list")
					}
				}
			}
		} else {
			return nil, errors.New("incrementArray: field value is not a map")
		}
	case Append:
		if v, ok := val.([]any); ok {
			field.append = make([]string, len(v))
			for i, elem := range v {
				field.append[i] = fmt.Sprintf("%v", elem)
			}
		} else {
			return nil, errors.New("append: field value is not a list")
		}
	}
	return field, nil
}

func (f *Field) DumpValue() (any, error) {
	switch f.Kind {
	case Array:
		res := make(map[string]any)
		for _, field := range f.array {
			v, err := field.DumpValue()
			if err != nil {
				return nil, err
			}
			res[field.ID] = v
		}
		return res, nil
	case Select:
		return f.selectOne[0], nil
	case SelectArray:
		return f.selectMultiple, nil
	case String:
		return f.str, nil
	case Text:
		return f.text, nil
	case Number:
		return f.number, nil
	case Password:
		return f.password, nil
	case NestSelect:
		child := f.nestSelectOne[0]
		childVal, err := child.DumpValue()
		if err != nil {
			return nil, err
		}
		return map[string]any{child.ID: childVal}, nil
	case NestSelectArray:
		res := make(map[string]any)
		for _, field := range f.nestSelectMultiple {
			v, err := field.DumpValue()
			if err != nil {
				return nil, err
			}
			res[field.ID] = v
		}
		return res, nil
	case IncrementArray:
		res := make(map[string][]any)
		for _, field := range f.incrementArray {
			if _, ok := res[field.ID]; !ok {
				res[field.ID] = make([]any, 0)
			}
			v, err := field.DumpValue()
			if err != nil {
				return nil, err
			}
			res[field.ID] = append(res[field.ID], v)
		}
		return res, nil
	case Append:
		return f.append, nil
	default:
		return nil, errors.Errorf("unsupported field kind '%s'", f.Kind)
	}
}
