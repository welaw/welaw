CREATE MATERIALIZED VIEW search_index AS 
    SELECT versions.uid as uid,
        setweight(to_tsvector('english'::regconfig, author.username), 'A') ||
        setweight(to_tsvector('english'::regconfig, author.full_name), 'A') ||
        setweight(to_tsvector('english'::regconfig, laws.title), 'A') ||
        setweight(to_tsvector('english'::regconfig, laws.short_title), 'A') ||
        setweight(to_tsvector('english'::regconfig, users.username), 'A') ||
        setweight(to_tsvector('english'::regconfig, users.full_name), 'A') as doc
    FROM versions
    INNER JOIN branches ON branches.uid = versions.branch_id
    INNER JOIN users ON users.uid = versions.user_id
    INNER JOIN laws ON laws.uid = branches.law_id
    INNER JOIN (
        SELECT branches.uid,
            username,
            full_name
        FROM users
        INNER JOIN branches ON branches.user_id = users.uid
        INNER JOIN versions ON versions.branch_id = branches.uid
        WHERE branches.name = (SELECT uid FROM users WHERE username = 'master' AND deleted_at IS NULL)::varchar
        ORDER BY versions.number
    ) author ON author.uid = branches.uid;

CREATE INDEX idx_fts_search ON search_index USING gin(doc);

CREATE MATERIALIZED VIEW word_index AS 
    SELECT word FROM ts_stat(
        'SELECT to_tsvector(''simple'', laws.title) ||
            to_tsvector(''simple'', laws.short_title) ||
            to_tsvector(''simple'', users.username) ||
            to_tsvector(''simple'', users.full_name)
        FROM versions
        INNER JOIN branches ON branches.uid = versions.branch_id
        INNER JOIN users ON users.uid = versions.user_id
        INNER JOIN laws ON laws.uid = branches.law_id');

CREATE INDEX words_idx ON word_index USING gin(word gin_trgm_ops);
