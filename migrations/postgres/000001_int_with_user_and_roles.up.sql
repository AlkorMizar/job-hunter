-- Table: user
CREATE TABLE IF NOT EXISTS "user"
(
    user_id SERIAL  NOT NULL,
    login varchar(40) COLLATE pg_catalog."default" NOT NULL,
    email varchar(255) COLLATE pg_catalog."default" NOT NULL ,
    date_created time with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password varchar(60) COLLATE pg_catalog."default" NOT NULL,
    last_check timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    full_name varchar(200) COLLATE pg_catalog."default" NOT NULL,
	UNIQUE(login,email),
    CONSTRAINT pk_user_id PRIMARY KEY (user_id)
);

-- Table: role
CREATE TABLE IF NOT EXISTS "role" (
  role_id SERIAL NOT NULL,
  name varchar(15) NOT NULL,
  UNIQUE (name),
  CONSTRAINT pk_role_id PRIMARY KEY (role_id)
);

-- Table: user_has_role
CREATE TABLE IF NOT EXISTS "user_has_role"
(
    fk_user_id integer NOT NULL,
    fk_role_id integer NOT NULL,
    CONSTRAINT pk_user_role PRIMARY KEY (fk_user_id, fk_role_id),
    CONSTRAINT fk_role FOREIGN KEY (fk_role_id)
        REFERENCES "role" (role_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_user FOREIGN KEY (fk_user_id)
        REFERENCES "user" (user_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
