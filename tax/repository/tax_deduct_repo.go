package repository

import (
	"database/sql"
	"fmt"
)

type TaxDeductConfigRepo struct {
	Db *sql.DB
}


func NewTaxDeductConfigRepo(db *sql.DB) TaxDeductConfigPort {
	return &TaxDeductConfigRepo{Db: db}
}



func (t *TaxDeductConfigRepo) UpdateById(id string , amount float64) (int64,error){
	

	query := ` UPDATE  
					tax_deduct_config
				
				SET
					amount = $1
					
				WHERE 
					deduct_id = $2 `

	
	stmt , err :=  t.Db.Prepare(query)

	if err !=nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(amount,id)
	
	if err != nil {
        return 0, err
    }

    // Get the number of rows affected
    numRows, err := res.RowsAffected()
    if err != nil {
        return 0, err
    }

    return numRows, nil
}

func (t *TaxDeductConfigRepo) FindById(id string) (*TaxDeductConfig,error){

	query := `
				SELECT 
					deduct_id , amount , description 
				FROM 
					tax_deduct_config 
				WHERE 
					deduct_id = $1 `

	stmt , err :=  t.Db.Prepare(query)

	if err !=nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query using the QueryRow method of the DB object
	row := stmt.QueryRow(id)

	var  tdc TaxDeductConfig

	// Scan the values returned by the query into the fields of the wallet struct
	err = row.Scan(&tdc.DeductId, &tdc.Amount, &tdc.Description)
	if err != nil {
		// If no rows are returned, check for sql.ErrNoRows error
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deduct config not found for ID: %s", id)
		}
		// Otherwise, return any other error
		return nil, err
	}

	// Return the populated wallet pointer
	return &tdc, nil

}



