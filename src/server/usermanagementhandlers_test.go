package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	consts "music-streaming/consts"
	sqlc "music-streaming/sql/sqlc"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestUserManagementHandlers(t *testing.T) {
	pg_pool, cache, err := getServerDependencies(t)
	if err != nil {
		t.Fatalf("Failed to create dependencies for server in TestSystemHandlers, Error:%v", err)
	}
	defer pg_pool.Close()

	server := NewServer(pg_pool, cache)
	ts := httptest.NewServer(server.router)
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
	conn, err := pg_pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("Failed to acquire DB connection")
	}
	defer conn.Release()
	q := sqlc.New(conn)

	testUser1 := sqlc.CreateDefaultUserParams{
		Username: pgtype.Text{String: "test1", Valid: true},
		Password: "test1",
		Email:    "test1",
	}
	testUser2 := sqlc.CreateDefaultUserParams{
		Username: pgtype.Text{String: "test2", Valid: true},
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

	testUser1HashedPassword := md5.Sum([]byte(testUser1.Password + salt))
	testUser1HashedPasswordHex := hex.EncodeToString(testUser1HashedPassword[:])
	testUser2HashedPassword := md5.Sum([]byte(testUser2.Password + salt))
	testUser2HashedPasswordHex := hex.EncodeToString(testUser2HashedPassword[:])
	version := consts.SubsonicVersion
	expectedResponse := `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`

	t.Run("/rest/getUser route", func(t *testing.T) {
		// Admin accessing user
		reqAdmin := fmt.Sprintf("%s/rest/getUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, adminName, adminHashedPasswordHex, salt, version, testUser1.Username.String)
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"><user username=\"test1\" email=\"test1\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"false\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user></subsonic-response>"
		assertGetRequest(t, reqAdmin, 200, expectedResponse)

		// User accessing self
		reqSelf := fmt.Sprintf("%s/rest/getUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, testUser1.Username.String)
		assertGetRequest(t, reqSelf, 200, expectedResponse)

		// User accessing another user
		reqOther := fmt.Sprintf("%s/rest/getUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, testUser2.Username.String)
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertGetRequest(t, reqOther, 200, expectedResponse)
	})

	t.Run("/rest/getUsers route", func(t *testing.T) {
		//Admin
		reqAdmin := fmt.Sprintf("%s/rest/getUsers?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, adminName, adminHashedPasswordHex, salt, version)
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"><user username=\"default\" email=\"default\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"true\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user><user username=\"test1\" email=\"test1\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"false\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user><user username=\"test2\" email=\"test2\" scrobblingEnabled=\"false\" ldapAuthenticated=\"false\" adminRole=\"false\" settingsRole=\"true\" streamRole=\"true\" jukeboxRole=\"false\" downloadRole=\"false\" uploadRole=\"false\" playlistRole=\"false\" coverArtRole=\"false\" commentRole=\"false\" podcastRole=\"false\" shareRole=\"false\" videoConversionRole=\"false\" maxBitRate=\"0\"></user></subsonic-response>"
		assertGetRequest(t, reqAdmin, 200, expectedResponse)
		//Not Admin
		reqOther := fmt.Sprintf("%s/rest/getUsers?u=%s&t=%s&s=%s&v=%s&c=Test", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version)
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertGetRequest(t, reqOther, 200, expectedResponse)
	})

	t.Run("/rest/createUser route", func(t *testing.T) {
		//Admin no params
		reqAdmin := fmt.Sprintf("%s/rest/createUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&password=%s&email=%s", ts.URL, adminName, adminHashedPasswordHex, salt, version, "test3", "test3", "test3")
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"></subsonic-response>"
		assertPostRequest(t, reqAdmin, 200, expectedResponse)
		//Admin random set of params
		reqAdmin = fmt.Sprintf("%s/rest/createUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&password=%s&email=%s&uploadRole=true&podcastRole=true", ts.URL, adminName, adminHashedPasswordHex, salt, version, "test4", "test4", "test4")
		assertPostRequest(t, reqAdmin, 200, expectedResponse)
		//Not Admin
		reqOther := fmt.Sprintf("%s/rest/createUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&password=%s&email=%s", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, "test5", "test5", "test5")
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertPostRequest(t, reqOther, 200, expectedResponse)
	})

	t.Run("/rest/updateUser route", func(t *testing.T) {
		//Admin random set of params
		reqAdmin := fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&settingsRole=false", ts.URL, adminName, adminHashedPasswordHex, salt, version, "test2")
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"></subsonic-response>"
		assertPostRequest(t, reqAdmin, 200, expectedResponse)
		//User == user & set roles
		reqSelf := fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&shareRole=true", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, "test1")
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"></subsonic-response>"
		assertPostRequest(t, reqSelf, 200, expectedResponse)
		//user == user & not set roles
		reqSelf = fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&shareRole=true", ts.URL, testUser2.Username.String, testUser2HashedPasswordHex, salt, version, "test2")
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertPostRequest(t, reqSelf, 200, expectedResponse)
		//User != user & not set roles
		reqSelf = fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&shareRole=true", ts.URL, testUser2.Username.String, testUser2HashedPasswordHex, salt, version, "test1")
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertPostRequest(t, reqSelf, 200, expectedResponse)
		//user != user & set roles
		reqSelf = fmt.Sprintf("%s/rest/updateUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&shareRole=true", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, "test2")
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertPostRequest(t, reqSelf, 200, expectedResponse)
	})

	t.Run("/rest/changePassword route", func(t *testing.T) {
		//Admin
		reqAdmin := fmt.Sprintf("%s/rest/changePassword?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&password=%s", ts.URL, adminName, adminHashedPasswordHex, salt, version, testUser1.Username.String, testUser1.Password)
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"></subsonic-response>"
		assertPostRequest(t, reqAdmin, 200, expectedResponse)
		//User == User
		reqSelf := fmt.Sprintf("%s/rest/changePassword?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&password=%s", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, testUser1.Username.String, testUser1.Password)
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"></subsonic-response>"
		assertPostRequest(t, reqSelf, 200, expectedResponse)
		//User != User
		reqOther := fmt.Sprintf("%s/rest/changePassword?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s&password=%s", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, testUser2.Username.String, "foo")
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertPostRequest(t, reqOther, 200, expectedResponse)
	})

	t.Run("/rest/deleteUser route", func(t *testing.T) {
		//Not Admin
		reqOther := fmt.Sprintf("%s/rest/deleteUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, testUser1.Username.String, testUser1HashedPasswordHex, salt, version, testUser1.Username.String)
		expectedResponse = `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="50" message="User is not authorized for the given operation."/></subsonic-response>`
		assertPostRequest(t, reqOther, 200, expectedResponse)

		//Admin
		reqAdmin := fmt.Sprintf("%s/rest/deleteUser?u=%s&t=%s&s=%s&v=%s&c=Test&username=%s", ts.URL, adminName, adminHashedPasswordHex, salt, version, testUser1.Username.String)
		expectedResponse = "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"></subsonic-response>"
		assertPostRequest(t, reqAdmin, 200, expectedResponse)
	})
}
