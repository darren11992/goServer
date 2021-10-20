
package httpserver

import (
	"os"
)
// tape is used to isolate the writing to files after "seeking" back to the start of them
// needed to isolate a fix bug with using seek
type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0) //New: basically empties a file
	t.file.Seek(0, 0)
	return t.file.Write(p)
}