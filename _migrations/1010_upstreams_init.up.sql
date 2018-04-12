INSERT INTO upstreams (
    ident,
    upstream_name,
    upstream_description,
    name,
    description,
    url,
    metadata,
    geo_coords
) VALUES (
    'usa',
    'US Congress',
    'Laws and votes from the US House of Representatives and the US Senate pulled from https://gpo.gov and https://www.congress.gov.',
    'The United States Congress',
    'The Senate and the House of Representatives.',
    'https:/www.congress.gov',
    '{
        "api_url":"https://www.gpo.gov/fdsys/",
        "wikipedia":"https://en.wikipedia.org/wiki/United_States_Congress",
        "twitter":"https://twitter.com/congressdotgov",
        "facebook":"https://www.facebook.com/US-Congress-133668506703823/"
    }'::jsonb,
    POINT(38.889931, -77.009003)
);

INSERT INTO upstream_tags (
    upstream_id,
    ident,
    name,
    ranking,
    number_type,
    description
) VALUES (
    (SELECT uid FROM upstreams WHERE ident = 'usa'),
    'chamber',
    'Chamber',
    0,
    true,
    'The House and Senate'
), (
    (SELECT uid FROM upstreams WHERE ident = 'usa'),
    'congress',
    'Congress',
    1,
    true,
    'Congress number'
), (
    (SELECT uid FROM upstreams WHERE ident = 'usa'),
    'version',
    'Bill Version',
    4,
    false,
    'The Bill Version acronym'
);

INSERT INTO upstream_groups (
    upstream_id,
    ident,
    name,
    description
) VALUES (
    (SELECT uid FROM upstreams WHERE ident = 'usa'),
    'house',
    'The House of Representatives',
    ''
), (
    (SELECT uid FROM upstreams WHERE ident = 'usa'),
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
    geo_coords
) VALUES (
    'uk',
    'UK Parliament',
    'Laws from parliament pulled from http://api.data.parliament.uk/',
    'The United Kingdom Parliament',
    'The House of Lords and House of Commons.',
    'https://www.parliament.uk/',
    '{"api_url":"http://www.data.parliament.uk/"}'::jsonb,
    POINT(51.4900, -0.2200)
);

INSERT INTO upstream_tags (
    upstream_id,
    ident,
    name,
    ranking,
    number_type,
    description
) VALUES (
    (SELECT uid FROM upstreams WHERE ident = 'uk'),
    'house',
    'House',
    0,
    false,
    'The House of Lords and House of Commons'
), (
    (SELECT uid FROM upstreams WHERE ident = 'uk'),
    'session',
    'Session',
    2,
    false,
    'Session number'
), (
    (SELECT uid FROM upstreams WHERE ident = 'uk'),
    'version',
    'Bill Version',
    4,
    false,
    'The version of the bill'
);

INSERT INTO upstream_groups (
    upstream_id,
    ident,
    name,
    description
) VALUES (
    (SELECT uid FROM upstreams WHERE ident = 'uk'),
    'lords',
    'The House of Lords',
    'The House of Lords'
), (
    (SELECT uid FROM upstreams WHERE ident = 'uk'),
    'commons',
    'The House of Commons',
    'The House of Commons'
);
