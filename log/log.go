package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

func print(tag string, format string, args []interface{}) {

	//time.
	//
	//format specifiers:
	//|
	//| mm:dd HH:MM:SECOND    yy zone
	//|	01/02 03:04:05.000 PM 06 0700
	//|       15           AM
	//
	var now string = time.Now().Format("06-01-02 15:04:05.000")

	//caller.
	//| Caller(0): print();
	//| Caller(1): I() or E();
	//| Caller(2): the user function.
	var _, path, line, okay = runtime.Caller(2)

	//message.
	var msg string = fmt.Sprintf(format, args...)

	if okay {
		_, file := filepath.Split(path)
		fmt.Printf("%s|%s|%s(%04d)|%s\n", now, tag, file, line, msg)
	} else {
		fmt.Printf("%s|%s|?|%s\n", now, tag, msg)
	}
}

func I(format string, args ...interface{}) {
	print("I", format, args)
}

func E(format string, args ...interface{}) {
	print("E", format, args)
}
