// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jsonschema supports the reading, writing, and manipulation
// of JSON Schemas.
package jsonschema

// The Schema struct models a JSON Schema and, because schemas are
// defined hierarchically, contains many references to itself.
// All fields are pointers and are nil if the associated values
// are not specified.
type Schema struct {
	Schema *string // $schema
	Id     *string // id keyword used for $ref resolution scope
	Ref    *string // $ref, i.e. JSON Pointers

	// http://json-schema.org/latest/json-schema-validation.html
	// 5.1.  Validation keywords for numeric instances (number and integer)
	MultipleOf       *SchemaNumber
	Maximum          *SchemaNumber
	ExclusiveMaximum *bool
	Minimum          *SchemaNumber
	ExclusiveMinimum *bool

	// 5.2.  Validation keywords for strings
	MaxLength *int64
	MinLength *int64
	Pattern   *string

	// 5.3.  Validation keywords for arrays
	AdditionalItems *SchemaOrBoolean
	Items           *SchemaOrSchemaArray
	MaxItems        *int64
	MinItems        *int64
	UniqueItems     *bool

	// 5.4.  Validation keywords for objects
	MaxProperties        *int64
	MinProperties        *int64
	Required             *[]string
	AdditionalProperties *SchemaOrBoolean
	Properties           *[]*NamedSchema
	PatternProperties    *[]*NamedSchema
	Dependencies         *[]*NamedSchemaOrStringArray

	// 5.5.  Validation keywords for any instance type
	Enumeration *[]SchemaEnumValue
	Type        *StringOrStringArray
	AllOf       *[]*Schema
	AnyOf       *[]*Schema
	OneOf       *[]*Schema
	Not         *Schema
	Definitions *[]*NamedSchema

	// 6.  Metadata keywords
	Title       *string
	Description *string
	Default     *interface{}

	// 7.  Semantic validation with "format"
	Format *string
}

// These helper structs represent "combination" types that generally can
// have values of one type or another. All are used to represent parts
// of Schemas.

// SchemaNumber represents a value that can be either an Integer or a Float.
type SchemaNumber struct {
	Integer *int64
	Float   *float64
}

func NewSchemaNumberWithInteger(i int64) *SchemaNumber {
	result := &SchemaNumber{}
	result.Integer = &i
	return result
}
func NewSchemaNumberWithFloat(f float64) *SchemaNumber {
	result := &SchemaNumber{}
	result.Float = &f
	return result
}

// SchemaOrBoolean represents a value that can be either a Schema or a Boolean.
type SchemaOrBoolean struct {
	Schema  *Schema
	Boolean *bool
}

func NewSchemaOrBooleanWithSchema(s *Schema) *SchemaOrBoolean {
	result := &SchemaOrBoolean{}
	result.Schema = s
	return result
}

func NewSchemaOrBooleanWithBoolean(b bool) *SchemaOrBoolean {
	result := &SchemaOrBoolean{}
	result.Boolean = &b
	return result
}

// StringOrStringArray represents a value that can be either
// a String or an Array of Strings.
type StringOrStringArray struct {
	String      *string
	StringArray *[]string
}

func NewStringOrStringArrayWithString(s string) *StringOrStringArray {
	result := &StringOrStringArray{}
	result.String = &s
	return result
}

func NewStringOrStringArrayWithStringArray(a []string) *StringOrStringArray {
	result := &StringOrStringArray{}
	result.StringArray = &a
	return result
}

// SchemaOrStringArray represents a value that can be either
// a Schema or an Array of Strings.
type SchemaOrStringArray struct {
	Schema      *Schema
	StringArray *[]string
}

// SchemaOrSchemaArray represents a value that can be either
// a Schema or an Array of Schemas.
type SchemaOrSchemaArray struct {
	Schema      *Schema
	SchemaArray *[]*Schema
}

func NewSchemaOrSchemaArrayWithSchema(s *Schema) *SchemaOrSchemaArray {
	result := &SchemaOrSchemaArray{}
	result.Schema = s
	return result
}

func NewSchemaOrSchemaArrayWithSchemaArray(a []*Schema) *SchemaOrSchemaArray {
	result := &SchemaOrSchemaArray{}
	result.SchemaArray = &a
	return result
}

// SchemaEnumValue represents a value that can be part of an
// enumeration in a Schema.
type SchemaEnumValue struct {
	String *string
	Bool   *bool
}

// NamedSchema is a name-value pair that is used to emulate maps
// with ordered keys.
type NamedSchema struct {
	Name  string
	Value *Schema
}

func NewNamedSchema(name string, value *Schema) *NamedSchema {
	return &NamedSchema{Name: name, Value: value}
}

// NamedSchemaOrSringArray is a name-value pair that is used
// to emulate maps with ordered keys.
type NamedSchemaOrStringArray struct {
	Name  string
	Value *SchemaOrStringArray
}

// Access named subschemas by name

func namedSchemaArrayElementWithName(array *[]*NamedSchema, name string) *Schema {
	if array == nil {
		return nil
	}
	for _, pair := range *array {
		if pair.Name == name {
			return pair.Value
		}
	}
	return nil
}

func (s *Schema) PropertyWithName(name string) *Schema {
	return namedSchemaArrayElementWithName(s.Properties, name)
}

func (s *Schema) PatternPropertyWithName(name string) *Schema {
	return namedSchemaArrayElementWithName(s.PatternProperties, name)
}

func (s *Schema) DefinitionWithName(name string) *Schema {
	return namedSchemaArrayElementWithName(s.Definitions, name)
}

func (s *Schema) AddProperty(name string, property *Schema) {
	*s.Properties = append(*s.Properties, NewNamedSchema(name, property))
}
