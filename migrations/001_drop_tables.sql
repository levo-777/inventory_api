-- Migration 001: Drop existing tables if they exist
-- This migration drops the items table to ensure a clean start

DROP TABLE IF EXISTS items CASCADE;
