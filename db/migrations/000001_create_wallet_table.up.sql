CREATE TABLE "wallet" (
	"id" serial NOT NULL,
	"name" TEXT NOT NULL,
	"balance" numeric(8,4),
	UNIQUE(name),
	CONSTRAINT "wallet_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);

ALTER TABLE "wallet" ADD CONSTRAINT "balance_nonnegative" check ("balance" >= 0);