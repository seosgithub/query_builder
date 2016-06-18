package query_builder

import (
	"os"
	"path"

	"github.com/jmcvetta/randutil"
)

// A place to put a temp file
func TempNonExistantFilePath() string {
	fn, err := randutil.AlphaString(10)
	if err != nil {
		panic(err)
	}

	return path.Join(os.TempDir(), fn)
}
