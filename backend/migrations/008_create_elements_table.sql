-- Create elements table for CRDT-based canvas elements
CREATE TABLE IF NOT EXISTS elements (
    id UUID PRIMARY KEY,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    content TEXT DEFAULT '',
    pos_x DOUBLE PRECISION DEFAULT 0,
    pos_y DOUBLE PRECISION DEFAULT 0,
    width DOUBLE PRECISION DEFAULT 100,
    height DOUBLE PRECISION DEFAULT 100,
    z_index INTEGER DEFAULT 0,
    rotation DOUBLE PRECISION DEFAULT 0,
    style JSONB DEFAULT '{}'::jsonb,
    version BIGINT NOT NULL DEFAULT 0, -- Lamport timestamp of last update
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    updated_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP -- Soft delete for tombstoning
);

-- Create indexes
CREATE INDEX idx_elements_workspace_id ON elements(workspace_id);
CREATE INDEX idx_elements_type ON elements(type);
CREATE INDEX idx_elements_z_index ON elements(z_index);
CREATE INDEX idx_elements_version ON elements(version);
CREATE INDEX idx_elements_deleted_at ON elements(deleted_at);

-- Composite index for queries
CREATE INDEX idx_elements_workspace_active ON elements(workspace_id, deleted_at) WHERE deleted_at IS NULL;

-- Comment on table
COMMENT ON TABLE elements IS 'Stores canvas elements with CRDT support';
COMMENT ON COLUMN elements.version IS 'Lamport timestamp for CRDT conflict resolution';
COMMENT ON COLUMN elements.deleted_at IS 'Tombstone timestamp for soft deletes (CRDT)';
