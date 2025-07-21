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

	SqlSetup(pg_pool)

	return pg_pool, cache, nil
}

func TestSubsonicMiddleware(t *testing.T) {
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

	t.Run("Should return an error for missing required parameter.", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/rest/ping?u=%s&t=%s&s=%s&v=%s", ts.URL, adminName, hex.EncodeToString(hashedPassword[:]), salt, subsonic.SubsonicVersion))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		bodyString, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"failed\" version=\"1.16.1\"><error code=\"10\" message=\"Required parameter is missing.\"/></subsonic-response>", string(bodyString))
	})

	t.Run("Should return an error for invalid incompatible client and server versions.", func(t *testing.T) {
		//Server upgrade
		resp, err := http.Get(fmt.Sprintf("%s/rest/ping?u=%s&t=%s&s=%s&v=%s&c=%s", ts.URL, adminName, hex.EncodeToString(hashedPassword[:]), salt, "2.0.0", "Test"))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		bodyString, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"failed\" version=\"1.16.1\"><error code=\"30\" message=\"Incompatible Subsonic REST protocol version. Server must upgrade.\"/></subsonic-response>", string(bodyString))

		//Client upgrade
		resp, err = http.Get(fmt.Sprintf("%s/rest/ping?u=%s&t=%s&s=%s&v=%s&c=%s", ts.URL, adminName, hex.EncodeToString(hashedPassword[:]), salt, "0.0.0", "Test"))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		bodyString, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"failed\" version=\"1.16.1\"><error code=\"20\" message=\"Incompatible Subsonic REST protocol version. Client must upgrade.\"/></subsonic-response>", string(bodyString))
	})

	t.Run("Should return a empty subsonic-response xml element.", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/rest/ping?u=%s&t=%s&s=%s&v=%s&c=%s", ts.URL, adminName, hex.EncodeToString(hashedPassword[:]), salt, subsonic.SubsonicVersion, "Test"))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		bodyString, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"ok\" version=\"1.16.1\"></subsonic-response>", string(bodyString))
	})

	t.Run("Should return error for wrong username or password.", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/rest/ping?u=%s&t=%s&s=%s&v=%s&c=%s", ts.URL, adminName, "WrongPassword", salt, subsonic.SubsonicVersion, "Test"))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		bodyString, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"failed\" version=\"1.16.1\"><error code=\"40\" message=\"Wrong username or password.\"/></subsonic-response>", string(bodyString))

		resp, err = http.Get(fmt.Sprintf("%s/rest/ping?u=%s&t=%s&s=%s&v=%s&c=%s", ts.URL, "NonExistingUser", hex.EncodeToString(hashedPassword[:]), salt, subsonic.SubsonicVersion, "Test"))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		bodyString, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "<subsonic-response xmlns=\"http://subsonic.org/restapi\" status=\"failed\" version=\"1.16.1\"><error code=\"40\" message=\"Wrong username or password.\"/></subsonic-response>", string(bodyString))
	})
}
