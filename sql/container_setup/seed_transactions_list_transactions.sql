-- Seed script for 50 transactions (100 ledger entries) with deterministic values

-- Ensure the users and accounts exist
INSERT INTO users (id, email, name) VALUES
('usr_1', 'hello+1@gmail.com', 'User 1')
ON CONFLICT (id) DO NOTHING;

INSERT INTO accounts (id, user_id, account_state, account_type) VALUES
('acct_1', 'usr_1', 'open', 'debit'),
('acct_2', 'usr_1', 'open', 'credit'),
('acct_3', 'usr_1', 'open', 'debit')
ON CONFLICT (id) DO NOTHING;

-- Insert 50 transactions and 100 ledger entries
DO $$
DECLARE
    i INT;
    txn_id TEXT;
    amount DECIMAL(10, 2);
    from_account TEXT;
    to_account TEXT;
BEGIN
    FOR i IN 1..50 LOOP
        txn_id := 'txn_' || i;
        amount := 100.00 + (i * 10.00);  -- Each transaction increases by $10
        
        -- Deterministic account selection
        CASE i % 3
            WHEN 0 THEN
                from_account := 'acct_1';
                to_account := 'acct_2';
            WHEN 1 THEN
                from_account := 'acct_2';
                to_account := 'acct_3';
            ELSE
                from_account := 'acct_3';
                to_account := 'acct_1';
        END CASE;

        -- Insert transaction
        INSERT INTO transactions (id, amount, status) VALUES
        (txn_id, amount, 'success');

        -- Insert ledger entries
        INSERT INTO ledger_entries (id, transaction_id, account_id, amount, direction) VALUES
        ('le_' || (i*2-1), txn_id, from_account, amount, 'debit'),
        ('le_' || (i*2), txn_id, to_account, amount, 'credit');
    END LOOP;
END $$;

-- final balances
-- acct_1: -350
-- acct_2: -170
-- acct_3: 520