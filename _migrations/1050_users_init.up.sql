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
        'usa-admin',
        'USA Administrator',
        'usa@welaw.org',
        '',
        '/assets/usa.png',
        (SELECT uid FROM upstreams WHERE ident = 'usa'),
        (SELECT uid FROM providers WHERE ident = 'welaw')
);

INSERT INTO users (username, full_name, email, biography, picture_url, upstream, provider)
    VALUES (
        'uk-admin',
        'UK Administrator',
        'uk@welaw.org',
        '',
        '/assets/uk.png',
        (SELECT uid FROM upstreams WHERE ident = 'uk'),
        (SELECT uid FROM providers WHERE ident = 'welaw')
);

INSERT INTO user_roles (user_id, role_id)
    VALUES ((SELECT uid FROM users WHERE username = 'master'), (SELECT uid FROM roles WHERE name = 'service')),
        ((SELECT uid FROM users WHERE username = 'usa-admin'), (SELECT uid FROM roles WHERE name = 'upstream-admin')),
        ((SELECT uid FROM users WHERE username = 'uk-admin'), (SELECT uid FROM roles WHERE name = 'upstream-admin'));
