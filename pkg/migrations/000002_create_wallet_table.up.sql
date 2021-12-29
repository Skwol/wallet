CREATE TABLE "wallet" (
	"id" serial NOT NULL,
	"name" TEXT NOT NULL,
	"account_id" bigint NOT NULL,
	"balance" numeric(8,4),
	UNIQUE(account_id, name),
	CONSTRAINT "wallet_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);

ALTER TABLE "wallet" ADD CONSTRAINT "wallet_fk0" FOREIGN KEY ("account_id") REFERENCES "account"("id");
ALTER TABLE "wallet" ADD CONSTRAINT "balance_nonnegative" check ("balance" >= 0);