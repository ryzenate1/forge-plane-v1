PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS _servers_migrate;
CREATE TABLE _servers_migrate (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    node_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'stopped' CHECK (status IN ('starting', 'running', 'stopping', 'stopped', 'crashed')),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE
);
INSERT INTO _servers_migrate
SELECT id, name, node_id,
       CASE WHEN status IN ('starting','running','stopping','stopped','crashed') THEN status ELSE 'stopped' END,
       created_at, updated_at
FROM servers;
DROP TABLE IF EXISTS servers;
ALTER TABLE _servers_migrate RENAME TO servers;

DROP TABLE IF EXISTS _nodes_migrate;
CREATE TABLE _nodes_migrate (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'offline' CHECK (status IN ('online', 'offline', 'maintenance')),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
INSERT INTO _nodes_migrate
SELECT id, name, ip_address,
       CASE WHEN status IN ('online','offline','maintenance') THEN status ELSE 'offline' END,
       created_at, updated_at
FROM nodes;
DROP TABLE IF EXISTS nodes;
ALTER TABLE _nodes_migrate RENAME TO nodes;

DROP TABLE IF EXISTS _backups_migrate;
CREATE TABLE _backups_migrate (
    id TEXT PRIMARY KEY,
    server_id TEXT NOT NULL,
    path TEXT NOT NULL,
    size_bytes INTEGER NOT NULL DEFAULT 0,
    checksum TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending',
    error_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE
);
INSERT INTO _backups_migrate (id, server_id, path, created_at)
SELECT id, server_id, path, created_at FROM backups;
DROP TABLE IF EXISTS backups;
ALTER TABLE _backups_migrate RENAME TO backups;

PRAGMA foreign_keys = ON;
