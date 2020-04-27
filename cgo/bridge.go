package cgo

/**
typedef void (*cb)(char* extra, char* data);
void callCb(cb callback, char* extra, char* arg) {
	callback(extra, arg);
}
*/
import "C"
