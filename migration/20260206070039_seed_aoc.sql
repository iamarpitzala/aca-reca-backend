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
    -- Revenue
    SELECT 'Revenue' AS account_type, 'GST on Income' AS account_tax, '200' AS code, 'Sales' AS name, 'Sales Revenue' AS description UNION ALL
    SELECT 'Revenue', 'GST Free Income', '222', 'Demo Sales', 'GST Free Income' UNION ALL
    SELECT 'Revenue', 'GST on Income', '260', 'Other Revenue', 'Other Revenue' UNION ALL
    SELECT 'Revenue', 'GST Free Income', '270', 'Interest Income', 'Interest Income' UNION ALL

    -- Cost of Goods Sold (Direct Costs → Expense)
    SELECT 'Expense', 'GST on Expenses', '310', 'Cost of Goods Sold', 'Direct Costs' UNION ALL

    -- Expenses
    SELECT 'Expense', 'GST on Expenses', '400', 'Advertising', 'Advertising Expenses' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '404', 'Bank Fees', 'Bank Fees' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '408', 'Cleaning', 'Cleaning Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '412', 'Consulting & Accounting', 'Consulting & Accounting' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '416', 'Depreciation', 'Depreciation' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '420', 'Entertainment', 'Entertainment' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '425', 'Freight & Courier', 'Freight & Courier' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '429', 'General Expenses', 'General Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '433', 'Insurance', 'Insurance' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '437', 'Interest Expense', 'Interest Expense' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '441', 'Legal Expenses', 'Legal Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '445', 'Light, Power & Heating', 'Utilities' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '449', 'Motor Vehicle Expenses', 'Motor Vehicle' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '453', 'Office Expenses', 'Office Expenses' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '461', 'Printing & Stationery', 'Printing & Stationery' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '469', 'Rent', 'Rent' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '473', 'Repairs & Maintenance', 'Repairs & Maintenance' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '477', 'Wages & Salaries', 'Payroll' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '478', 'Superannuation', 'Superannuation' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '485', 'Subscriptions', 'Subscriptions' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '489', 'Telephone & Internet', 'Telephone & Internet' UNION ALL
    SELECT 'Expense', 'GST on Expenses', '493', 'Travel - National', 'Travel National' UNION ALL
    SELECT 'Expense', 'GST Free Expenses', '494', 'Travel - International', 'Travel International' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '497', 'Bank Revaluations', 'Bank Revaluations' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '498', 'Unrealised Currency Gains', 'Unrealised Currency Gains' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '499', 'Realised Currency Gains', 'Realised Currency Gains' UNION ALL
    SELECT 'Expense', 'BAS Excluded', '505', 'Income Tax Expense', 'Income Tax' UNION ALL

    -- Assets
    SELECT 'Asset', 'BAS Excluded', '610', 'Accounts Receivable', 'Current Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '620', 'Prepayments', 'Current Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '630', 'Inventory', 'Current Asset' UNION ALL
    SELECT 'Asset', 'GST on Expenses', '710', 'Office Equipment', 'Fixed Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '711', 'Accumulated Depreciation – Office Equipment', 'Contra Asset' UNION ALL
    SELECT 'Asset', 'GST on Expenses', '720', 'Computer Equipment', 'Fixed Asset' UNION ALL
    SELECT 'Asset', 'BAS Excluded', '721', 'Accumulated Depreciation – Computer Equipment', 'Contra Asset' UNION ALL

    -- Liabilities
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

    -- Equity
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
