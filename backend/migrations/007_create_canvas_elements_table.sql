-- Migration: Create canvas_elements table
-- This table stores all canvas elements (text, shapes, images, drawings, etc.)
-- using JSONB for flexible element data structure

CREATE TABLE IF NOT EXISTS canvas_elements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,

    -- Element metadata
    element_type VARCHAR(50) NOT NULL, -- text, shape, image, drawing, sticky, list, connector, group
    element_data JSONB NOT NULL DEFAULT '{}', -- Flexible JSONB storage for element properties

    -- Common properties (duplicated for query optimization)
    z_index INTEGER NOT NULL DEFAULT 0,
    parent_id UUID REFERENCES canvas_elements(id) ON DELETE SET NULL, -- For grouping

    -- Audit fields
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP -- Soft delete support
);

-- Indexes for performance optimization
CREATE INDEX idx_canvas_elements_workspace_id ON canvas_elements(workspace_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_canvas_elements_type ON canvas_elements(element_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_canvas_elements_z_index ON canvas_elements(workspace_id, z_index) WHERE deleted_at IS NULL;
CREATE INDEX idx_canvas_elements_parent_id ON canvas_elements(parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_canvas_elements_created_by ON canvas_elements(created_by);
CREATE INDEX idx_canvas_elements_deleted_at ON canvas_elements(deleted_at);

-- JSONB indexes for efficient queries
CREATE INDEX idx_canvas_elements_data_gin ON canvas_elements USING GIN (element_data);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_canvas_elements_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_canvas_elements_updated_at
    BEFORE UPDATE ON canvas_elements
    FOR EACH ROW
    EXECUTE FUNCTION update_canvas_elements_updated_at();

-- Comments for documentation
COMMENT ON TABLE canvas_elements IS 'Stores all canvas elements with flexible JSONB data structure';
COMMENT ON COLUMN canvas_elements.element_type IS 'Type of element: text, shape, image, drawing, sticky, list, connector, group';
COMMENT ON COLUMN canvas_elements.element_data IS 'JSONB containing all element-specific properties (position, size, style, content, etc.)';
COMMENT ON COLUMN canvas_elements.z_index IS 'Layer order (higher = on top)';
COMMENT ON COLUMN canvas_elements.parent_id IS 'Parent element ID for grouped elements';
