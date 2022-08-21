package rerrors

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

/*
	Rerrors
	Rerrors is errors (which is built in error)
	error types:
	use string not int, clearly
	1, config 	error
	2, sys		error
	3, runtime 	error
	4, internal resource error
	5, external resource error
	6, unexpected 		 error
	7, custom 		 	 error
	8, other 		 	 error
*/

type Rerrors struct {
	data      interface{}
	errorType string
	cause     error
	code      string
}

const (
	ErrorTypeConfig      = "errorTypeConfig"
	ErrorTypeSys         = "errorTypeSys"
	ErrorTypeRuntime     = "errorTypeRuntime"
	ErrorTypeInternalRes = "errorTypeInternalRes"
	ErrorTypeExternalRes = "errorTypeExternalRes"
	ErrorTypeUnexpected  = "errorTypeUnexpected"
	ErrorTypeCustom      = "errorTypeCustom"
	ErrorTypeOther       = "errorTypeOther"
)

func NewErrors(msg string, errorType string, code string) *Rerrors {
	return NewErrorsWithData(msg, errorType, code, nil)
}

func NewErrorsWithData(msg string, errorType string, code string, data interface{}) *Rerrors {
	err := errors.New(msg)
	return &Rerrors{
		cause:     err,
		errorType: errorType,
		data:      data,
		code:      code,
	}
}

func WrapErrors(err error, msg string, errorType string, code string) *Rerrors {
	return WrapErrorsWithData(err, msg, errorType, code, nil)
}
func WrapErrorsWithData(err error, msg string, errorType string, code string, data interface{}) *Rerrors {
	if err == nil {
		return nil
	}
	newerr := errors.Wrap(err, msg)
	return &Rerrors{
		cause:     newerr,
		errorType: errorType,
		data:      data,
		code:      code,
	}
}

func (err *Rerrors) Error() string {
	return "errorType:" + err.errorType + " errMsg:" + err.cause.Error()
}

func (err *Rerrors) Cause() error {
	return err.cause
}

func (err *Rerrors) Type() string {
	return err.errorType
}

func (err *Rerrors) Data() interface{} {
	return err.data
}

/*
 * Format func formats errors
 * not the same as errors
 * v|+v -> print stack exactly
 * s| print f.Error()
 */
func (err *Rerrors) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(s, "errorType:%v \n", err.errorType)
		fmt.Fprintf(s, "errorCode:%v \nmsg:", err.code)
		fmt.Fprintf(s, "%+v\n", err.Cause())
		fmt.Fprintf(s, "errorData:%+v\n", err.Data())
		return
	case 's':
		io.WriteString(s, err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", err.Error())
	}
}

/*
func (f *Rerrors) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "errorType:%v \n", f.errorType)
			fmt.Fprintf(s, "errorCode:%v \nmsg:", f.code)
			fmt.Fprintf(s, "%+v", f.Cause())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.Error())
	case 'q':
		fmt.Fprintf(s, "%q", f.Error())
	}
}
*/
