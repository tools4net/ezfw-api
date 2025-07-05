-- V1 to V2 Migration Script
-- This script migrates existing SingBox and Xray configurations to the new V2 architecture
-- It creates default nodes and converts configurations to service instances

-- Step 1: Create a default localhost node if it doesn't exist
INSERT OR IGNORE INTO nodes (
    id, name, description, hostname, ip_address, port, status, 
    tags, created_at, updated_at
) VALUES (
    'node-localhost-default',
    'Localhost Migration Node',
    'Default node created during V1 to V2 migration for existing configurations',
    'localhost',
    '127.0.0.1',
    22,
    'active',
    '["migration", "localhost", "default"]',
    datetime('now'),
    datetime('now')
);

-- Step 2: Migrate existing Xray configurations to service instances
INSERT OR IGNORE INTO service_instances (
    id,
    node_id,
    name,
    description,
    service_type,
    status,
    port,
    protocol,
    config,
    tags,
    created_at,
    updated_at
)
SELECT 
    'service-xray-' || x.id as id,
    'node-localhost-default' as node_id,
    COALESCE(x.name, 'Migrated Xray Config ' || substr(x.id, 1, 8)) as name,
    COALESCE(x.description, 'Migrated from V1 Xray configuration') as description,
    'xray' as service_type,
    'stopped' as status,
    COALESCE(
        -- Try to extract port from first inbound
        CASE 
            WHEN x.inbounds IS NOT NULL AND x.inbounds != '' AND x.inbounds != 'null' THEN
                CAST(json_extract(x.inbounds, '$[0].port') AS INTEGER)
            ELSE 443
        END,
        443
    ) as port,
    'tcp' as protocol,
    json_object(
        'id', x.id,
        'name', x.name,
        'description', x.description,
        'log', CASE WHEN x.log_config IS NOT NULL AND x.log_config != '' AND x.log_config != 'null' THEN json(x.log_config) ELSE NULL END,
        'api', CASE WHEN x.api_config IS NOT NULL AND x.api_config != '' AND x.api_config != 'null' THEN json(x.api_config) ELSE NULL END,
        'dns', CASE WHEN x.dns_config IS NOT NULL AND x.dns_config != '' AND x.dns_config != 'null' THEN json(x.dns_config) ELSE NULL END,
        'routing', CASE WHEN x.routing_config IS NOT NULL AND x.routing_config != '' AND x.routing_config != 'null' THEN json(x.routing_config) ELSE NULL END,
        'policy', CASE WHEN x.policy_config IS NOT NULL AND x.policy_config != '' AND x.policy_config != 'null' THEN json(x.policy_config) ELSE NULL END,
        'inbounds', CASE WHEN x.inbounds IS NOT NULL AND x.inbounds != '' AND x.inbounds != 'null' THEN json(x.inbounds) ELSE json('[]') END,
        'outbounds', CASE WHEN x.outbounds IS NOT NULL AND x.outbounds != '' AND x.outbounds != 'null' THEN json(x.outbounds) ELSE json('[]') END,
        'transport', CASE WHEN x.transport_config IS NOT NULL AND x.transport_config != '' AND x.transport_config != 'null' THEN json(x.transport_config) ELSE NULL END,
        'stats', CASE WHEN x.stats_config IS NOT NULL AND x.stats_config != '' AND x.stats_config != 'null' THEN json(x.stats_config) ELSE NULL END,
        'reverse', CASE WHEN x.reverse_config IS NOT NULL AND x.reverse_config != '' AND x.reverse_config != 'null' THEN json(x.reverse_config) ELSE NULL END,
        'fakedns', CASE WHEN x.fakedns_config IS NOT NULL AND x.fakedns_config != '' AND x.fakedns_config != 'null' THEN json(x.fakedns_config) ELSE NULL END,
        'metrics', CASE WHEN x.metrics_config IS NOT NULL AND x.metrics_config != '' AND x.metrics_config != 'null' THEN json(x.metrics_config) ELSE NULL END,
        'observatory', CASE WHEN x.observatory_config IS NOT NULL AND x.observatory_config != '' AND x.observatory_config != 'null' THEN json(x.observatory_config) ELSE NULL END,
        'burstObservatory', CASE WHEN x.burst_observatory_config IS NOT NULL AND x.burst_observatory_config != '' AND x.burst_observatory_config != 'null' THEN json(x.burst_observatory_config) ELSE NULL END,
        'created_at', x.created_at,
        'updated_at', x.updated_at
    ) as config,
    '["migrated", "xray", "v1"]' as tags,
    COALESCE(x.created_at, datetime('now')) as created_at,
    COALESCE(x.updated_at, datetime('now')) as updated_at
FROM xray_configs x
WHERE NOT EXISTS (
    SELECT 1 FROM service_instances s 
    WHERE s.id = 'service-xray-' || x.id
);

-- Step 3: Migrate existing SingBox configurations to service instances
INSERT OR IGNORE INTO service_instances (
    id,
    node_id,
    name,
    description,
    service_type,
    status,
    port,
    protocol,
    config,
    tags,
    created_at,
    updated_at
)
SELECT 
    'service-singbox-' || s.id as id,
    'node-localhost-default' as node_id,
    COALESCE(s.name, 'Migrated SingBox Config ' || substr(s.id, 1, 8)) as name,
    COALESCE(s.description, 'Migrated from V1 SingBox configuration') as description,
    'singbox' as service_type,
    'stopped' as status,
    COALESCE(
        -- Try to extract port from first inbound
        CASE 
            WHEN s.inbounds IS NOT NULL AND s.inbounds != '' AND s.inbounds != 'null' THEN
                CAST(json_extract(s.inbounds, '$[0].listen_port') AS INTEGER)
            ELSE 443
        END,
        443
    ) as port,
    'tcp' as protocol,
    json_object(
        'id', s.id,
        'name', s.name,
        'description', s.description,
        'log', CASE WHEN s.log_config IS NOT NULL AND s.log_config != '' AND s.log_config != 'null' THEN json(s.log_config) ELSE NULL END,
        'dns', CASE WHEN s.dns_config IS NOT NULL AND s.dns_config != '' AND s.dns_config != 'null' THEN json(s.dns_config) ELSE NULL END,
        'ntp', CASE WHEN s.ntp_config IS NOT NULL AND s.ntp_config != '' AND s.ntp_config != 'null' THEN json(s.ntp_config) ELSE NULL END,
        'inbounds', CASE WHEN s.inbounds IS NOT NULL AND s.inbounds != '' AND s.inbounds != 'null' THEN json(s.inbounds) ELSE json('[]') END,
        'outbounds', CASE WHEN s.outbounds IS NOT NULL AND s.outbounds != '' AND s.outbounds != 'null' THEN json(s.outbounds) ELSE json('[]') END,
        'route', CASE WHEN s.route_config IS NOT NULL AND s.route_config != '' AND s.route_config != 'null' THEN json(s.route_config) ELSE NULL END,
        'experimental', CASE WHEN s.experimental_config IS NOT NULL AND s.experimental_config != '' AND s.experimental_config != 'null' THEN json(s.experimental_config) ELSE NULL END,
        'services', CASE WHEN s.services_config IS NOT NULL AND s.services_config != '' AND s.services_config != 'null' THEN json(s.services_config) ELSE NULL END,
        'endpoints', CASE WHEN s.endpoints_config IS NOT NULL AND s.endpoints_config != '' AND s.endpoints_config != 'null' THEN json(s.endpoints_config) ELSE NULL END,
        'certificate', CASE WHEN s.certificate_config IS NOT NULL AND s.certificate_config != '' AND s.certificate_config != 'null' THEN json(s.certificate_config) ELSE NULL END,
        'created_at', s.created_at,
        'updated_at', s.updated_at
    ) as config,
    '["migrated", "singbox", "v1"]' as tags,
    COALESCE(s.created_at, datetime('now')) as created_at,
    COALESCE(s.updated_at, datetime('now')) as updated_at
FROM singbox_configs s
WHERE NOT EXISTS (
    SELECT 1 FROM service_instances si 
    WHERE si.id = 'service-singbox-' || s.id
);

-- Step 4: Update the default node's last_seen timestamp
UPDATE nodes 
SET last_seen = datetime('now'), updated_at = datetime('now')
WHERE id = 'node-localhost-default';

-- Step 5: Create migration summary view (temporary, for verification)
CREATE TEMPORARY VIEW migration_summary AS
SELECT 
    'Migration Summary' as operation,
    (
        SELECT COUNT(*) FROM xray_configs
    ) as total_xray_configs,
    (
        SELECT COUNT(*) FROM singbox_configs  
    ) as total_singbox_configs,
    (
        SELECT COUNT(*) FROM service_instances WHERE service_type = 'xray' AND tags LIKE '%migrated%'
    ) as migrated_xray_services,
    (
        SELECT COUNT(*) FROM service_instances WHERE service_type = 'singbox' AND tags LIKE '%migrated%'
    ) as migrated_singbox_services,
    (
        SELECT COUNT(*) FROM nodes WHERE tags LIKE '%migration%'
    ) as migration_nodes_created;

-- Display migration summary
SELECT * FROM migration_summary;