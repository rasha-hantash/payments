-- Insert users
INSERT INTO users (id, email, name) VALUES
('usr_1', 'hello+1@gmail.com', 'User 1'),
('usr_2','hello+2@gmail.com', 'User 2'),
('usr_3', 'hello+3@gmail.com', 'User 3');

-- Insert accounts
-- Assuming accounts table has columns for account ID, user ID, and balance
INSERT INTO accounts (id, user_id, account_state, account_type) VALUES
('acct_1', 'usr_1', 'open', 'debit'), -- Sufficient funds for withdrawal
('acct_2', 'usr_1', 'open', 'credit'),   -- Target account for successful withdrawal
('acct_3', 'usr_2', 'open','debit'),  -- Insufficient funds
('acct_5', 'usr_3', 'open', 'credit'); -- Exists but credit account does not

-- Note: acct-4 and acct_6 are intentionally omitted to simulate "account does not exist" scenarios

-- Insert transactions
-- Assuming transactions table has columns for transaction ID, debit account ID, credit account ID, amount, and status
-- Successful withdraw
INSERT INTO transactions (id, amount, status) VALUES
('txn_1', 100, 'success'),
('txn_2', 50, 'success');

-- Insufficient funds (attempted transaction, might not be recorded if validation happens before insertion)
-- Credit account does not exist (attempted transaction, might not be recorded for the same reason)

-- Insert ledger entries for the successful withdrawal
-- Assuming ledger_entries table records individual account movements per transaction
INSERT INTO ledger_entries (id, transaction_id, account_id, amount, direction) VALUES
('le_1', 'txn_1', 'acct_1', 100, 'credit'), -- Debit from acct_1
('le_2', 'txn_1', 'acct_2', 100, 'debit'), 
('le_3', 'txn_2', 'acct_1', 50, 'credit'), -- Debit from acct_1
('le_4', 'txn_2', 'acct_2', 50, 'debit'); -- Credit to acct_2

-- Note: For "insufficient funds" and "account does not exist" scenarios, ledger entries are not created as these transactions would fail validation checks before affecting account balances.