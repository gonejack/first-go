// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package genny

// StringQueue is a queue of Strings.
type StringQueue struct {
	items []string
}

func NewStringQueue() *StringQueue {
	return &StringQueue{items: make([]string, 0)}
}
func (q *StringQueue) Push(item string) {
	q.items = append(q.items, item)
}
func (q *StringQueue) Pop() string {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}