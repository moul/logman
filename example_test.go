package logman_test

import (
	"fmt"

	"moul.io/logman"
)

func Example() {
	// new log manager
	manager := logman.Manager{
		Path:     "./path/to/dir",
		MaxFiles: 10,
	}

	/*
		// cleanup old log files for a specific app name
		err := manager.GCWithName("my-app")
		checkErr(err)

		// cleanup old log files for any app sharing this log directory
		err = manager.GC()
		checkErr(err)
	*/

	// list existing log files
	files, err := manager.Files()
	checkErr(err)
	fmt.Println(files)

	// - create an WriteCloser
	// - automatically delete old log files if it hits a limit
	writer, err := manager.New("my-app")
	checkErr(err)
	defer writer.Close()
	writer.Write([]byte("hello world!\n"))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
