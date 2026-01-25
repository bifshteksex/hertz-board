-- Migration: Create assets table
-- This table stores metadata for uploaded files (images, documents, etc.)

CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    uploaded_by UUID NOT NULL REFERENCES users(id),

    -- File metadata
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL, -- File size in bytes

    -- URLs
    url TEXT NOT NULL, -- Full path to file in MinIO
    thumbnail_url TEXT, -- Path to thumbnail (for images)

    -- Image-specific metadata (nullable for non-images)
    width INTEGER,
    height INTEGER,

    -- Audit fields
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP -- Soft delete support
);

-- Indexes for performance
CREATE INDEX idx_assets_workspace_id ON assets(workspace_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_uploaded_by ON assets(uploaded_by);
CREATE INDEX idx_assets_content_type ON assets(content_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_created_at ON assets(created_at DESC);
CREATE INDEX idx_assets_deleted_at ON assets(deleted_at);

-- Comments for documentation
COMMENT ON TABLE assets IS 'Stores metadata for uploaded files (images, documents, etc.)';
COMMENT ON COLUMN assets.url IS 'Full path to file in MinIO storage';
COMMENT ON COLUMN assets.thumbnail_url IS 'Path to thumbnail image (for images only)';
COMMENT ON COLUMN assets.width IS 'Image width in pixels (for images only)';
COMMENT ON COLUMN assets.height IS 'Image height in pixels (for images only)';
