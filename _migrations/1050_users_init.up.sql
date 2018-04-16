INSERT INTO providers (ident)
    VALUES ('welaw'), ('google'), ('amazon'), ('microsoft');

-- TODO  replace with uuid_nil()

INSERT INTO users (username, full_name, email, biography, picture_url, upstream, provider)
    VALUES (
        'master',
        '',
        '',
        '',
        '',
        null,
        (SELECT uid FROM providers WHERE ident = 'welaw')
);

INSERT INTO users (username, full_name, email, biography, picture_url, upstream, provider)
    VALUES (
        'congress-admin',
        'US Congress Administrator',
        'congress@welaw.org',
        '',
        '/assets/congress.png',
        (SELECT uid FROM upstreams WHERE ident = 'congress'),
        (SELECT uid FROM providers WHERE ident = 'welaw')
);

INSERT INTO users (username, full_name, email, biography, picture_url, upstream, provider)
    VALUES (
        'parliament-admin',
        'UK Parliament Administrator',
        'parliament@welaw.org',
        '',
        '/assets/parliament.png',
        (SELECT uid FROM upstreams WHERE ident = 'parliament'),
        (SELECT uid FROM providers WHERE ident = 'welaw')
);

INSERT INTO user_roles (user_id, role_id, scope)
    VALUES ((SELECT uid FROM users WHERE username = 'master'), (SELECT uid FROM roles WHERE name = 'service'), ''),
        (
            (SELECT uid FROM users WHERE username = 'congress-admin'),
            (SELECT uid FROM roles WHERE name = 'upstream-admin'),
            'congress'
        ),
        (
            (SELECT uid FROM users WHERE username = 'parliament-admin'),
            (SELECT uid FROM roles WHERE name = 'upstream-admin'),
            'parliament'
        );
