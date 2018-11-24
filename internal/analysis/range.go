package analysis

type Range struct {
	End    int
	Begin  int
	stride int
	next   int
}

func (l *Range) Next() (n int) {
	n = l.next
	l.next += l.stride
	return
}
