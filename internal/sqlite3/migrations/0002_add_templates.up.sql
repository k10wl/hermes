DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS template_dependencies;

CREATE TABLE templates (
    id INTEGER PRIMARY KEY,
    name TEXT,
    template TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE template_dependencies (
    core_id INTEGER,
    linked_id INTEGER,
    FOREIGN KEY (core_id) REFERENCES templates(id),
    FOREIGN KEY (linked_id) REFERENCES templates(id)
);
