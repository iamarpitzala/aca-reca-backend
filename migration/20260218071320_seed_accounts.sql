-- +goose Up
-- +goose StatementBegin

INSERT INTO tbl_account (
    id,
    owner_user_id,
    account_type_id,
    account_tax_id,
    code,
    name,
    description
)
SELECT
    gen_random_uuid()::varchar(40),
    u.id AS owner_user_id,
    at.id,
    tax.id,
    a.code,
    a.name,
    a.description
FROM (
    -- Revenue
    SELECT 'Revenue' AS account_type, 'GST Free Income' AS account_tax, '200' AS code, 'Patient Fee Account' AS name, 'Patient Fee Account' AS description UNION ALL
    SELECT 'Revenue','GST on Income','201','Commission Received','Commission Received' UNION ALL
    SELECT 'Revenue','GST on Income','202','Other Income','Other Income' UNION ALL

    -- Expenses
    SELECT 'Expense','GST on Expenses','400','Home Office (GST)','Home Office Expenses (GST)' UNION ALL
    SELECT 'Expense','GST Free Expenses','401','Home Office (GST Free)','Home Office Expenses (GST Free)' UNION ALL
    SELECT 'Expense','GST Free Expenses','402','Laboratory Work (GST Free)','Laboratory Work Expenses (GST Free)' UNION ALL
    SELECT 'Expense','GST on Expenses','403','Laboratory Work (GST)','Laboratory Work Expenses (GST)' UNION ALL
    SELECT 'Expense','GST Free Expenses','404','Subscription/Membership (GST Free)','Subscription/Membership Expenses (GST Free)' UNION ALL
    SELECT 'Expense','GST on Expenses','405','Subscription/Membership (GST)','Subscription/Membership Expenses (GST)' UNION ALL
    SELECT 'Expense','GST Free Expenses','406','Bank Fees','Bank Fees' UNION ALL
    SELECT 'Expense','GST Free Expenses','407','Merchant Fees','Merchant Fees' UNION ALL

    -- Assets
    SELECT 'Asset','BAS Excluded','610','Accounts Receivable','Current Asset' UNION ALL
    SELECT 'Asset','BAS Excluded','620','Prepayments','Current Asset' UNION ALL
    SELECT 'Asset','BAS Excluded','630','Inventory','Current Asset' UNION ALL

    -- Liabilities
    SELECT 'Liability','BAS Excluded','800','Accounts Payable','Current Liability' UNION ALL
    SELECT 'Liability','BAS Excluded','820','GST','GST Payable' UNION ALL
    SELECT 'Liability','BAS Excluded','900','Loan','Non-Current Liability' UNION ALL

    -- Equity
    SELECT 'Equity','BAS Excluded','960','Retained Earnings','Retained Earnings' UNION ALL
    SELECT 'Equity','BAS Excluded','970','Owner A Share Capital','Share Capital'
) a
JOIN tbl_account_type at ON at.name = a.account_type
JOIN tbl_account_tax tax ON tax.name = a.account_tax
CROSS JOIN tbl_user u
WHERE NOT EXISTS (
    SELECT 1 FROM tbl_account ac
    WHERE ac.owner_user_id = u.id AND ac.code = a.code AND (ac.deleted_at IS NULL)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd