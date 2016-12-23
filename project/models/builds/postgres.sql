

--管理员表
CREATE SEQUENCE "public"."seq_admin_id" INCREMENT 1 START 10000;
CREATE TABLE "public"."admin" (
	"id" int8 NOT NULL DEFAULT nextval('seq_admin_id'),
	"account" varchar(50) NOT NULL,
	"password" varchar(50) NOT NULL,
	"name" varchar(50) NOT NULL,
	"role" varchar(50) NOT NULL DEFAULT 'nobody',
	"changed" timestamptz NOT NULL DEFAULT now(),
	"created" timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE,
	CONSTRAINT "account" UNIQUE ("account") NOT DEFERRABLE INITIALLY IMMEDIATE
)
WITH (OIDS=FALSE);


--添加管理员
INSERT INTO "public"."admin"("account","password","name","role") VALUES ('admin','admin','admin','system');

