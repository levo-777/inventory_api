-- Migration 002: Create the items table
-- This migration creates the items table with the proper schema

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the items table
CREATE TABLE items (
    -- id is the primary key for the table (UUID)
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- name is the name of the item
    name VARCHAR(255) NOT NULL,
    -- stock is the current stock level
    stock INTEGER NOT NULL DEFAULT 0,
    -- price is the price of the item
    price DECIMAL(10, 2) NOT NULL,
    -- created_at is the timestamp when the item was created
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- updated_at is the timestamp when the item was last updated
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- deleted_at is used for soft deletes
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_items_name ON items (name);
CREATE INDEX IF NOT EXISTS idx_items_stock ON items (stock);
CREATE INDEX IF NOT EXISTS idx_items_price ON items (price);
CREATE INDEX IF NOT EXISTS idx_items_created_at ON items (created_at);
CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items (deleted_at);
