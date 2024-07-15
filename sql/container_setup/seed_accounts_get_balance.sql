INSERT INTO accounts (id, account_state, account_type) VALUES 
    ('acct_1', 'open', 'internal'), -- user-1 internal account 
    ('acct_2', 'open', 'external'),
    ('acct_3', 'open', 'external'); -- user-1 external account (ex: bank account)

-- Insert transactions for account-1
INSERT INTO transactions (id, amount, status) VALUES 
    ('txn_1', 500, 'in_process'),
    ('txn_2', 50, 'paid');

INSERT INTO ledger_entries (id, transaction_id, custom_id, account_id, direction, amount) VALUES 
    ('le_2', 'txn_1', 'txn_1', 'acct_2', 'debit', 500),
    ('le_1', 'txn_1', 'txn_1', 'acct_1', 'credit', 500),
    ('le_3', 'txn_2', 'txn_2', 'acct_1', 'debit', 50),
    ('le_4', 'txn_2', 'txn_2', 'acct_2', 'credit', 50);
