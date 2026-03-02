-- Migration 04: Add rfq_number support to trips table
ALTER TABLE trips ADD COLUMN IF NOT EXISTS rfq_number VARCHAR(100);

-- Create sequence for RFQ numbers starting from 1000
CREATE SEQUENCE IF NOT EXISTS rfq_id_seq START 1000;
