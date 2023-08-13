package gormdb_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anang5u/simple-gorm-unit-test/gormdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCreateUser(t *testing.T) {
	// Membuat mock database connection
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	// Inisialisasi GORM dengan mock database connection
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Tidak menampilkan log
	})
	assert.NoError(t, err)

	fullName := "John Doe"
	email := "johndoe@example.com"

	// Menyiapkan ekspetasi pada mock
	mock.ExpectBegin() // Mengharapkan transaksi dimulai
	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(sqlmock.AnyArg(), fullName, email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mengharapkan hasil query berhasil
	mock.ExpectCommit() // Mengharapkan transaksi di-commit

	// Memanggil fungsi CreateUser
	user, err := gormdb.CreateUser(db, fullName, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, fullName, user.FullName)
	assert.Equal(t, email, user.Email)

	bb, _ := json.MarshalIndent(user, " ", " ")
	log.Println(string(bb))

	// Memastikan bahwa semua ekspetasi pada mock telah terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	// Membuat mock database connection
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	// Inisialisasi GORM dengan mock database connection
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Tidak menampilkan log
	})
	assert.NoError(t, err)

	userID := uuid.New()

	// Menyiapkan ekspetasi pada mock
	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "full_name", "email", "created"}).
				AddRow(userID, "John Doe", "johndoe@example.com", time.Now()))

	// Memanggil fungsi GetUserByID
	user, err := gormdb.GetUserByID(db, userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.FullName)
	assert.Equal(t, "johndoe@example.com", user.Email)

	// Memastikan bahwa semua ekspetasi pada mock telah terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())

	bb, _ := json.MarshalIndent(user, " ", " ")
	log.Println(string(bb))
}

func TestUpdateUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})

	// Data testing
	userID := uuid.New()
	newEmail := "newemail@example.com"

	// Menyiapkan ekspetasi pada mock
	mock.ExpectBegin() // Mengharapkan transaksi dimulai
	mock.ExpectExec(`UPDATE "users"`).
		WithArgs(newEmail, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mengharapkan hasil query berhasil
	mock.ExpectCommit() // Mengharapkan transaksi di-commit

	// Melakukan pengujian
	rowsAffected, err := gormdb.UpdateUser(db, userID, newEmail)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected) // Memastikan 1 baris diubah

	// Memastikan semua ekspetasi terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())

	log.Println(err)
	log.Println(rowsAffected)
}

func TestDeleteUser(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{})

	// Data testing
	userID := uuid.New()

	// Menyiapkan ekspetasi pada mock
	mock.ExpectBegin() // Mengharapkan transaksi dimulai
	mock.ExpectExec(`DELETE FROM "users"`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Mengharapkan hasil query berhasil (1 row affected)
	mock.ExpectCommit() // Mengharapkan transaksi di-commit

	// Melakukan pengujian
	err := gormdb.DeleteUser(db, userID)

	assert.NoError(t, err)

	// Memastikan semua ekspetasi terpenuhi
	assert.NoError(t, mock.ExpectationsWereMet())
}
