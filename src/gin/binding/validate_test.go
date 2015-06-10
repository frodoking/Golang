package binding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Struct1 struct {
	Value float64 `binding:"required"`
}

type Struct2 struct {
	RequiredValue string `binding:"required"`
	Value         float64
}

type Struct3 struct {
	Integer    int
	String     string
	BasicSlice []int
	Boolean    bool

	RequiredInteger       int       `binding:"required"`
	RequiredString        string    `binding:"required"`
	RequiredAnotherStruct Struct1   `binding:"required"`
	RequiredBasicSlice    []int     `binding:"required"`
	RequiredComplexSlice  []Struct2 `binding:"required"`
	RequiredBoolean       bool      `binding:"required"`
}

func NewStruct() Struct3 {
	return Struct3{
		RequiredInteger:       2,
		RequiredString:        "hello",
		RequiredAnotherStruct: Struct1{1.5},
		RequiredBasicSlice:    []int{1, 2, 3, 4, 5},
		RequiredComplexSlice: []Struct2{
			{RequiredValue: "A"},
			{RequiredValue: "B"},
		},
		RequiredBoolean: true,
	}
}

func TestValidateGoodObject(t *testing.T) {
	test := NewStruct()
	assert.Nil(t, validate(&test))
}

type Object map[string]interface{}
type MyObjects []Object

func TestValidateSlice(t *testing.T) {
	var obj MyObjects
	var obj2 Object
	var nu = 10

	assert.NoError(t, validate(obj))
	assert.NoError(t, validate(&obj))
	assert.NoError(t, validate(obj2))
	assert.NoError(t, validate(&obj2))
	assert.NoError(t, validate(nu))
	assert.NoError(t, validate(&nu))
}
