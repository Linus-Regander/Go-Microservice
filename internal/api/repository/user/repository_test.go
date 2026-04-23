package user

import (
	"testing"
	"time"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	"github.com/stretchr/testify/require"
)

func Test_DeleteQuery(t *testing.T) {
	testCases := map[string]struct {
		expectedArgs  []interface{}
		expectedQuery string
		userId        string
	}{
		"DeleteQuery - correct query and args - with user id": {
			userId: "7c2f1c3e-6c2b-4f3a-9d8d-1e9f2b8c4a61",
			expectedArgs: []interface{}{
				"7c2f1c3e-6c2b-4f3a-9d8d-1e9f2b8c4a61",
			},
			expectedQuery: "DELETE FROM users WHERE id = ?",
		},
		"DeleteQuery - correct query and args - without user id": { // should theoretically not happen since service handles empty, invalid ids.
			userId: "",
			expectedArgs: []interface{}{
				"",
			},
			expectedQuery: "DELETE FROM users WHERE id = ?",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ur := &UserRepository{}

			query, args := ur.deleteQuery(tc.userId)
			require.Equal(t, tc.expectedArgs, args)
			require.Equal(t, tc.expectedQuery, query)
		})
	}
}

func Test_InsertQuery(t *testing.T) {
	testCases := map[string]struct {
		expectedArgs  []interface{}
		expectedQuery string
		user          userModel.User
	}{
		"InsertQuery - correct query and args - with user body": {
			user: userModel.User{
				ID:        "7c2f1c3e-6c2b-4f3a-9d8d-1e9f2b8c4a61",
				Username:  "Test",
				Name:      "Test Testsson",
				Role:      userModel.Admin,
				CreatedAt: time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
				UpdatedAt: time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
			},
			expectedArgs: []interface{}{
				"7c2f1c3e-6c2b-4f3a-9d8d-1e9f2b8c4a61",
				"Test",
				"Test Testsson",
				userModel.Admin,
				time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
				time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
			},
			expectedQuery: "INSERT INTO users (id, username, name, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		},
		"InsertQuery - correct query and args - with empty user body": {
			user: userModel.User{},
			expectedArgs: []interface{}{
				"",
				"",
				"",
				userModel.Role(""),
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedQuery: "INSERT INTO users (id, username, name, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ur := &UserRepository{}

			query, args := ur.insertQuery(tc.user)
			require.Equal(t, tc.expectedArgs, args)
			require.Equal(t, tc.expectedQuery, query)
		})
	}
}

func Test_UpdateQuery(t *testing.T) {
	testCases := map[string]struct {
		expectedArgs []interface{}
		user         userModel.User
	}{
		"UpdateQuery - correct query and args - with user body": {
			user: userModel.User{
				ID:        "7c2f1c3e-6c2b-4f3a-9d8d-1e9f2b8c4a61",
				Username:  "Test",
				Name:      "Test Testsson",
				Role:      userModel.Admin,
				CreatedAt: time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
				UpdatedAt: time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
			},
			expectedArgs: []interface{}{
				"7c2f1c3e-6c2b-4f3a-9d8d-1e9f2b8c4a61",
				"Test",
				"Test Testsson",
				userModel.Admin,
				time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
				time.Date(2020, time.February, 3, 4, 5, 6, 0, time.UTC),
			},
		},
		"UpdateQuery - correct query and args - with empty user body": {
			user: userModel.User{},
			expectedArgs: []interface{}{
				"",
				"",
				"",
				userModel.Role(""),
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ur := &UserRepository{}

			query, args := ur.updateQuery(tc.user)
			require.Subset(t, args, tc.expectedArgs)
			require.Contains(t, query, "UPDATE users")
			require.Contains(t, query, "WHERE id = ?")

			//
			// TODO: Add ability to test args in query, since they are none-deterministic right now.
			//
		})
	}
}

func Test_SelectAllQuery(t *testing.T) {
	testCases := map[string]struct {
		expectedArgs  []interface{}
		expectedQuery string
		userParams    userModel.Params
	}{
		"SelectQuery - correct query and args - with params": {
			userParams: userModel.Params{
				Name: "Test Testsson",
			},
			expectedArgs: []interface{}{
				"Test Testsson",
			},
			expectedQuery: "SELECT id, username, name, role, created_at, updated_at FROM users WHERE name = ?",
		},
		"SelectQuery - correct query and args - without params": {
			userParams:    userModel.Params{},
			expectedArgs:  []interface{}(nil),
			expectedQuery: "SELECT id, username, name, role, created_at, updated_at FROM users",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ur := &UserRepository{}

			query, args := ur.selectQuery(tc.userParams)
			require.Equal(t, tc.expectedArgs, args)
			require.Equal(t, tc.expectedQuery, query)

			//
			// TODO: Add ability to test args in query, since they are none-deterministic right now.
			//
		})
	}
}
