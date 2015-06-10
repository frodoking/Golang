package binding

import (
	"gopkg.in/bluesuncorp/validator.v5"
	"reflect"
	"sync"
)

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ StructValidator = &DefaultValidator{}

func (this *DefaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		this.lazyInit()
		if err := this.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

func (this *DefaultValidator) lazyInit() {
	this.once.Do(func() {
		this.validate = validator.New("binding", validator.BakedInValidators)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}

	return valueType
}
