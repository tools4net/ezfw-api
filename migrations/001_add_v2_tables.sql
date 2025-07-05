-- Migration script for V2 API: Add nodes and service_instances tables
-- This script adds the new V2 tables while preserving existing data

-- Create nodes table
CREATE TABLE IF NOT EXISTS nodes (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    hostname TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    port INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'inactive',
    os_info TEXT,
    agent_info TEXT,
    tags TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    last_seen DATETIME
);

-- Create service_instances table
CREATE TABLE IF NOT EXISTS service_instances (
    id TEXT PRIMARY KEY,
    node_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    service_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'stopped',
    port INTEGER NOT NULL,
    protocol TEXT NOT NULL,
    config TEXT,
    tags TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    last_seen DATETIME,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status);
CREATE INDEX IF NOT EXISTS idx_nodes_hostname ON nodes(hostname);
CREATE INDEX IF NOT EXISTS idx_service_instances_node_id ON service_instances(node_id);
CREATE INDEX IF NOT EXISTS idx_service_instances_type ON service_instances(service_type);
CREATE INDEX IF NOT EXISTS idx_service_instances_status ON service_instances(status);

-- Optional: Create a default localhost node for development/testing
-- Uncomment the following lines if you want a default node
/*
INSERT OR IGNORE INTO nodes (
    id, name, description, hostname, ip_address, port, status, 
    tags, created_at, updated_at
) VALUES (
    'node-localhost-default',
    'Localhost Development Node',
    'Default development node running on localhost',
    'localhost',
    '127.0.0.1',
    22,
    'active',
    '["development", "localhost"]',
    datetime('now'),
    datetime('now')
);
*/