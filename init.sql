CREATE TABLE tax_deduct_config (
    deduct_id CHARACTER(10) PRIMARY KEY,
    amount DECIMAL(15, 2) -- Assuming maximum precision of 15 digits with 2 decimal places
    description CHARACTER(100)
); 

-- Inserting sample data into tax_deduct_config table
INSERT INTO tax_deduct_config (deduct_type, amount, description) VALUES
    ('personal', 60000.00, 'Personal allowance'),
    ('k-receipt', 50000.00, 'k-receipt allowance');
