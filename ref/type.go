package ref

import (
	"context"
	"reflect"
)

// errorType error 的反射类型。
// errorType the reflection type of error.
var errorType = reflect.TypeOf((*error)(nil)).Elem()

// contextType context.Context 的反射类型。
// contextType the reflection type of context.Context.
var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

// IsErrorType t 是否是 error 类型。
// IsErrorType returns whether `t` is error type.
func IsErrorType(t reflect.Type) bool {
	return t == errorType || t.Implements(errorType)
}

// IsContextType t 是否是 context.Context 类型。
// IsContextType returns whether `t` is context.Context type.
func IsContextType(t reflect.Type) bool {
	return t == contextType || t.Implements(contextType)
}

// IsFuncType t 是否是 func 类型。
// IsFuncType returns whether `t` is func type.
func IsFuncType(t reflect.Type) bool {
	return t.Kind() == reflect.Func
}

// IsStructPointer 返回是否是结构体的指针类型。
// ptr 一般指一层指针，因为多层指针在 web 开发中很少用到，甚至应该在纯业务代码中禁止使用多层指针。
// IsStructPointer returns whether it is the pointer type of structure.
func IsStructPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct
}
