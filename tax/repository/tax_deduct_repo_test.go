package repository

import (
	"fmt"
	"testing"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTaxDeductConfigRepo_UpdateById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTaxDeductConfigRepo(db)

	expectedID := "1"
	expectedAmount := 100.0

	mock.ExpectPrepare(`UPDATE tax_deduct_config SET amount = \$1 WHERE deduct_id = \$2`).
		ExpectExec().
		WithArgs(expectedAmount, expectedID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	numRows, err := repo.UpdateById(expectedID, expectedAmount)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), numRows)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaxDeductConfigRepo_FindById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewTaxDeductConfigRepo(db)

	expectedID := "1"

	rows := sqlmock.NewRows([]string{"deduct_id", "amount", "description"}).
		AddRow(expectedID, 100.0, "Description")

	mock.ExpectPrepare(`SELECT deduct_id\s*,\s*amount\s*,\s*description\s*FROM tax_deduct_config\s*WHERE deduct_id = \$1`).
		ExpectQuery().
		WithArgs(expectedID).
		WillReturnRows(rows)

	_, err = repo.FindById(expectedID)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaxDeductConfigRepo_FindById_NotFound(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := NewTaxDeductConfigRepo(db)

    expectedID := "1"

    // Expecting the prepare query
    rows := sqlmock.NewRows([]string{"deduct_id", "amount", "description"})
    mock.ExpectPrepare(`SELECT deduct_id\s*,\s*amount\s*,\s*description\s*FROM tax_deduct_config\s*WHERE deduct_id = \$1`).
        ExpectQuery().
        WithArgs(expectedID).
        WillReturnRows(rows)

    // Call the method under test
    _, err = repo.FindById(expectedID)

    // Assert the error message
    expectedErrorMsg := fmt.Sprintf("deduct config not found for ID: %s", expectedID)
    if err.Error() != expectedErrorMsg {
        t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
    }

    assert.NoError(t, mock.ExpectationsWereMet())
}