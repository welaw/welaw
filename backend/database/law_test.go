package database_test

import (
	"os"

	"github.com/go-kit/kit/log"
)

const (
	testEnvPath = "../../.env"
)

func newTestLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}

//func TestListBranchVersions(t *testing.T) {
//err := godotenv.Load(testEnvPath)
//require.NoError(t, err)

//logger := newTestLogger()
//db := backend.NewTestDatabase(logger)

//sets, err := db.ListBranchVersions("us", "114s3084", "master")
//require.NoError(t, err)

//fmt.Printf("result count: %v\n", len(sets))
//for _, set := range sets {
//fmt.Printf("set: %+v\n", set)
//}

//}

//func TestListLawBranches(t *testing.T) {
//err := godotenv.Load(testEnvPath)
//require.NoError(t, err)

//logger := newTestLogger()
//db := backend.NewTestDatabase(logger)

//sets, err := db.ListLawBranches("us", "114s3084")
//require.NoError(t, err)

//fmt.Printf("result count: %v\n", len(sets))
//for _, set := range sets {
//fmt.Printf("set: %+v\n", set)
//}

//}
