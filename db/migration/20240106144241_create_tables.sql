-- +migrate Up notransaction
CREATE TYPE "gender" AS ENUM (
    'MALE',
    'FEMALE'
);

CREATE TABLE IF NOT EXISTS "provinces" (
    "id" BIGINT PRIMARY KEY,
    "name" VARCHAR(255),
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "cities" (
    "id" BIGINT PRIMARY KEY,
    "province_id" BIGINT,
    "name" VARCHAR(255),
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "candidates" (
    "id" BIGINT PRIMARY KEY,
    "full_name" VARCHAR(255),
    "email" VARCHAR(255),
    "phone" VARCHAR(50),
    "password" TEXT,
    "date_of_birth" DATE,
    "latitude" TEXT NOT NULL,
    "longitude" TEXT NOT NULL,
    "gender" gender,
    "city_id" BIGINT,
    "province_id" BIGINT,
    "last_education" TIMESTAMP,
    "last_experience" TIMESTAMP,
    "login_date" TIMESTAMP,
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "educations" (
    "id" BIGINT PRIMARY KEY,
    "candidate_id" BIGINT,
    "institution_name" VARCHAR(255),
    "major" VARCHAR(255),
    "start_year" DATE,
    "end_year" DATE,
    "until_now" BOOLEAN,
    "gpa" DOUBLE PRECISION,
    "flag" VARCHAR(50),
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "experiences" (
    "id" BIGINT PRIMARY KEY,
    "candidate_id" BIGINT,
    "company_name" VARCHAR(255),
    "company_address" TEXT,
    "start_year" DATE,
    "end_year" DATE,
    "until_now" BOOLEAN,
    "position" VARCHAR(255),
    "job_desc" TEXT,
    "flag" VARCHAR(50),
    "created_at" TIMESTAMP NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
    "deleted_at" TIMESTAMP
);

ALTER TABLE "educations" ADD FOREIGN KEY ("candidate_id") REFERENCES "candidates" ("id");

ALTER TABLE "experiences" ADD FOREIGN KEY ("candidate_id") REFERENCES "candidates" ("id");

ALTER TABLE "cities" ADD FOREIGN KEY ("province_id") REFERENCES "provinces" ("id");

-- +migrate Down