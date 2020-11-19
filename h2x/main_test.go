package main

import (
	"testing"
)

func BenchmarkX2dm(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x2dm("1A5C")
		x2dm("1A5C")
		x2dm("1A5C")
	}
}

func BenchmarkX2dmod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x2dmod("1A5C")
		x2dmod("1A5C")
		x2dmod("1A5C")
	}
}

//func BenchmarkTestWriteNameValuePairs(b *testing.B) {
//	m := make(map[string]string)
//	m["foo"] = "bar"
//	m["bad"] = "good"
//	for i := 0; i < b.N; i++ {
//		_ = writeNameValuePairs(m)
//
//		//fmt.Printf("%d bytes written\n", len(b))
//	}
//}
//
//func writeNameValuePairs(val map[string]string) []byte {
//	var abc = make([]byte, 0, 16)
//
//	buf := bytes.NewBuffer(abc)
//	for k, v := range val {
//		if buf.Len() > 0 {
//			buf.WriteByte('&')
//		}
//		buf.WriteString(k)
//		buf.WriteByte('=')
//		buf.WriteString(v)
//	}
//	return buf.Bytes()
//}
//
