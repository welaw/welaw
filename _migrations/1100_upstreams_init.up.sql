INSERT INTO upstreams (
    ident,
    upstream_name,
    upstream_description,
    name,
    description,
    url,
    metadata,
    geo_coords,
    user_id
) VALUES (
    'congress',
    'US Congress',
    'Laws and votes from the US House of Representatives and the US Senate, fetched from https://gpo.gov and https://www.congress.gov.',
    'The United States Congress',
    'The Senate and the House of Representatives.',
    'https:/www.congress.gov',
    '{
        "api_url":"https://www.gpo.gov/fdsys/",
        "wikipedia":"https://en.wikipedia.org/wiki/United_States_Congress",
        "twitter":"https://twitter.com/congressdotgov",
        "facebook":"https://www.facebook.com/US-Congress-133668506703823/"
    }'::jsonb,
    POINT(38.889931, -77.009003),
    (SELECT uid FROM users WHERE username = 'congress-admin')
);

INSERT INTO upstream_groups (
    upstream_id,
    ident,
    name,
    description
) VALUES (
    (SELECT uid FROM upstreams WHERE ident = 'congress'),
    'house',
    'The House of Representatives',
    ''
), (
    (SELECT uid FROM upstreams WHERE ident = 'congress'),
    'senate',
    'The Senate',
    ''
);

INSERT INTO upstreams (
    ident,
    upstream_name,
    upstream_description,
    name,
    description,
    url,
    metadata,
    geo_coords,
    user_id
) VALUES (
    'parliament',
    'UK Parliament',
    'Laws from the UK Parliament, fetched from http://api.data.parliament.uk/',
    'The United Kingdom Parliament',
    'The House of Lords and House of Commons.',
    'https://www.parliament.uk/',
    '{"api_url":"http://www.data.parliament.uk/"}'::jsonb,
    POINT(51.4900, -0.2200),
    (SELECT uid FROM users WHERE username = 'parliament-admin')
);

INSERT INTO upstream_groups (
    upstream_id,
    ident,
    name,
    description
) VALUES (
    (SELECT uid FROM upstreams WHERE ident = 'parliament'),
    'lords',
    'The House of Lords',
    'The House of Lords'
), (
    (SELECT uid FROM upstreams WHERE ident = 'parliament'),
    'commons',
    'The House of Commons',
    'The House of Commons'
);
