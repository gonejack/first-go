package util

type TicketBucket struct {
	bucket chan int
}

func (tb *TicketBucket) Put() {
	select {
	case tb.bucket <- 0:
	default:
	}
}

func (tb *TicketBucket) Get() int {
	return <-tb.bucket
}

func NewTicketBucket(size int) (tb *TicketBucket) {
	tb = &TicketBucket{
		bucket: make(chan int, size),
	}

	for range IdxRange(size) {
		tb.Put()
	}

	return
}
