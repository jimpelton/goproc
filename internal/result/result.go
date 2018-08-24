package result

type Result struct {
	r chan interface{}
}

func NewResult() *Result {
	return &Result{
		r: make(chan interface{}, 1),
	}
}

func (r *Result) Done(v interface{}) {
	r.r <- v
}

func (r *Result) Wait() interface{} {
	return <-r.r
}
