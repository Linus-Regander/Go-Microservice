package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	userModel "github.com/Linus-Regander/Go-Microservice/internal/api/model/user"
	"github.com/Linus-Regander/Go-Microservice/pkg/database/query"
	"github.com/blockloop/scan/v2"
	"github.com/huandu/go-sqlbuilder"
)

// UserRepository holds information about a UserRepository
type UserRepository struct {
	db  *sql.DB
	log *log.Logger
}

// New returns a new instance of user repository,
func New(log *log.Logger, db *sql.DB) *UserRepository {
	//
	// TODO: Setup db with migrations.
	//

	return &UserRepository{
		db:  db,
		log: log,
	}
}

// DeleteUser deletes a user from store.
func (ur *UserRepository) DeleteUser(ctx context.Context, userId string) error {
	conn, err := ur.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func(conn *sql.Conn) {
		if err = conn.Close(); err != nil {
			ur.log.Print(fmt.Errorf("error closing connection: %w", err))
		}
	}(conn)

	query, args := ur.deleteQuery(userId)

	if _, err = conn.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

// InsertUser inserts a new user to store.
func (ur *UserRepository) InsertUser(ctx context.Context, user userModel.User) error {
	conn, err := ur.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func(conn *sql.Conn) {
		if err = conn.Close(); err != nil {
			ur.log.Print(fmt.Errorf("error closing connection: %w", err))
		}
	}(conn)

	query, args := ur.insertQuery(user)

	if _, err = conn.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user in store.
func (ur *UserRepository) UpdateUser(ctx context.Context, user userModel.User) error {
	conn, err := ur.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func(conn *sql.Conn) {
		if err = conn.Close(); err != nil {
			ur.log.Print(fmt.Errorf("error closing connection: %w", err))
		}
	}(conn)

	query, args := ur.updateQuery(user)

	if _, err = conn.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

// SelectAllUsers returns users, with filtering and pagination data from store.
func (ur *UserRepository) SelectAllUsers(ctx context.Context, userParams userModel.Params) (userModel.Users, userModel.Pagination, error) {
	var (
		users      userModel.Users
		pagination userModel.Pagination
	)

	conn, err := ur.db.Conn(ctx)
	if err != nil {
		return users, pagination, err
	}
	defer func(conn *sql.Conn) {
		if err = conn.Close(); err != nil {
			ur.log.Print(fmt.Errorf("error closing connection: %w", err))
		}
	}(conn)

	query, args := ur.selectQuery(userParams)

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return users, pagination, err
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			ur.log.Print(fmt.Errorf("error closing rows: %w", err))
		}
	}(rows)

	if err = scan.RowsStrict(&users, rows); err != nil {
		return users, pagination, err
	}

	//
	// TODO: Get pagination data from select.
	//

	return users, pagination, nil
}

//
// Unexported query functions.
// For testability.
//

func (ur *UserRepository) deleteQuery(userId string) (string, []interface{}) {
	deleteBuilder := sqlbuilder.NewDeleteBuilder()

	return deleteBuilder.
		DeleteFrom(userModel.UsersTableName).
		Where(deleteBuilder.Equal(userModel.UserIdFieldName, userId)).
		Build()
}

func (ur *UserRepository) insertQuery(user userModel.User) (string, []interface{}) {
	insertBuilder := sqlbuilder.NewInsertBuilder()

	return insertBuilder.
		InsertInto(userModel.UsersTableName).
		Cols(query.ColumnNames(user)...).
		Values(query.Values(user)...).
		Build()
}

func (ur *UserRepository) updateQuery(user userModel.User) (string, []interface{}) {
	updateBuilder := sqlbuilder.
		NewUpdateBuilder().
		Update(userModel.UsersTableName)

	assignments := make([]string, 0)

	for column, value := range query.ColumnValues(user) {
		assignments = append(
			assignments,
			updateBuilder.Assign(column, value),
		)
	}

	updateBuilder.Set(assignments...)

	return updateBuilder.Where(
		updateBuilder.Equal(userModel.UserIdFieldName, user.ID),
	).Build()
}

func (ur *UserRepository) selectQuery(userParams userModel.Params) (string, []interface{}) {
	selectBuilder := sqlbuilder.
		NewSelectBuilder().
		Select(query.ColumnNames(userModel.User{})...).
		From(userModel.UsersTableName)

	conditions := query.ColumnValues(userParams)

	if len(conditions) > 0 {
		for condition, arg := range conditions {
			if arg == nil || arg == "" {
				continue
			}

			selectBuilder.Where(selectBuilder.Equal(condition, arg))
		}
	}

	return selectBuilder.Build()
}
