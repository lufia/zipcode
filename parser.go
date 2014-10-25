package postal

type Parser interface {
	Parse(c <-chan interface{}, c1 chan<- interface{})
}

type garbageStrip func(entry *Entry) *Entry

func (f garbageStrip) Parse(c <-chan interface{}, c1 chan<- interface{}) {
	for v := range c {
		if err, ok := v.(error); ok {
			c1 <- err
			break
		}
		c1 <- f(v.(*Entry))
	}
}
