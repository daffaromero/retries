CREATE TABLE payment (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36), 
  customer_id VARCHAR(36),
  amount DECIMAL(10, 2),
  status VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);