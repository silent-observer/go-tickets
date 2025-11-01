CREATE TABLE projects (
    id serial PRIMARY KEY,
    name varchar(20) UNIQUE NOT NULL
);
INSERT INTO projects(name) VALUES ('default');

CREATE TABLE tickets (
    id serial PRIMARY KEY,
    project_id integer NOT NULL,
    ticket_num integer NOT NULL,
    title text NOT NULL,
    content text NOT NULL,
    status_id integer NOT NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    edited_date timestamptz NOT NULL DEFAULT now(),
    created_by integer NOT NULL
);

CREATE TABLE users (
    id serial PRIMARY KEY,
    username varchar(120) UNIQUE NOT NULL
);

CREATE TABLE statuses (
    id serial PRIMARY KEY,
    project_id integer NOT NULL,
    name varchar(50) NOT NULL,
    icon varchar(20) NOT NULL
);

CREATE TABLE ticket_links (
    ticket_from integer NOT NULL,
    ticket_to integer NOT NULL,
    type varchar(20) NOT NULL,
    PRIMARY KEY (ticket_from, ticket_to)
);

CREATE TABLE tags (
    id serial PRIMARY KEY,
    project_id integer NOT NULL,
    name varchar(50) NOT NULL
);

CREATE TABLE ticket_tags (
    ticket_id integer NOT NULL,
    tag_id integer NOT NULL,
    val varchar(50) NOT NULL,
    PRIMARY KEY (ticket_id, tag_id)
);

CREATE TABLE ticket_assigned (
    ticket_id integer NOT NULL,
    user_id integer NOT NULL,
    PRIMARY KEY (ticket_id, user_id)
);

CREATE TABLE boards (
    id serial PRIMARY KEY,
    project_id integer NOT NULL,
    name varchar(50) NOT NULL
);

CREATE TABLE board_tickets (
    board_id integer NOT NULL,
    ticket_id integer NOT NULL,
    PRIMARY KEY (board_id, ticket_id)
);