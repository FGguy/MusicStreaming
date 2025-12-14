package controller

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	consts "music-streaming/internal/consts"
	"music-streaming/internal/data"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type HttpTestCase struct {
	Name     string
	Req      string
	FormBody string
	Expected string
	Status   int
}

func assertGetRequest(t *testing.T, req string, expectedStatus int, expectedBody string) {
	resp, err := http.Get(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error reading body, got %v", err)
	}

	assert.Equal(t, expectedStatus, resp.StatusCode)
	assert.Equal(t, expectedBody, string(bodyBytes))
}

func assertPostRequest(t *testing.T, req string, formBody string, expectedStatus int, expectedBody string) {
	resp, err := http.Post(req, "application/x-www-form-urlencoded", strings.NewReader(formBody))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error reading body, got %v", err)
	}

	assert.Equal(t, expectedStatus, resp.StatusCode)
	assert.Equal(t, expectedBody, string(bodyBytes))
}

func TestSubsonicMiddleware(t *testing.T) {
	err := os.Chdir("../../..")
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	if err := godotenv.Load(".env"); err != nil {
		t.Fatal("Error loading .env file")
	}
	// Setup dependencies
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
	hashedPassword := md5.Sum([]byte(adminPassword + salt))
	hashedPasswordHex := hex.EncodeToString(hashedPassword[:])
	baseURL := ts.URL + "/rest/ping"

	testCases := []HttpTestCase{
		{
			Name:     "Should return an error for missing required parameter.",
			Req:      fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=%s", baseURL, adminName, hashedPasswordHex, salt, consts.SubsonicVersion),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="10" message="Required parameter is missing."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "Should return an error for invalid incompatible client and server versions (server must upgrade).",
			Req:      fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=2.0.0&c=Test", baseURL, adminName, hashedPasswordHex, salt),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="30" message="Incompatible Subsonic REST protocol version. Server must upgrade."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "Should return an error for invalid incompatible client and server versions (client must upgrade).",
			Req:      fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=0.0.0&c=Test", baseURL, adminName, hashedPasswordHex, salt),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="20" message="Incompatible Subsonic REST protocol version. Client must upgrade."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "Should return a empty subsonic-response xml element.",
			Req:      fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=%s&c=Test", baseURL, adminName, hashedPasswordHex, salt, consts.SubsonicVersion),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "Should return error for wrong password.",
			Req:      fmt.Sprintf("%s?u=%s&t=WrongPassword&s=%s&v=%s&c=Test", baseURL, adminName, salt, consts.SubsonicVersion),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="40" message="Wrong username or password."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
		{
			Name:     "Should return error for non-existing user.",
			Req:      fmt.Sprintf("%s?u=NonExistingUser&t=%s&s=%s&v=%s&c=Test", baseURL, hashedPasswordHex, salt, consts.SubsonicVersion),
			Expected: `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="0" message="A generic error."></error></subsonic-response>`,
			Status:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assertGetRequest(t, tc.Req, tc.Status, tc.Expected)
		})
	}
}
