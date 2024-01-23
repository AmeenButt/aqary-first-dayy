package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"assesment.sqlc.dev/app/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	// Replace with your actual package path
	// Replace with your actual package path
)

func TestRegisterUserRoutes(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:admin@localhost:5432/test_sqlc_practice")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	ctx := context.Background()
	router := gin.Default()
	queries := postgres.NewStore(conn)
	RegisterUserRoutes(router, conn, &ctx, queries)

	// You can add assertions here to check if routes are correctly registered
	assertRouteExists(t, router, "POST", "/users/create")
	assertRouteExists(t, router, "POST", "/users/sign-in")
	assertRouteExists(t, router, "GET", "/users/get")
	assertRouteExists(t, router, "GET", "/users/get-all-users")
	assertRouteExists(t, router, "POST", "/users/upload-profile")
}

func assertRouteExists(t *testing.T, router *gin.Engine, method, path string) {
	routes := router.Routes()
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			return // Route found, test passed
		}
	}

	t.Errorf("Route %s %s not found", method, path)
}

func TestUserRoutes(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:admin@localhost:5432/test_sqlc_practice")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	queries := postgres.NewStore(conn)
	ctx := context.Background()
	router := gin.Default()
	RegisterUserRoutes(router, conn, &ctx, queries)

	// Replace the following with your actual test scenarios
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"CreateUser", "POST", "/users/create", http.StatusNoContent},
		{"SignIn", "POST", "/users/sign-in", http.StatusNoContent},
		{"GetUser", "GET", "/users/get", http.StatusBadRequest},
		{"GetAllUser", "GET", "/users/get-all-users", 401},
		{"UploadProfile", "POST", "/users/upload-profile", 401},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, recorder.Code)
			}
		})
	}
}
