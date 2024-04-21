CREATE TABLE tax_rates (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    IncomeRangeFrom DECIMAL(15, 2),
    IncomeRangeTo DECIMAL(15, 2),
    TaxRate DECIMAL(5, 2)
);

-- data init
INSERT INTO TaxRates (IncomeRangeFrom, IncomeRangeTo, TaxRate)
VALUES
    (0, 150000, 0),
    (150001, 500000, 0.10),
    (500001, 1000000, 0.15),
    (1000001, 2000000, 0.20),
    (2000001, NULL, 0.35);