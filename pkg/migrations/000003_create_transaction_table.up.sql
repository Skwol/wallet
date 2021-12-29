CREATE TYPE transaction_type AS ENUM ('deposit', 'withdraw', 'transfer');

CREATE TABLE "transaction" (
	"id" serial NOT NULL,
	"sender_id" bigint NOT NULL,
	"receiver_id" bigint NOT NULL,
	"amount" numeric(8,4) NOT NULL,
	"date" timestamp NOT NULL,
	"tran_type" transaction_type NOT NULL,
	CONSTRAINT "transaction_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);

ALTER TABLE "transaction" ADD CONSTRAINT "transaction_fk_sender" FOREIGN KEY ("sender_id") REFERENCES "wallet"("id");
ALTER TABLE "transaction" ADD CONSTRAINT "transaction_fk_receiver" FOREIGN KEY ("receiver_id") REFERENCES "wallet"("id");
ALTER TABLE "transaction" ADD CONSTRAINT "amount_morethenzero" CHECK ("amount" > 0);