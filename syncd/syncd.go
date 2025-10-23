package syncd

import "sync"

func IncrementRC(mu *sync.Mutex, total_rc, rc_to_add *int64) {
	mu.Lock()
	*total_rc += *rc_to_add
	mu.Unlock()
}

func CatchErr(mu *sync.Mutex, errs *[]error, err *error) {
	mu.Lock()
	*errs = append(*errs, *err)
	mu.Unlock()
}
