INSERT INTO providers (ident) VALUES ('welaw'), ('google');

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
        'US Congress Admin',
        'congress@welaw.org',
        '',
        '/assets/congress.png',
        (SELECT uid FROM upstreams WHERE ident = 'congress'),
        (SELECT uid FROM providers WHERE ident = 'welaw')
);

INSERT INTO users (username, full_name, email, biography, picture_url, upstream, provider)
    VALUES (
        'parliament-admin',
        'UK Parliament Admin',
        'parliament@welaw.org',
        '',
        '/assets/uk.png',
        (SELECT uid FROM upstreams WHERE ident = 'uk'),
        (SELECT uid FROM providers WHERE ident = 'welaw')
);

INSERT INTO user_roles (user_id, role_id)
    VALUES ((SELECT uid FROM users WHERE username = 'master'), (SELECT uid FROM roles WHERE name = 'service')),
        ((SELECT uid FROM users WHERE username = 'congress-admin'), (SELECT uid FROM roles WHERE name = 'upstream-admin')),
        ((SELECT uid FROM users WHERE username = 'parliament-admin'), (SELECT uid FROM roles WHERE name = 'upstream-admin'));
