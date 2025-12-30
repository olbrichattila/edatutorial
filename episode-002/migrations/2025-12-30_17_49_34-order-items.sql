CREATE TABLE order_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    product_id CHAR(36) NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES order_heads(id)
);