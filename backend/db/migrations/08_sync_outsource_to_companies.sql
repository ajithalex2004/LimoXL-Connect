-- Migration 08: Sync outsource_companies to companies table
-- This ensures all outsource partners can be selected for direct assignment in dispatch

INSERT INTO companies (id, name, type, verified, created_at, updated_at)
SELECT 
    id, 
    name, 
    'SUPPLY'::VARCHAR(50), 
    true, 
    created_at, 
    updated_at
FROM outsource_companies
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    type = 'SUPPLY',
    verified = true;

-- Also ensure any partners with 'BOTH' type in companies are preserved if they exist in outsource_companies
UPDATE companies c
SET type = 'BOTH'
FROM outsource_companies oc
WHERE c.id = oc.id AND (c.type = 'BOTH' OR c.type = 'SUPPLY');
