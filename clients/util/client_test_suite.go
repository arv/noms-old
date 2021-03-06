package util

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/attic-labs/noms/d"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	LdbFlagName string
	TempDir     string
	LdbDir      string
	out         *os.File
}

func (suite *ClientTestSuite) SetupSuite() {
	dir, err := ioutil.TempDir(os.TempDir(), "nomstest")
	d.Chk.NoError(err)
	out, err := ioutil.TempFile(dir, "out")
	d.Chk.NoError(err)

	suite.TempDir = dir
	suite.LdbDir = path.Join(dir, "ldb")
	suite.out = out
}

func (suite *ClientTestSuite) TearDownSuite() {
	suite.out.Close()
	defer d.Chk.NoError(os.RemoveAll(suite.TempDir))
}

func (suite *ClientTestSuite) Run(m func(), args []string) string {
	origArgs := os.Args
	origOut := os.Stdout
	origErr := os.Stderr

	ldbFlagName := suite.LdbFlagName

	if ldbFlagName == "" {
		ldbFlagName = "-ldb"
	}

	os.Args = append([]string{"cmd", ldbFlagName, suite.LdbDir}, args...)
	os.Stdout = suite.out

	defer func() {
		os.Args = origArgs
		os.Stdout = origOut
		os.Stderr = origErr
	}()

	m()

	_, err := suite.out.Seek(0, 0)
	d.Chk.NoError(err)
	b, err := ioutil.ReadAll(os.Stdout)
	d.Chk.NoError(err)

	_, err = suite.out.Seek(0, 0)
	d.Chk.NoError(err)
	err = suite.out.Truncate(0)
	d.Chk.NoError(err)

	return string(b)
}
