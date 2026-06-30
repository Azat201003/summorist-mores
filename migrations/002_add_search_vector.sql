ALTER TABLE metas ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
    to_tsvector('simple', title || ' ' || description)
) STORED;

CREATE INDEX idx_metas_search ON metas USING GIN (search_vector);

