CREATE TABLE "account" (
	"id" serial NOT NULL,
	"username" TEXT NOT NULL UNIQUE,
	CONSTRAINT "account_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);