package d

import "github.com/stretchr/testify/assert"

// Try calls f() and recovers from any panics caused by using d.Exp. The error generated by the call to d.Exp is returned.
func Try(f func()) (err error) {
	defer nomsRecover(&err)
	f()
	return
}

type UsageError struct {
	msg string
}

func (e UsageError) Error() string {
	return e.msg
}

func nomsRecover(errp *error) {
	if r := recover(); r != nil {
		if _, ok := r.(UsageError); !ok {
			panic(r)
		}
		*errp = r.(error)
	}
}

func IsUsageError(a *assert.Assertions, f func()) {
	e := Try(f)
	a.IsType(UsageError{}, e)
}
