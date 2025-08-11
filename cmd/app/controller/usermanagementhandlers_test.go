package controller

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	consts "music-streaming/internal/consts"
	"music-streaming/internal/data"
	sqlc "music-streaming/internal/sql/sqlc"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestUserManagementHandlers(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		t.Fatal("Error loading .env file")
	}

	dataLayer, err := data.NewTest(context.Background())
	if err != nil {
		t.Fatalf("Failed initializing data layer. Error: %s", err)
	}
	defer dataLayer.Pg_pool.Close()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed loading server configuration file. Error: %s", err)
	}
	app := NewApplication(dataLayer, config)

	ts := httptest.NewServer(app.Router)
	defer ts.Close()

	adminName, adminNameDefined := os.LookupEnv("ADMIN_NAME")
	adminPassword, adminPasswordDefined := os.LookupEnv("ADMIN_PASSWORD")
	salt := "abcdef"
	if !adminNameDefined || !adminPasswordDefined {
		t.Fatalf("Failed to lookup admin credentials from environment")
	}
	adminHashedPassword := md5.Sum([]byte(adminPassword + salt))
	adminHashedPasswordHex := hex.EncodeToString(adminHashedPassword[:])

	// Create test users
	ctx := context.Background()
	conn, err := dataLayer.Pg_pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Failed to acquire DB connection")
	}
	defer conn.Release()
	q := sqlc.New(conn)

	testUser1 := sqlc.CreateDefaultUserParams{
		Username: "test1",
		Password: "test1",
		Email:    "test1",
	}
	testUser2 := sqlc.CreateDefaultUserParams{
		Username: "test2",
		Password: "test2",
		Email:    "test2",
	}

	_, err = q.CreateDefaultUser(ctx, testUser1)
	if err != nil {
		t.Fatalf("Failed to create test user1: %v", err)
	}
	_, err = q.CreateDefaultUser(ctx, testUser2)
	if err != nil {
		t.Fatalf("Failed to create test user2: %v", err)
	}

	var (
		testUser1HashedPassword    = md5.Sum([]byte(testUser1.Password + salt))
		testUser1HashedPasswordHex = hex.EncodeToString(testUser1HashedPassword[:])
		testUser2HashedPassword    = md5.Sum([]byte(testUser2.Password + salt))
		testUser2HashedPasswordHex = hex.EncodeToString(testUser2HashedPassword[:])
		version                    = consts.SubsonicVersion
	)

	testCases := []HttpTestCase{
		// getUser route tests
		{
			Name:     "/rest/getUser - Admin accessing user",
			Req:      fmt.Sprintf("%s/rest/getUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, adminName, adminHashedPasswordHex, salt, version, testUser1.Username),
			Expected: "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"><user username=\"test1\" email=\"test1\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"false\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user></subsonic-response>",
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/getUser - User accessing self",
			Req:      fmt.Sprintf("%s/rest/getUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version, testUser1.Username),
			Expected: "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"><user username=\"test1\" email=\"test1\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"false\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user></subsonic-response>",
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/getUser - User accessing another user (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/getUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version, testUser2.Username),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		// getUsers route tests
		{
			Name:     "/rest/getUsers - Admin accessing all users",
			Req:      fmt.Sprintf("%s/rest/getUsers?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, adminName, adminHashedPasswordHex, salt, version),
			Expected: "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"><user username=\"default\" email=\"default\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"true\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user><user username=\"test1\" email=\"test1\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"false\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user><user username=\"test2\" email=\"test2\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"false\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user></subsonic-response>",
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/getUsers - Non-admin user accessing all users (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/getUsers?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assertGetRequest(t, tc.Req, tc.Status, tc.Expected)
		})
	}

	testCases = []HttpTestCase{
		// createUser route tests
		{
			Name:     "/rest/createUser - Admin creating user with basic params",
			Req:      fmt.Sprintf("%s/rest/createUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, adminName, adminHashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s&password=%s&email=%s", "test3", "test3", "test3"),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/createUser - Admin creating user with additional roles",
			Req:      fmt.Sprintf("%s/rest/createUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, adminName, adminHashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s&password=%s&email=%s&uploadRole=true&podcastRole=true", "test4", "test4", "test4"),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/createUser - Non-admin user creating user (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/createUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s&password=%s&email=%s", "test5", "test5", "test5"),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},

		// updateUser route tests
		{
			Name:     "/rest/updateUser - Admin updating another user's roles",
			Req:      fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, adminName, adminHashedPasswordHex, salt, version),
			FormBody: "username=test2&settingsRole=false",
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/updateUser - User updating self with settings role",
			Req:      fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version),
			FormBody: "username=test1&coverArtRole=true",
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/updateUser - User updating self without settings role (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser2.Username, testUser2HashedPasswordHex, salt, version),
			FormBody: "username=test2&email=test@example.com",
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/updateUser - User updating another user without settings role (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser2.Username, testUser2HashedPasswordHex, salt, version),
			FormBody: "username=test1&email=shouldfail@example.com",
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/updateUser - User updating another user with settings role (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version),
			FormBody: "username=test2&uploadRole=true",
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},

		// changePassword route tests
		{
			Name:     "/rest/changePassword - Admin changing user password",
			Req:      fmt.Sprintf("%s/rest/changePassword?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, adminName, adminHashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s&password=%s", testUser1.Username, testUser1.Password),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/changePassword - User changing own password",
			Req:      fmt.Sprintf("%s/rest/changePassword?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s&password=%s", testUser1.Username, testUser1.Password),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/changePassword - User changing another user's password (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/changePassword?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s&password=%s", testUser2.Username, "foo"),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},

		// deleteUser route tests
		{
			Name:     "/rest/deleteUser - Non-admin user deleting user (unauthorized)",
			Req:      fmt.Sprintf("%s/rest/deleteUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username, testUser1HashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s", testUser1.Username),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "/rest/deleteUser - Admin deleting user",
			Req:      fmt.Sprintf("%s/rest/deleteUser?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, adminName, adminHashedPasswordHex, salt, version),
			FormBody: fmt.Sprintf("username=%s", testUser1.Username),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assertPostRequest(t, tc.Req, tc.FormBody, tc.Status, tc.Expected)
		})
	}
}
