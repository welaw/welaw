package vcs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
	apiv1 "github.com/welaw/welaw/api/v1"
)

func TestCreateLaw(t *testing.T) {
	// load env
	//err := godotenv.Load(testEnvPath)
	//require.NoError(t, err)

	dir, err := ioutil.TempDir("", "testrepo")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	logger := log.NewLogfmtLogger(os.Stderr)
	vc, err := NewVcs(dir, logger)
	require.NoError(t, err)

	hash, err := vc.CreateLaw(&apiv1.LawSet{
		Law: &apiv1.Law{
			Upstream: "test upstream",
			Ident:    "test ident",
			Title:    "test title",
		},
		Branch: &apiv1.Branch{
			Name: "testbranch",
		},
		Version: &apiv1.Version{
			Body: "test body",
			Msg:  "test msg",
		},
		Author: &apiv1.Author{
			Username: "testuser",
			Email:    "test@welaw.org",
		},
	})

	require.NoError(t, err)
	require.NotNil(t, hash)
}

//func TestCreateBranch(t *testing.T) {
//}

//func TestCreateVersion(t *testing.T) {
//}
