package data_test

import (
	"errors"
	"testing"

	data "github.com/bashbunni/project-management/database"
	"github.com/bashbunni/project-management/database/dbconn"
	"github.com/bashbunni/project-management/database/models"
	"github.com/matryer/is"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"github.com/tauraamui/xerror"
)

type testReader struct {
	onReadCallback func()
	readData       []byte
	readError      error
}

func (t testReader) Read(b []byte) (int, error) {
	t.onReadCallback()
	n := copy(b, t.readData)
	return n, t.readError
}

type testPlainPromptReader struct {
	testUsername string
	testError    error
}

func (t testPlainPromptReader) ReadPlain(string) (string, error) {
	return t.testUsername, t.testError
}

type testPasswordPromptReader struct {
	testPassword string
	testError    error
}

func (t testPasswordPromptReader) ReadPassword(string) ([]byte, error) {
	return []byte(t.testPassword), t.testError
}

type multipleAttemptPasswordPromptReader struct {
	attemptCount, maxCalls int
	passwordsToAttempt     []string
	testError              error
}

func (t *multipleAttemptPasswordPromptReader) ReadPassword(string) ([]byte, error) {
	if t.attemptCount >= t.maxCalls {
		return nil, xerror.New("TESTING ERROR: multipleAttempts exceeds maximum call limit")
	}
	password := []byte(t.passwordsToAttempt[t.attemptCount])
	t.attemptCount++
	return password, t.testError
}

type DBSetupTestSuite struct {
	suite.Suite
	dbMock                    dbconn.MockGormWrapper
	resetOpenDBConn           func()
	resetFs                   func()
	resetUC                   func()
	resetPlainPromptReader    func()
	resetPasswordPromptReader func()
}

func (suite *DBSetupTestSuite) SetupSuite() {
	suite.resetOpenDBConn = data.OverloadOpenDBConnection(
		func(string) (dbconn.GormWrapper, error) {
			return suite.dbMock, nil
		},
	)
}

func (suite *DBSetupTestSuite) TearDownSuite() {
	suite.resetOpenDBConn()
}

func (suite *DBSetupTestSuite) SetupTest() {
	suite.dbMock = dbconn.Mock()
	suite.resetFs = data.OverloadFS(afero.NewMemMapFs())
	suite.resetUC = data.OverloadUC(func() (string, error) {
		return "/testroot/.cache", nil
	})
	suite.resetPlainPromptReader = data.OverloadPlainPromptReader(
		testPlainPromptReader{
			testUsername: "testadmin",
		},
	)

	suite.resetPasswordPromptReader = data.OverloadPasswordPromptReader(
		testPasswordPromptReader{
			testPassword: "testpassword",
		},
	)
}

func (suite *DBSetupTestSuite) TearDownTest() {
	suite.dbMock = nil
	suite.resetFs()
	suite.resetUC()
	suite.resetPlainPromptReader()
	suite.resetPasswordPromptReader()
}

func (suite *DBSetupTestSuite) TestCreateFullFilePathForDBWithSingleRootUserDir() {
	is := is.New(suite.T())
	is.NoErr(data.Setup())

	created := suite.dbMock.Created()
	is.Equal(len(created), 1)
	user := models.User{}
	is.NoErr(dbconn.Replace(&user, created[0]))
	is.Equal(user.Name, "testadmin")
}

func (suite *DBSetupTestSuite) TestConnectWithoutHavingToRunSetupFirst() {
	is := is.New(suite.T())
	is.NoErr(data.Setup())

	conn, err := data.Connect()
	is.NoErr(err)
	is.True(conn != nil)
}

func (suite *DBSetupTestSuite) TestCreateFileAndThenRemovedOnDestroy() {
	is := is.New(suite.T())
	is.NoErr(data.Setup())

	is.NoErr(data.Destroy())

	is.Equal(data.Destroy().Error(), "remove /testroot/.cache/tacusci/dragondaemon/dd.db: file does not exist")
}

func (suite *DBSetupTestSuite) TestReturnErrorFromSetupDueToROFileSystem() {
	is := is.New(suite.T())
	suite.resetFs = data.OverloadFS(afero.NewReadOnlyFs(afero.NewMemMapFs()))
	is.Equal(data.Setup().Error(), "unable to create database file: operation not permitted")
}

func (suite *DBSetupTestSuite) TestReturnErrorFromSetupDueToDBAlreadyExisting() {
	is := is.New(suite.T())
	is.NoErr(data.Setup())
	is.Equal(data.Setup().Error(), "database file already exists: /testroot/.cache/tacusci/dragondaemon/dd.db")
}

func (suite *DBSetupTestSuite) TestReturnErrorFromSetupDueToPathResolutionFailure() {
	is := is.New(suite.T())
	suite.resetUC = data.OverloadUC(func() (string, error) {
		return "", xerror.New("test cache dir error")
	})
	is.Equal(data.Setup().Error(), "unable to resolve dd.db database file location: test cache dir error")
}

func (suite *DBSetupTestSuite) TestSetupUnableToConnectDueToUCError() {
	is := is.New(suite.T())
	suite.resetUC = data.OverloadUC(func() (string, error) {
		return "", xerror.New("test cache dir error")
	})
	conn, err := data.Connect()
	is.True(err != nil && conn == nil)
	is.Equal(err.Error(), "unable to resolve dd.db database file location: test cache dir error")
}

func (suite *DBSetupTestSuite) TestSetupUnableToConnectAfterFileCreate() {
	is := is.New(suite.T())
	callCount := 0
	suite.resetOpenDBConn = data.OverloadOpenDBConnection(func(s string) (dbconn.GormWrapper, error) {
		callCount++
		if callCount >= 1 {
			return nil, errors.New("test open db conn error")
		}
		return dbconn.Mock(), nil
	})
	defer suite.resetOpenDBConn()
	is.Equal(data.Setup().Error(), "unable to open db connection: test open db conn error")
}

func (suite *DBSetupTestSuite) TestUnableToResolveDBPathHandlesAndReturnsWrappedError() {
	is := is.New(suite.T())
	is.NoErr(data.Setup())

	suite.resetUC = data.OverloadUC(func() (string, error) {
		return "", xerror.New("test cache dir error")
	})

	is.Equal(
		data.Destroy().Error(), "unable to delete database file: unable to resolve dd.db database file location: test cache dir error",
	)
}

func (suite *DBSetupTestSuite) TestSetupHandlesCreateRootUserErrorAndReturnsWrappedError() {
	suite.resetOpenDBConn = data.OverloadOpenDBConnection(
		func(string) (dbconn.GormWrapper, error) {
			return dbconn.Mock().SetError(errors.New("test create failed")), nil
		},
	)
	defer suite.resetOpenDBConn()

	is := is.New(suite.T())
	is.Equal(data.Setup().Error(), "unable to create root user entry: test create failed")
}

func (suite *DBSetupTestSuite) TestUsernamePromptErrorHandlesAndReturnWrappedError() {
	is := is.New(suite.T())
	resetPlainPromptReader := data.OverloadPlainPromptReader(
		testPlainPromptReader{
			testError: xerror.New("testing read username error"),
		},
	)
	defer resetPlainPromptReader()

	resetPasswordPromptReader := data.OverloadPasswordPromptReader(
		testPasswordPromptReader{
			testPassword: "testpassword",
		},
	)
	defer resetPasswordPromptReader()

	is.Equal(data.Setup().Error(), "failed to prompt for root username: testing read username error")
}

func (suite *DBSetupTestSuite) TestSetupReturnsErrorFromTooManyPasswordAttempts() {
	is := is.New(suite.T())
	resetPlainPromptReader := data.OverloadPlainPromptReader(
		testPlainPromptReader{
			testUsername: "testadmin",
		},
	)
	defer resetPlainPromptReader()

	resetPasswordPromptReader := data.OverloadPasswordPromptReader(
		&multipleAttemptPasswordPromptReader{
			maxCalls: 6,
			passwordsToAttempt: []string{
				"1stpair", "1stpairnomatch", "2ndpair", "2ndpairnomatch", "3rdpair", "3rdpairnomatch",
			},
		},
	)
	defer resetPasswordPromptReader()

	is.Equal(
		data.Setup().Error(),
		"failed to prompt for root password: tried entering new password at least 3 times",
	)
}

func TestDBSetupTestSuite(t *testing.T) {
	suite.Run(t, &DBSetupTestSuite{})
}

func TestOpenDBConnection(t *testing.T) {
	is := is.New(t)
	conn, err := data.OpenDBConnection("file::memory:?cache=shared")
	is.NoErr(err)
	is.True(conn != nil)
}

func TestPlainPromptReaderShouldReadFromReadableAndReturnValue(t *testing.T) {
	is := is.New(t)
	calledCount := 0
	plainReader := data.NewStdinPlainReader(
		testReader{
			readData: []byte("testuser\n"),
			onReadCallback: func() {
				calledCount++
			},
		},
	)

	value, err := plainReader.ReadPlain("")
	is.NoErr(err)
	is.Equal(value, "testuser")
	is.Equal(calledCount, 1)
}

func TestConnectHandlesAutoMigrateErrorAndReturnsWrappedError(t *testing.T) {
	resetOpenDBConn := data.OverloadOpenDBConnection(
		func(string) (dbconn.GormWrapper, error) {
			return dbconn.Mock().SetAutoMigrateError(errors.New("test automigrate failed")), nil
		},
	)
	defer resetOpenDBConn()

	is := is.New(t)
	conn, err := data.Connect()
	is.True(err != nil && conn == nil)
	is.Equal(err.Error(), "unable to run automigrations: test automigrate failed")
}
