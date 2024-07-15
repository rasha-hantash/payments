CREATE EXTENSION IF NOT EXISTS citext;

-- Create users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email citext UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL DEFAULT 'system',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL DEFAULT 'system',
    history JSONB[] DEFAULT '{}'
);

-- todo remove account_type
-- Create accounts table
CREATE TABLE accounts (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id),
    account_type text NOT NULL,
    account_state text NOT NULL, -- e.g., 'open', 'closed', 'frozen'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL DEFAULT 'system',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL DEFAULT 'system',
    history JSONB[] DEFAULT '{}'
);

CREATE TABLE payment_methods (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    method_type VARCHAR(50) NOT NULL, -- e.g., 'ACH', 'paypal', 'venmo'
    account_number VARCHAR(255),
    routing_number VARCHAR(255),
    card_number VARCHAR(255),
    expiration_date DATE,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL DEFAULT 'system',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL DEFAULT 'system',
    history JSONB[] DEFAULT '{}'
);

-- Create transactions table
CREATE TABLE transactions (
    id TEXT PRIMARY KEY,
    external_payment_method_id TEXT REFERENCES payment_methods(id),
    is_internal BOOLEAN,
    amount BIGINT NOT NULL,
    status text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL DEFAULT 'system',
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL DEFAULT 'system',
    history JSONB[] DEFAULT '{}'
);

CREATE TABLE ledger_entries (
    id TEXT PRIMARY KEY ,
    transaction_id TEXT NOT NULL REFERENCES transactions(id),
    account_id TEXT NOT NULL REFERENCES accounts(id),
    direction text NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL DEFAULT 'system'
);



-- -- Create index on transactions for faster querying
-- CREATE INDEX idx_transactions_account_id ON transactions(account_id);

-- Create function to update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE OR REPLACE FUNCTION history_trigger_function()
    RETURNS trigger
    LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
    NEW.history := array_append(OLD.history, (to_jsonb(OLD) - 'history'));
    RETURN NEW;
END;
$BODY$;

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER payment_methods_transactions_updated_at BEFORE UPDATE ON payment_methods FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- DROP TRIGGER IF EXISTS userprofile_history_trigger_fn on userprofile;
CREATE TRIGGER users_history_trigger_fn BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE history_trigger_function ();
CREATE TRIGGER accounts_history_trigger_fn BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE history_trigger_function ();
CREATE TRIGGER transactions_history_trigger_fn BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE history_trigger_function ();
CREATE TRIGGER payment_methods_history_trigger_fn BEFORE UPDATE ON payment_methods FOR EACH ROW EXECUTE FUNCTION history_trigger_function();


-- TODO add COMMENT ON statements