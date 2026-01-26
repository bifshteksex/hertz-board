-- Create operations table for CRDT synchronization
CREATE TABLE IF NOT EXISTS operations (
    id UUID PRIMARY KEY,
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    element_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    op_type VARCHAR(20) NOT NULL CHECK (op_type IN ('create', 'update', 'delete', 'move')),
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    timestamp BIGINT NOT NULL, -- Lamport timestamp
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient queries
CREATE INDEX idx_operations_workspace_id ON operations(workspace_id);
CREATE INDEX idx_operations_element_id ON operations(element_id);
CREATE INDEX idx_operations_user_id ON operations(user_id);
CREATE INDEX idx_operations_timestamp ON operations(timestamp);
CREATE INDEX idx_operations_workspace_timestamp ON operations(workspace_id, timestamp);
CREATE INDEX idx_operations_created_at ON operations(created_at);

-- Create composite index for sync queries
CREATE INDEX idx_operations_sync ON operations(workspace_id, user_id, timestamp);

-- Comment on table
COMMENT ON TABLE operations IS 'Stores CRDT operations for real-time synchronization';
COMMENT ON COLUMN operations.timestamp IS 'Lamport timestamp for operation ordering';
COMMENT ON COLUMN operations.data IS 'Operation-specific data (element properties for create/update)';
