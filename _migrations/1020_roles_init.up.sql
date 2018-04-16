INSERT INTO roles (name)
    VALUES ('admin'),
        ('upstream-admin'),
        ('lawmaker'),
        ('user'),
        ('service'),
        ('banned'),
        ('disabled');

INSERT INTO operations (name, description)
    VALUES 
    ('comment_create', 'Create comment'),
    ('comment_list', 'List comments'),
    ('comment_delete', 'Delete comment'),
    ('comment_update', 'Update comment'),
    ('law_create', 'Create law'),
    ('laws_create', 'Brach create laws'),
    ('law_delete', 'Delete law'),
    ('law_view', 'View law'),
    ('law_update', 'Update law'),
    ('repos_save', 'Admin save repos'),
    ('repos_load', 'Admin load repos'),
    ('user_create', 'Admin create user account'),
    ('user_delete', 'Admin delete user account'),
    ('user_view', 'Admin view user account'),
    ('user_update', 'Admin update user account'),
    ('upstream_create', 'Admin create upstream'),
    ('upstream_delete', 'Admin delete upstream'),
    ('upstream_view', 'Admin view upstream'),
    ('upstream_update', 'Admin update upstream'),
    ('vote_create', 'Create vote'),
    ('votes_create', 'Batch create votes');

INSERT INTO role_operations (role_id, operation_id)
    VALUES 
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'repos_load')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'repos_save')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'comment_create')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'user_create')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'user_delete')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'user_view')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'user_update')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'upstream_create')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'upstream_delete')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'upstream_view')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'upstream_update')),
    ((SELECT uid FROM roles WHERE name = 'admin'), (SELECT uid FROM operations WHERE name = 'votes_create')),
    ((SELECT uid FROM roles WHERE name = 'upstream-admin'), (SELECT uid FROM operations WHERE name = 'law_create')),
    ((SELECT uid FROM roles WHERE name = 'upstream-admin'), (SELECT uid FROM operations WHERE name = 'laws_create')),
    ((SELECT uid FROM roles WHERE name = 'upstream-admin'), (SELECT uid FROM operations WHERE name = 'user_create')),
    ((SELECT uid FROM roles WHERE name = 'upstream-admin'), (SELECT uid FROM operations WHERE name = 'user_update')),
    ((SELECT uid FROM roles WHERE name = 'upstream-admin'), (SELECT uid FROM operations WHERE name = 'upstream_view')),
    ((SELECT uid FROM roles WHERE name = 'upstream-admin'), (SELECT uid FROM operations WHERE name = 'upstream_update')),
    ((SELECT uid FROM roles WHERE name = 'upstream-admin'), (SELECT uid FROM operations WHERE name = 'votes_create'));
