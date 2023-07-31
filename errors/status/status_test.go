package status

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func printUnicode(code uint64) {
	uc := fmt.Sprintf("\\u%x", code)
	str, err := strconv.Unquote(`"` + uc + `"`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d %s %s \n", code, uc, str)

}
func TestError(t *testing.T) {
	e := Error("test error")
	fmt.Println(e)

	printUnicode(38382)
	printUnicode(65510)
}

func TestNew(t *testing.T) {
	errcode := "test_code"
	err1msg := "test error 1"
	err1 := New(errcode, err1msg)
	assert.NotNil(t, err1)
	assert.NotEmpty(t, err1.Error())
	assert.Equal(t, nil, err1.Unwrap())
	assert.Equal(t, errcode, err1.Code())
	assert.Equal(t, err1msg, err1.Message())
	assert.Equal(t, []any(nil), err1.Details())
	fmt.Println(err1)

	suberr := errors.New("sub error")
	err2 := New(errcode, suberr)
	assert.NotNil(t, err2)
	assert.NotEmpty(t, err2.Error())
	assert.Equal(t, suberr, err2.Unwrap())
	assert.Equal(t, errcode, err2.Code())
	assert.Equal(t, "", err2.Message())
	assert.Equal(t, []any(nil), err2.Details())
	fmt.Println(err2)

	err3 := New(errcode, err1msg).SetError(suberr)
	assert.NotNil(t, err3)
	assert.NotEmpty(t, err3.Error())
	assert.Equal(t, suberr, err3.Unwrap())
	assert.Equal(t, errcode, err3.Code())
	assert.Equal(t, err1msg, err3.Message())
	assert.Equal(t, []any(nil), err3.Details())
	fmt.Println(err3)
}
