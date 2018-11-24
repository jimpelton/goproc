package analysis

type Body interface {
	// F is called by each goroutine and passed in the shared buffer.
	F(Range)
	// A Body must be able to copy itself -- each goroutine is given a copy
	// of the body.
	Copy() Body
}
