package errors

import (
	"github.com/Binaretech/classroom-auth/lang"
	"github.com/gofiber/fiber/v2"
)

type ServerError interface {
	error
	GetCode() uint
	GetMessage() interface{}
}

type Error struct {
	msg  interface{}
	code uint
}

func (err *Error) Error() string {
	return ""
}

func (err *Error) GetMessage() interface{} {
	return err.msg
}

func (err *Error) GetCode() uint {
	return err.code
}

func New(msg interface{}, code uint) *Error {
	return &Error{msg: msg, code: code}
}

type WrappedError struct {
	err  error
	msg  string
	code uint
}

func (err *WrappedError) Error() string {
	return err.err.Error()
}

func (err *WrappedError) GetCode() uint {
	return err.code
}

func (err *WrappedError) GetMessage() interface{} {
	return err.msg
}

func WrapError(e error, message string, code uint) *WrappedError {
	return &WrappedError{err: e, code: code, msg: message}
}

type InternalError struct {
	err error
}

func NewInternalError(err error) *InternalError {
	return &InternalError{err: err}
}

func (err *InternalError) Error() string {
	return err.err.Error()
}

func (err *InternalError) GetCode() uint {
	return fiber.StatusInternalServerError
}

func (err *InternalError) GetMessage() interface{} {
	return lang.Trans("internal error")
}

type BadRequest struct {
	msg interface{}
	err error
}

func NewBadRequest(msg interface{}) *BadRequest {
	return &BadRequest{msg: msg}
}

func NewBadRequestWrapped(msg interface{}, err error) *BadRequest {
	return &BadRequest{msg: msg, err: err}
}

func (err *BadRequest) Error() string {
	if err.err == nil {
		return ""
	}

	return err.err.Error()
}

func (err *BadRequest) GetCode() uint {
	return fiber.StatusBadRequest
}

func (err *BadRequest) GetMessage() interface{} {
	return err.msg
}

type Unauthenticated struct {
}

func NewUnauthenticatedError() *Unauthenticated {
	return &Unauthenticated{}
}

func (err *Unauthenticated) Error() string {
	return ""
}

func (err *Unauthenticated) GetCode() uint {
	return fiber.StatusUnauthorized
}

func (err *Unauthenticated) GetMessage() interface{} {
	return lang.Trans("unauthenticated")
}
