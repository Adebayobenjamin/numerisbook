-- Drop indexes first
DROP INDEX IF EXISTS idx_audit_trails_customer_id ON audit_trails;
DROP INDEX IF EXISTS idx_audit_trails_invoice_id ON audit_trails;
DROP INDEX IF EXISTS idx_invoice_reminders_customer_id ON invoice_reminders;
DROP INDEX IF EXISTS idx_invoice_reminders_invoice_id ON invoice_reminders;
DROP INDEX IF EXISTS idx_payment_info_invoice_id ON payment_info;
DROP INDEX IF EXISTS idx_payments_invoice_id ON payments;
DROP INDEX IF EXISTS idx_invoice_items_invoice_id ON invoice_items;
DROP INDEX IF EXISTS idx_invoices_customer_id ON invoices;

-- Drop unique constraint
ALTER TABLE invoice_reminders 
DROP CONSTRAINT IF EXISTS uk_invoice_schedule;

-- Drop tables in reverse order of creation (to handle foreign key dependencies)
DROP TABLE IF EXISTS audit_trails;
DROP TABLE IF EXISTS invoice_reminders;
DROP TABLE IF EXISTS payment_info;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS senders;
