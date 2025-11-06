package go_schema

type FieldKind string

const (
	Select      FieldKind = "select"      // single select
	SelectArray FieldKind = "selectarray" // multiple select

	String   FieldKind = "string"   // line edit
	Text     FieldKind = "text"     // text edit
	Number   FieldKind = "number"   // number edit
	Password FieldKind = "password" // password edit
	Append   FieldKind = "append"   // array line edit

	Array           FieldKind = "array"           // nest fields
	NestSelect      FieldKind = "nestselect"      // nest single field select
	NestSelectArray FieldKind = "nestselectarray" // nest multiple fields select
	IncrementArray  FieldKind = "incrementarray"  // increment nest field
)
