-- Migration: Create canvas_snapshots table
-- This table stores version snapshots of the entire canvas state

CREATE TABLE IF NOT EXISTS canvas_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,

    -- Snapshot metadata
    version INTEGER NOT NULL, -- Sequential version number
    description TEXT, -- Optional description of changes

    -- Snapshot data (full canvas state)
    snapshot_data JSONB NOT NULL, -- Complete serialized canvas state
    element_count INTEGER NOT NULL DEFAULT 0, -- Number of elements in this snapshot

    -- Audit fields
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Ensure unique version numbers per workspace
    UNIQUE(workspace_id, version)
);

-- Indexes for performance
CREATE INDEX idx_canvas_snapshots_workspace_id ON canvas_snapshots(workspace_id);
CREATE INDEX idx_canvas_snapshots_version ON canvas_snapshots(workspace_id, version DESC);
CREATE INDEX idx_canvas_snapshots_created_at ON canvas_snapshots(workspace_id, created_at DESC);
CREATE INDEX idx_canvas_snapshots_created_by ON canvas_snapshots(created_by);

-- JSONB index for efficient queries
CREATE INDEX idx_canvas_snapshots_data_gin ON canvas_snapshots USING GIN (snapshot_data);

-- Function to get next version number for a workspace
CREATE OR REPLACE FUNCTION get_next_snapshot_version(workspace_uuid UUID)
RETURNS INTEGER AS $$
DECLARE
    next_version INTEGER;
BEGIN
    SELECT COALESCE(MAX(version), 0) + 1
    INTO next_version
    FROM canvas_snapshots
    WHERE workspace_id = workspace_uuid;

    RETURN next_version;
END;
$$ LANGUAGE plpgsql;

-- Comments for documentation
COMMENT ON TABLE canvas_snapshots IS 'Stores version snapshots of canvas state for versioning and restore';
COMMENT ON COLUMN canvas_snapshots.version IS 'Sequential version number starting from 1';
COMMENT ON COLUMN canvas_snapshots.snapshot_data IS 'Complete serialized canvas state (all elements)';
COMMENT ON COLUMN canvas_snapshots.element_count IS 'Number of elements in this snapshot for quick reference';
