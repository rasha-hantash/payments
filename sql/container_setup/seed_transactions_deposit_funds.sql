-- Seed users (assuming a users table exists)
INSERT INTO users (id, email, name) VALUES
('usr_1', 'hello@gmail.com', 'User 1'),
('usr_2', 'hello+2@gmail.com','User 2'),
('usr_3', 'hello+3@gmail.com', 'User 3');

INSERT INTO accounts (id, account_state, account_type) VALUES 
    ('acct_1', 'open', 'internal'), -- user-1 internal account 
    ('acct_2', 'open', 'external'),
    ('acct_4', 'open', 'external'),-- user-1 external account (ex: bank account)
    ('acct_5', 'open', 'internal'); -- user-1 external account (ex: bank account)


-- Seed payment methods for successful deposit
-- INSERT INTO payment_methods (id, user_id, method_type, is_verified) VALUES
-- ('acct_1', 1, 'ACH', TRUE),
-- ('acct_2', 1, 'ACH', TRUE);

-- -- Seed payment method for insert error scenarios
-- -- Credit account for debit account does not exist scenario
-- INSERT INTO payment_methods (id, user_id, method_type, is_verified) VALUES
-- ('acct_4', 2, 'ACH', TRUE);
-- -- Debit account for credit account does not exist scenario
-- INSERT INTO payment_methods (id, user_id, method_type, is_verified) VALUES
-- ('acct_5', 3, 'ACH', TRUE);

-- Seed a successful transaction (assuming linkage to payment_methods or accounts is handled elsewhere)
-- INSERT INTO transactions (id, amount, status) VALUES
-- ('txn_1', 100, 'success');

-- Note: For the error scenarios, no transactions are seeded as they represent failed attempts to insert due to missing accounts.