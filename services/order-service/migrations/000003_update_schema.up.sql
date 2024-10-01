ALTER TABLE orders
ALTER COLUMN id SET TYPE VARCHAR;
ALTER COLUMN product_id SET NAME product_ids;
ALTER COLUMN product_ids SET TYPE TEXT[];
ALTER COLUMN quantity SET NAME products_details;
ALTER COLUMN products_details SET TYPE JSONB;
ADD COLUMN settlement_status VARCHAR;
ADD COLUMN total_payment INT;
ADD COLUMN admin_fee INT;
ADD COLUMN grand_total INT;
ADD COLUMN payment_link VARCHAR;

/*
updated schema:
  string id = 1;
  string customer_id = 2;
  repeated string product_ids = 3;
  repeated ProductDetails products_details = 4;
  string settlement_status = 5;
  int32 total_payment = 6;
  int32 admin_fee = 7;
  int32 grand_total = 8;
  string payment_link = 9;

starting schema:
CREATE TABLE "orders" (
  "id"    INT PRIMARY KEY,
  "customer_id" INT,
  "product_id"  INT,
  "quantity"    INT
)

*/
