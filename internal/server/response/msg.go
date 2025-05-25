package response

const (
	cErrForm uint8 = iota + 1
	cErrDBOperation
	cErrReachLimit
	cErrUnexpected
	_
	cErrNotFound
	cErrTelegram
	cErrResourceConflict
)

type Msg struct {
	Code uint8       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

var (
	ErrForm = &Msg{
		Code: cErrForm,
		Msg:  "illegal request parameter",
	}
	ErrDBOperation = &Msg{
		Code: cErrDBOperation,
		Msg:  "server internal error",
	}
	ErrReachLimit = &Msg{
		Code: cErrReachLimit,
		Msg:  "server busy",
	}
	ErrUnexpected = &Msg{
		Code: cErrUnexpected,
		Msg:  "unexpected error",
	}
	ErrNotFound = &Msg{
		Code: cErrNotFound,
		Msg:  "not found",
	}
	ErrTelegram = &Msg{
		Code: cErrTelegram,
		Msg:  "telegram error",
	}
	ErrResourceConflict = &Msg{
		Code: cErrResourceConflict,
		Msg:  "resource conflict",
	}
)
