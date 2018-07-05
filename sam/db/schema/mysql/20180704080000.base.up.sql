-- all known organisations (crust instances) and our relation towards them
CREATE TABLE organisations (
  id               BIGINT UNSIGNED NOT NULL,
  fqn              TEXT            NOT NULL, -- fully qualified name of the organisation
  label            TEXT            NOT NULL, -- display name of the organisation

  PRIMARY KEY (id)
);

-- Keeps all known teams
CREATE TABLE teams (
  id               BIGINT UNSIGNED NOT NULL,
  label            TEXT            NOT NULL, -- display name of the team
  handle           TEXT            NOT NULL, -- team handle string

  PRIMARY KEY (id)
);

-- Keeps all known channels
CREATE TABLE channels (
  id               BIGINT UNSIGNED NOT NULL,
  label            TEXT            NOT NULL, -- display name of the channel
  meta             JSON            NOT NULL,

  PRIMARY KEY (id)
);

-- Keeps all known users, home and external organisation
--   changes are stored in audit log
CREATE TABLE users (
  id               BIGINT UNSIGNED NOT NULL,
  username         TEXT            NOT NULL,
  password         TEXT,
  meta             JSON            NOT NULL,

  rel_organisation BIGINT UNSIGNED NOT NULL REFERENCES organisation(id),

  PRIMARY KEY (id)
);

-- Keeps team memberships
CREATE TABLE team_members (
  rel_team         BIGINT UNSIGNED NOT NULL REFERENCES organisation(id),
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),

  PRIMARY KEY (rel_team, rel_user)
);

-- handles membership, visibility and level of access
CREATE TABLE channel_members (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),

  -- display only messages from-to to particular user/organisation
  messages_since   DATETIME            NULL,
  messages_until   DATETIME            NULL,

  PRIMARY KEY (rel_channel, rel_user)
);

CREATE TABLE messages (
  id               BIGINT UNSIGNED NOT NULL,

  -- basic set:
  --    NULL                            for common text messages
  --    image/*; disposition=inline     for inline images (does not follow stdandard http headers, but
  --    crust/system                    body holds one of defined system messages (user joining, parting...)
  --    crust/reaction                  body holds reaction to a message
  --    crust/preview                   body holds preview of the message it relates to in form of metadata (no binary blobs!)
  --    crust/url                       body holds only URL
  --    crust/pin                       message it relates to is pinned, body holds pin details
  --    crust/flag                      message it relates to is flaged, body holds flag information
  --    */*                             body holds reference to the uploaded file and its metadata
  --
  mimetype         VARCHAR(255)         NULL,

  -- null body only valid when rel_message is set => message removal
  body             TEXT                 NULL,

  -- the contributor
  rel_user         BIGINT  UNSIGNED NOT NULL REFERENCES users(id),

  -- message's channel
  rel_channel      BIGINT  UNSIGNED NOT NULL REFERENCES channels(id),

  -- replies, edits, reactions, flags, attachments...
  rel_message      BIGINT  UNSIGNED NOT NULL REFERENCES messages(id),

  PRIMARY KEY (id)
);


CREATE TABLE channel_views (
  rel_channel      BIGINT UNSIGNED NOT NULL REFERENCES channels(id),
  rel_user         BIGINT UNSIGNED NOT NULL REFERENCES users(id),

  -- timestamp of last view, should be enough to find out which messaghr
  viewed_at        DATETIME        NOT NULL DEFAULT NOW(),

  -- new messages count since last view
  new_since        INT    UNSIGNED NOT NULL DEFAULT 0,

  PRIMARY KEY (rel_user, rel_channel)
);