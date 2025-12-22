-- Account Service Schema Design
-- Professional Banking Backend for Educational Purposes

-- Customers table to store banking-specific info
-- Links to external Auth system via external_id
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(100) UNIQUE NOT NULL, -- Reference to the User ID from the Auth service
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Accounts table supporting multiple currencies per customer
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id),
    account_number VARCHAR(34) UNIQUE NOT NULL, -- IBAN or internal unique format
    currency CHAR(3) NOT NULL, -- ISO 4217 Currency Code (e.g., PLN, USD, EUR)
    balance NUMERIC(20, 4) NOT NULL DEFAULT 0.0000,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, frozen, closed
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Ledger for all balance-changing operations (Audit Trail)
-- Each entry records the balance before and after for strict auditability
CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id),
    type VARCHAR(50) NOT NULL, -- deposit, withdrawal, transfer_in, transfer_out, fx_exchange
    amount NUMERIC(20, 4) NOT NULL,
    balance_after NUMERIC(20, 4) NOT NULL,
    reference_id UUID, -- To link related entries (e.g., in transfers)
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_accounts_customer_id ON accounts(customer_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_account_id ON ledger_entries(account_id);
CREATE INDEX IF NOT EXISTS idx_customers_external_id ON customers(external_id);

-- Seed Initial Data
INSERT INTO customers (id, external_id, full_name) 
VALUES ('de305d54-75b4-431b-adb2-eb6b9e546014', 'test_user', 'Jan Kowalski')
ON CONFLICT (external_id) DO NOTHING;

INSERT INTO accounts (id, customer_id, account_number, currency, balance, status)
VALUES ('de305d54-75b4-431b-adb2-eb6b9e546014', 'de305d54-75b4-431b-adb2-eb6b9e546014', 'PL12345678900000000012345678', 'PLN', 12500.50, 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO ledger_entries (account_id, type, amount, balance_after, description)
VALUES 
('de305d54-75b4-431b-adb2-eb6b9e546014', 'deposit', 10000.00, 10000.00, 'Wpłata początkowa'),
('de305d54-75b4-431b-adb2-eb6b9e546014', 'deposit', 2500.50, 12500.50, 'Premia świąteczna')
ON CONFLICT DO NOTHING;
