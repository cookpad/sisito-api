package sisito

import (
	"github.com/bouk/monkey"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

func tempFile(content string, callback func(f *os.File)) {
	tmpfile, _ := ioutil.TempFile("", "sisito")
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString(content)
	tmpfile.Sync()
	tmpfile.Seek(0, 0)
	callback(tmpfile)
}

func readResponse(res *http.Response) (string, int) {
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	return string(content), res.StatusCode
}

func patchInstanceMethod(receiver interface{}, methodName string, replacementf func(**monkey.PatchGuard) interface{}) {
	var guard *monkey.PatchGuard
	replacement := replacementf(&guard)
	guard = monkey.PatchInstanceMethod(
		reflect.TypeOf(receiver), methodName, replacement)
}
