-- Migration 03: Multi-Tenancy Support
-- Tenants table (one per operator company)
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id) UNIQUE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(50) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SUSPENDED', 'TRIAL')),
    plan VARCHAR(50) DEFAULT 'STARTER' CHECK (plan IN ('STARTER', 'PROFESSIONAL', 'ENTERPRISE')),
    max_users INT DEFAULT 10,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Per-tenant feature/microservice flags
CREATE TABLE IF NOT EXISTS tenant_features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE,
    feature_key VARCHAR(100) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, feature_key)
);

-- SuperAdmin support: make company_id nullable on users
ALTER TABLE users ALTER COLUMN company_id DROP NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_super_admin BOOLEAN DEFAULT FALSE;

-- Seed the initial tenant for the existing admin company
INSERT INTO tenants (company_id, name, slug, status, plan)
SELECT id, name, LOWER(REPLACE(name, ' ', '-')), 'ACTIVE', 'PROFESSIONAL'
FROM companies
WHERE id = 'da0e254c-ac52-4d34-bc7c-03ff6cdc6b18'
ON CONFLICT (company_id) DO NOTHING;

-- Seed all features as enabled for existing tenant
INSERT INTO tenant_features (tenant_id, feature_key, is_enabled)
SELECT t.id, f.key, true
FROM tenants t
CROSS JOIN (VALUES
    ('dispatch'),
    ('outsource_marketplace'),
    ('fleet_management'),
    ('team_management'),
    ('invoicing'),
    ('partner_portal'),
    ('analytics')
) AS f(key)
ON CONFLICT (tenant_id, feature_key) DO NOTHING;

-- Seed SuperAdmin user
INSERT INTO users (company_id, role, email, password_hash, name, is_super_admin)
VALUES (NULL, 'SUPER_ADMIN', 'superadmin@limoxlink.com', '$2a$10$Nqwszx9V96zm9wCeFtGIreYzMf0ojeUMpYV0PHhS0x4qLi2/avToG', 'Super Admin', true)
ON CONFLICT (email) DO UPDATE
SET is_super_admin = true, role = 'SUPER_ADMIN';
