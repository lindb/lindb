package apierrors

import (
	"errors"
	"fmt"
)

type ExceptionCode int

const (
	AuthExceptionType ExceptionCode = iota
	ExprExceptionType
	IndexExceptionType
	LinExceptionType
	LinSQLExceptionType
	NoLeaderExceptionType
	NoPartitionExceptionType
	NoShardExceptionType
	QueryExceptionType
	RPCExceptionType
	TimeoutExceptionType
	BadGatewayExceptionType
	BadRequestExceptionType
)

func (e ExceptionCode) String() string {
	switch e {
	case RPCExceptionType:
		return "RPCException"
	case AuthExceptionType:
		return "AuthException"
	case ExprExceptionType:
		return "ExprException"
	case IndexExceptionType:
		return "IndexException"
	case LinExceptionType:
		return "LinException"
	case LinSQLExceptionType:
		return "LinSQLException"
	case NoLeaderExceptionType:
		return "NoLeaderException"
	case NoPartitionExceptionType:
		return "NoPartitionException"
	case NoShardExceptionType:
		return "NoShardException"
	case QueryExceptionType:
		return "QueryException"
	case TimeoutExceptionType:
		return "TimeoutException"
	case BadGatewayExceptionType:
		return "BadGatewayException"
	case BadRequestExceptionType:
		return "BadRequestException"
	default:
		return "UnknownException"
	}
}

// ErrTooManyTags is the error returned by memory-database when
// writes exceed the max limit of tag identifiers.
var ErrTooManyTags = errors.New("too many tags")

// ErrUserNotExists user not exists
var ErrUserNotExists = errors.New("user not exists")

// ErrInvalidJWTToken invalid jwt token
var ErrInvalidJWTToken = errors.New("invalid jwt token")

// ErrSignatureMismatch means signature did not match.
var ErrSignatureMismatch = errors.New("signature does not match")

// BadRequest represents invalid request
type BadRequestException struct {
	Msg string
}

// Error error message
func (e BadRequestException) Error() error {
	return fmt.Errorf(ErrorFormat(e.Msg, BadRequestExceptionType.String()))
}

func NewBadRequestException(message string) {
	panic(BadRequestException{message}.Error())
}

// BadGateway represents proxy error
type BadGatewayException struct {
	Msg string
}

// Error error message
func (e BadGatewayException) Error() error {
	return fmt.Errorf(ErrorFormat(e.Msg, BadGatewayExceptionType.String()))
}

func NewBadGatewayException(message string) {
	panic(BadGatewayException{message}.Error())
}

type RPCException struct {
	Msg string
}

func (r RPCException) Error() error {
	return fmt.Errorf(ErrorFormat(r.Msg, RPCExceptionType.String()))
}

func NewRPCException(message string) {
	panic(RPCException{message}.Error())
}

type ShutdownHookException struct {
	Msg string
}

type AuthException struct {
	Msg string
}

func (a AuthException) Error() error {
	return fmt.Errorf(ErrorFormat(a.Msg, AuthExceptionType.String()))
}

func NewAuthException(message string) {
	panic(AuthException{message}.Error())
}

type ExprException struct {
	Msg string
}

func (e ExprException) Error() error {
	return fmt.Errorf(ErrorFormat(e.Msg, ExprExceptionType.String()))
}

func NewExprException(message string) {
	panic(ExprException{message}.Error())
}

type IndexException struct {
	Msg string
}

func (i IndexException) Error() error {
	return fmt.Errorf(ErrorFormat(i.Msg, IndexExceptionType.String()))
}

func NewIndexException(message string) {
	panic(IndexException{message}.Error())
}

type LinException struct {
	Msg string
}

func (l LinException) Error() error {
	return fmt.Errorf(ErrorFormat(l.Msg, LinExceptionType.String()))
}

func NewLinException(message string) {
	panic(LinException{message}.Error())
}

type LinSQLException struct {
	Msg string
}

func (l LinSQLException) Error() error {
	return fmt.Errorf(ErrorFormat(l.Msg, LinSQLExceptionType.String()))
}

func NewLinSQLException(message string) {
	panic(LinSQLException{message}.Error())
}

type NoLeaderException struct {
	Msg string
}

func (n NoLeaderException) Error() error {
	return fmt.Errorf(ErrorFormat(n.Msg, NoLeaderExceptionType.String()))
}

func NewNoLeaderException(message string) {
	panic(NoLeaderException{message}.Error())
}

type NoPartitionException struct {
	Msg string
}

func (n NoPartitionException) Error() error {
	return fmt.Errorf(ErrorFormat(n.Msg, NoPartitionExceptionType.String()))
}

func NewNoPartitionException(message string) {
	panic(NoPartitionException{message}.Error())
}

type NoShardException struct {
	Msg string
}

func (n NoShardException) Error() error {
	return fmt.Errorf(ErrorFormat(n.Msg, NoShardExceptionType.String()))
}

func NewNoShardException(message string) {
	panic(NoShardException{message}.Error())
}

type QueryException struct {
	Msg string
}

func (q QueryException) Error() error {
	return fmt.Errorf(ErrorFormat(q.Msg, QueryExceptionType.String()))
}

func NewQueryException(message string) {
	panic(QueryException{message}.Error())
}

type TimeoutException struct {
	Msg string
}

func (t TimeoutException) Error() error {
	return fmt.Errorf(ErrorFormat(t.Msg, TimeoutExceptionType.String()))
}

func NewTimeoutException(message string) {
	panic(TimeoutException{message}.Error())
}

func ErrorFormat(message string, exceptionType string) string {
	return fmt.Sprintf("%s:%s", exceptionType, message)
}
