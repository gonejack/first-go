package main

import (
	"fmt"
	"unicode/utf8"
)

type treeLabel struct {
	label string
	value string
	subs  []*treeLabel
}

func (t *treeLabel) leftPadToLen(str string, target int) string {
	for t.displayLen(str) < target {
		str = " " + str
	}
	return str
}
func (t *treeLabel) displayLen(str string) (n int) {
	for _, r := range []rune(str) {
		if utf8.RuneLen(r) > 1 {
			n += 2
		} else {
			n += 1
		}
	}
	return
}
func (t *treeLabel) printL(labelLen int) (text string) {
	if t.label != "" {
		label := t.leftPadToLen(t.label, labelLen)
		text += fmt.Sprintf("%s:%s\n", label, t.value)
	}

	max := 0
	for idx := range t.subs {
		if max < t.displayLen(t.subs[idx].label) {
			max = t.displayLen(t.subs[idx].label)
		}
	}
	next := labelLen + max
	for next > max && next > labelLen+2 {
		next -= 1
	}
	for i := range t.subs {
		text += t.subs[i].printL(next)
	}

	return
}
func (t *treeLabel) print() string {
	return t.printL(0)
}
func (t *treeLabel) append(label, value string) (sub *treeLabel) {
	sub = &treeLabel{label: label, value: value}
	t.subs = append(t.subs, sub)
	return
}
func (t *treeLabel) group(group string) (sub *treeLabel) {
	return t.append(group, "")
}

func main() {
	var tr treeLabel

	abc := tr.group("abc")
	{
		abc.append("a1", "v1")
		abc.append("a2", "v2")
	}
	def := tr.group("def")
	{
		def.append("d1", "d1")
		def.append("d2", "d2")
	}

	print(tr.print())
}
