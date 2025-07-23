package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"music-streaming/util/subsonic"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func getServerDependencies(t *testing.T) (pg_pool *pgxpool.Pool, cache *redis.Client, err error) {
	err = godotenv.Load()
	if err != nil {
		return nil, nil, err
	}

	//remember to defer closing outside of call
	postgresString, ok := os.LookupEnv(("TEST_POSTGRES_CONNECTION_STRING"))
	if !ok {
		t.Fatalf("Failed to lookup postgres connection string")
	}
	pg_pool, err = pgxpool.New(context.Background(), postgresString)
	if err != nil {
		return nil, nil, err
	}

	redisString, ok := os.LookupEnv(("TEST_REDIS_CONNECTION_STRING"))
	if !ok {
		t.Fatalf("Failed to lookup redis connection string")
	}
	opt, err := redis.ParseURL(redisString)
	if err != nil {
		return nil, nil, err
	}
	cache = redis.NewClient(opt)

	SqlSetup(pg_pool, true)

	return pg_pool, cache, nil
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

func assertPostRequest(t *testing.T, req string, expectedStatus int, expectedBody string) {
	resp, err := http.Post(req, "application/x-www-form-urlencoded", nil)
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
	// Setup dependencies
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
	hashedPassword := md5.Sum([]byte(adminPassword + salt))
	hashedPasswordHex := hex.EncodeToString(hashedPassword[:])
	baseURL := ts.URL + "/rest/ping"

	t.Run("Should return an error for missing required parameter.", func(t *testing.T) {
		req := fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=%s", baseURL, adminName, hashedPasswordHex, salt, subsonic.SubsonicVersion)
		expected := `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="10" message="Required parameter is missing."/></subsonic-response>`
		assertGetRequest(t, req, 200, expected)
	})

	t.Run("Should return an error for invalid incompatible client and server versions.", func(t *testing.T) {
		req1 := fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=2.0.0&c=Test", baseURL, adminName, hashedPasswordHex, salt)
		expected1 := `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="30" message="Incompatible Subsonic REST protocol version. Server must upgrade."/></subsonic-response>`
		assertGetRequest(t, req1, 200, expected1)

		req2 := fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=0.0.0&c=Test", baseURL, adminName, hashedPasswordHex, salt)
		expected2 := `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="20" message="Incompatible Subsonic REST protocol version. Client must upgrade."/></subsonic-response>`
		assertGetRequest(t, req2, 200, expected2)
	})

	t.Run("Should return a empty subsonic-response xml element.", func(t *testing.T) {
		req := fmt.Sprintf("%s?u=%s&t=%s&s=%s&v=%s&c=Test", baseURL, adminName, hashedPasswordHex, salt, subsonic.SubsonicVersion)
		expected := `<subsonic-response xmlns="http://subsonic.org/restapi" status="ok" version="1.16.1"></subsonic-response>`
		assertGetRequest(t, req, 200, expected)
	})

	t.Run("Should return error for wrong username or password.", func(t *testing.T) {
		req1 := fmt.Sprintf("%s?u=%s&t=WrongPassword&s=%s&v=%s&c=Test", baseURL, adminName, salt, subsonic.SubsonicVersion)
		expected := `<subsonic-response xmlns="http://subsonic.org/restapi" status="failed" version="1.16.1"><error code="40" message="Wrong username or password."/></subsonic-response>`
		assertGetRequest(t, req1, 200, expected)

		req2 := fmt.Sprintf("%s?u=NonExistingUser&t=%s&s=%s&v=%s&c=Test", baseURL, hashedPasswordHex, salt, subsonic.SubsonicVersion)
		assertGetRequest(t, req2, 200, expected)
	})
}
