package logman_test

import "moul.io/logman"

func Example() {
	writer, _ := logman.NewWriteCloser("./path/to/logdir/", "my-app")
	defer writer.Close()
	writer.Write([]byte("hello world!"))
}
