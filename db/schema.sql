CREATE SCHEMA r1;

SET search_path TO r1, public;

CREATE TABLE services
(
    service_id        UUID PRIMARY KEY,
    name              VARCHAR NOT NULL,
    description       TEXT,
    owner_info        TEXT,
    industry_category VARCHAR,
    client_rating     FLOAT,
    created_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE service_instances
(
    instance_id   UUID PRIMARY KEY,
    service_id    UUID REFERENCES services (service_id),
    version       VARCHAR NOT NULL,
    host          VARCHAR NOT NULL,
    port          INTEGER NOT NULL,
    url           VARCHAR NOT NULL,
    api_spec      TEXT,
    latitude      FLOAT,
    longitude     FLOAT,
    health_status VARCHAR,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_checked  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE service_instance_history
(
    history_id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    instance_id UUID PRIMARY KEY,
    service_id  UUID REFERENCES services (service_id),
    version     VARCHAR NOT NULL,
    url         VARCHAR NOT NULL,
    metrics     JSON,
    started_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    stopped_at   TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE registries
(
    registry_id UUID PRIMARY KEY,
    url         VARCHAR NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE service_reviews
(
    review_id  UUID PRIMARY KEY,
    service_id UUID REFERENCES services (service_id),
    rating     INTEGER NOT NULL,
    review     TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE registry_group
(
    group_id    UUID PRIMARY KEY,
    registry_id UUID REFERENCES registries (registry_id),
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);