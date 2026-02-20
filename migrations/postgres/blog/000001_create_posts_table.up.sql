CREATE TABLE posts (
    id         UUID PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    content    TEXT NOT NULL,

    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    deleted_by UUID,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX idx_posts_name_created_by ON posts (name, created_by) WHERE deleted_at IS NULL;