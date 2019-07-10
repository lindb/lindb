package apierrors

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_BadRequestException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "BadRequestException:test1", r.Error())
			}

		}
	}()
	NewBadRequestException("test1")
}

func Test_BadGatewayException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "BadGatewayException:test1", r.Error())
			}

		}
	}()
	NewBadGatewayException("test1")
}

func Test_RpcException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "RPCException:test1", r.Error())
			}

		}
	}()
	NewRPCException("test1")
}

func Test_AuthException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "AuthException:test1", r.Error())
			}

		}
	}()
	NewAuthException("test1")
}

func Test_ExprException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "ExprException:test1", r.Error())
			}

		}
	}()
	NewExprException("test1")
}

func Test_IndexException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "IndexException:test1", r.Error())
			}

		}
	}()
	NewIndexException("test1")
}

func Test_LinException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "LinException:test1", r.Error())
			}

		}
	}()
	NewLinException("test1")
}

func Test_LinSQLException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "LinSQLException:test1", r.Error())
			}

		}
	}()
	NewLinSQLException("test1")
}

func Test_NoLeaderException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "NoLeaderException:test1", r.Error())
			}

		}
	}()
	NewNoLeaderException("test1")
}

func Test_NoPartitionException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "NoPartitionException:test1", r.Error())
			}

		}
	}()
	NewNoPartitionException("test1")
}

func Test_NoShardException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "NoShardException:test1", r.Error())
			}

		}
	}()
	NewNoShardException("test1")
}

func Test_QueryException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "QueryException:test1", r.Error())
			}

		}
	}()
	NewQueryException("test1")
}

func Test_TimeoutException(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			r, ok := err.(error)
			if ok {
				assert.Equal(t, "TimeoutException:test1", r.Error())
			}

		}
	}()
	NewTimeoutException("test1")
}
