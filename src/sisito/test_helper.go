package sisito

import (
	"io/ioutil"
	"os"
)

func tempFile(content string, callback func(f *os.File)) {
	tmpfile, _ := ioutil.TempFile("", "sisito")
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString(content)
	tmpfile.Sync()
	tmpfile.Seek(0, 0)
	callback(tmpfile)
}
