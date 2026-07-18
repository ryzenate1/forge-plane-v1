-- Add indexes for performance optimization

-- Index for server name lookups
CREATE INDEX IF NOT EXISTS idx_servers_name ON servers(name);

-- Index for server status queries
CREATE INDEX IF NOT EXISTS idx_servers_status ON servers(status);

-- Index for server node lookups
CREATE INDEX IF NOT EXISTS idx_servers_node_id ON servers(node_id);

-- Composite index for per-node status queries
CREATE INDEX IF NOT EXISTS idx_servers_node_status ON servers(node_id, status);

-- Index for node name lookups (already UNIQUE, but index aids search patterns)
CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name);

-- Index for node status
CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status);

-- Index for node IP address lookups
CREATE INDEX IF NOT EXISTS idx_nodes_ip_address ON nodes(ip_address);

-- Index for backup server lookups
CREATE INDEX IF NOT EXISTS idx_backups_server_id ON backups(server_id);

-- Composite index for per-server backup listings ordered by time
CREATE INDEX IF NOT EXISTS idx_backups_server_created ON backups(server_id, created_at);

-- Index for backup timestamp queries (e.g. retention cleanup)
CREATE INDEX IF NOT EXISTS idx_backups_created_at ON backups(created_at);

-- Index for backup status queries
CREATE INDEX IF NOT EXISTS idx_backups_status ON backups(status);