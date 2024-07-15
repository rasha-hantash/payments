-- Drop tables
DROP TABLE IF EXISTS ledger_entries;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;


-- Drop history triggers
DROP TRIGGER IF EXISTS users_history_trigger_fn ON users;
DROP TRIGGER IF EXISTS accounts_history_trigger_fn ON accounts;
DROP TRIGGER IF EXISTS transactions_history_trigger_fn ON transactions;
DROP TRIGGER IF EXISTS payment_methods_history_trigger_fn ON payment_methods;

-- Drop updated_at triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS payment_methods_transactions_updated_at ON payment_methods;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS history_trigger_function();

-- Drop extensions
DROP EXTENSION IF EXISTS citext;
