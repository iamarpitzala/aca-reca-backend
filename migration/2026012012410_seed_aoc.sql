-- +goose Up
-- +goose StatementBegin

INSERT INTO tbl_account (
    account_type_id,
    account_tax_id,
    code,
    name,
    description
)
SELECT
    account_type.id,
    account_tax.id,
    accounts.code,
    accounts.name,
    accounts.description
FROM (
    -- Revenue Accounts (Updated list - only Patient Fee Account, Commission Received, Other Income)
    SELECT 'Revenue' AS account_type, 'GST Free Income' AS account_tax, '200' AS code, 'Patient Fee Account' AS name, 'Patient Fee Account' AS description UNION ALL
    SELECT 'Revenue', 'GST on Income' AS account_tax, '201' AS code, 'Commission Received' AS name, 'Commission Received' AS description UNION ALL
    SELECT 'Revenue', 'GST on Income' AS account_tax, '202' AS code, 'Other Income' AS name, 'Other Income' AS description UNION ALL
    -- Expense Accounts (Updated list - Computer expense removed)
    SELECT 'Expense', 'GST on Expenses', '400', 'Home Office (GST)', 'Home Office Expenses (GST)' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '401', 'Home Office (GST Free)', 'Home Office Expenses (GST Free)' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '402', 'Laboratory Work (GST Free)', 'Laboratory Work Expenses (GST Free)' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '403', 'Laboratory Work (GST)', 'Laboratory Work Expenses (GST)' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '404', 'Subscription/Membership (GST Free)', 'Subscription/Membership Expenses (GST Free)' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '405', 'Subscription/Membership (GST)', 'Subscription/Membership Expenses (GST)' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '406', 'Bank Fees', 'Bank Fees' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '407', 'Merchant Fees', 'Merchant Fees' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '408', 'Motor Vehicle - Set Rate', 'Motor Vehicle Expenses - Set Rate' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '409', 'M/V Insurance', 'Motor Vehicle Insurance' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '410', 'M/V Registration', 'Motor Vehicle Registration' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '411', 'M/V Fuel', 'Motor Vehicle Fuel' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '412', 'M/V Repairs/Maintenance', 'Motor Vehicle Repairs and Maintenance' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '413', 'Management Fee (Gross Up)', 'Management Fee (Gross Up)' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '414', 'Materials/Dental Supplies', 'Materials and Dental Supplies' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '415', 'Office Supplies', 'Office Supplies' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '416', 'Postage', 'Postage Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '417', 'Protective Clothing', 'Protective Clothing Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '418', 'Internet', 'Internet Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '419', 'Telephone', 'Telephone Expenses' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '420', 'Telephone and Internet (GST Free)', 'Telephone and Internet Expenses (GST Free)' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '421', 'Travel/Accommodation (GST Free)', 'Travel and Accommodation Expenses (GST Free)' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '422', 'Travel/Accommodation (GST)', 'Travel and Accommodation Expenses (GST)' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '423', 'Tolls / Parking', 'Tolls and Parking Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '424', 'Waste Disposal', 'Waste Disposal Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '425', 'Repairs and Maintenance', 'Repairs and Maintenance Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '426', 'Sundries', 'Sundry Expenses' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '610', 'Accounts Receivable', 'Current Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '620', 'Prepayments', 'Current Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '630', 'Inventory', 'Current Asset' UNION ALL
    SELECT 'Asset', 'GST on Expenses', '710', 'Office Equipment', 'Fixed Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '711', 'Accumulated Depreciation – Office Equipment', 'Contra Asset' UNION ALL
    SELECT 'Asset', 'GST on Expenses', '720', 'Computer Equipment', 'Fixed Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '721', 'Accumulated Depreciation – Computer Equipment', 'Contra Asset' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '800', 'Accounts Payable', 'Current Liability' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '801', 'Unpaid Expense Claims', 'Current Liability' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '804', 'Wages Payable – Payroll', 'Payroll Liability' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '820', 'GST', 'GST Payable' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '825', 'PAYG Withholdings Payable', 'PAYG' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '826', 'Superannuation Payable', 'Superannuation' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '830', 'Income Tax Payable', 'Income Tax' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '840', 'Historical Adjustment', 'Historical Adjustment' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '850', 'Suspense', 'Suspense' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '860', 'Rounding', 'Rounding' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '877', 'Tracking Transfers', 'Tracking Transfers' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '880', 'Owner A Drawings', 'Drawings' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '881', 'Owner A Funds Introduced', 'Capital Introduced' UNION ALL
    SELECT 'Liability', 'BAS Excluded', '900', 'Loan', 'Non-Current Liability' UNION ALL
    SELECT 'Equity', 'BAS Excluded', '960', 'Retained Earnings', 'Retained Earnings' UNION ALL
    SELECT 'Equity', 'BAS Excluded', '970', 'Owner A Share Capital', 'Share Capital'
) AS accounts
JOIN tbl_account_type AS account_type ON accounts.account_type = account_type.name
JOIN tbl_account_tax AS account_tax ON accounts.account_tax = account_tax.name
ON CONFLICT (code) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
