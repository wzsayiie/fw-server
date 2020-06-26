package log

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func print(tag string, format string, args []interface{}) {

	var msg string = fmt.Sprintf(format, args...)

	//Caller(0): print();
	//Caller(1): I() or E();
	//Caller(2): the user function.
	_, path, line, okay := runtime.Caller(2)

	if okay {
		_, file := filepath.Split(path)
		fmt.Printf("%s|%s(%04d)|%s\n", tag, file, line, msg)
	} else {
		fmt.Printf("%s|?|%s\n", tag, msg)
	}
}

func I(format string, args ...interface{}) {
	print("I", format, args)
}

func E(format string, args ...interface{}) {
	print("E", format, args)
}
