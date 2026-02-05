package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDeleteProduct(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()

	dialector := postgres.New(postgres.Config{Conn: dbRaw})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	repo := NewSearchRepository(db)
	productID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM \"product_search_models\"").
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(productID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdvancedSearch(t *testing.T) {
	dbRaw, mock, _ := sqlmock.New()
	defer dbRaw.Close()
	dialector := postgres.New(postgres.Config{Conn: dbRaw})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	repo := NewSearchRepository(db)

	rowsCount := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM \"product_search_models\"").WillReturnRows(rowsCount)

	rowsItems := sqlmock.NewRows([]string{"id", "name", "price", "category"}).
		AddRow(uuid.New(), "Pasta", 15.0, "Italian")

	mock.ExpectQuery("SELECT \\* FROM \"product_search_models\" WHERE name ILIKE .*").
		WillReturnRows(rowsItems)

	result, err := repo.AdvancedSearch("Pasta", "", 0, 0, 1, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.TotalItems)
	assert.NoError(t, mock.ExpectationsWereMet())
}
