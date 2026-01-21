package validation

import (
	"fmt"
	"strings"
	"github.com/go-playground/validator/v10"
)

func ProcessErr(err error) string {
	if err == nil {
		return ""
	}

	// ValidationErrors 是一个错误切片, 它保存了每个字段违反的每个约束信息
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, 0, 3)
		for _, validationErr := range validationErrs {
			msgs = append(msgs, fmt.Sprintf("字段 [%s] 不满足条件 [%s]",
			validationErr.Field(), validationErr.Tag()))
		}
		return strings.Join(msgs, ";")
	} else {
		return "invalid error"
	}
}