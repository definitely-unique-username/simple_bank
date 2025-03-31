DROP INDEX IF EXISTS idx_accounts_owner;
DROP INDEX IF EXISTS idx_entries_account_id;
DROP INDEX IF EXISTS idx_transfers_from_account_id;
DROP INDEX IF EXISTS idx_transfers_to_account_id;
DROP INDEX IF EXISTS idx_transfers_from_to_account_id;

DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS accounts;

