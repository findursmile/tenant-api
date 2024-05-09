# Table tenant
DEFINE TABLE tenant SCHEMALESS
    PERMISSIONS
        FOR create, update, delete, select where id = $auth.id;

DEFINE FIELD name ON TABLE tenant TYPE string;
DEFINE FIELD email ON TABLE tenant TYPE string
    VALUE string::lowercase($value)
    ASSERT string::is::email($value);
DEFINE FIELD password ON TABLE tenant TYPE string;
DEFINE FIELD mobile ON TABLE tenant TYPE string;
DEFINE FIELD country_code ON TABLE tenant TYPE string;
DEFINE FIELD status ON TABLE tenant TYPE string;

DEFINE FIELD created ON tenant VALUE time::now() READONLY TYPE datetime;
DEFINE FIELD updated ON tenant VALUE time::now() TYPE datetime;

DEFINE INDEX email ON tenant FIELDS email UNIQUE;

DEFINE SCOPE tenant SESSION 1d
	SIGNIN (
		SELECT * FROM tenant WHERE email = $email AND crypto::argon2::compare(password, $password)
	)
	SIGNUP (
		CREATE tenant CONTENT {
			name: $name,
			email: $email,
			mobile: $mobile,
			country_code: $country_code,
			status: $status,
			password: crypto::argon2::generate($password)
		}
	);



# Table event
DEFINE TABLE event SCHEMALESS
    PERMISSIONS
        FOR create, update, delete, select where tenant = $auth.id;

DEFINE FIELD name ON TABLE event TYPE string;
DEFINE FIELD cover_photo ON TABLE event TYPE option<string>;
DEFINE FIELD title ON TABLE event TYPE string;
DEFINE FIELD event_date ON TABLE event TYPE datetime;
DEFINE FIELD event_end_at ON TABLE event TYPE datetime;
DEFINE FIELD status ON TABLE event TYPE string;
DEFINE FIELD tenant ON TABLE event TYPE record<tenant>;

DEFINE FIELD created ON event VALUE time::now() READONLY;
DEFINE FIELD updated ON event VALUE time::now() TYPE datetime;


# RELATION owner_of  - Tenant owner_of event
DEFINE TABLE owner_of TYPE RELATION FROM tenant TO event;

DEFINE EVENT owner ON TABLE event WHEN $before == NONE THEN {
    LET $tenant = (select * from tenant where id = $after.tenant);
    RELATE $tenant->owner_of->$after;
};


# Table image
DEFINE TABLE image SCHEMALESS
    PERMISSIONS
        FOR create, update, delete, select where event.tenant = $auth.id;

DEFINE FIELD image_uri ON TABLE image TYPE string;
DEFINE FIELD status ON TABLE image TYPE string;

DEFINE FIELD created ON image VALUE time::now() READONLY TYPE datetime;
DEFINE FIELD updated ON image VALUE time::now() TYPE datetime;


# RELATION event_of  - Event event_of image
DEFINE TABLE event_of TYPE RELATION FROM event TO image;

# Table encoding
DEFINE TABLE face_encoding SCHEMALESS;
DEFINE FIELD encoding ON TABLE face_encoding TYPE array<decimal>;
DEFINE FIELD position ON TABLE face_encoding TYPE array<decimal>;

DEFINE FIELD created ON face_encoding VALUE time::now() READONLY TYPE datetime;
DEFINE FIELD updated ON face_encoding VALUE time::now() TYPE datetime;

# RELATION face_of  - image face_of face_encoding
DEFINE TABLE face_of TYPE RELATION FROM image TO face_encoding;

