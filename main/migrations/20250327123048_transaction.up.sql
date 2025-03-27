-- Migration UP: Create transactions table
CREATE TABLE transactions (
                              id INT AUTO_INCREMENT PRIMARY KEY,
                              source_id INT,
                              external_source_id BIGINT,
                              destination_id INT,
                              external_destination_id BIGINT,
                              type VARCHAR(10) NOT NULL,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                              is_deleted BOOLEAN DEFAULT FALSE,
                              deleted_at TIMESTAMP NULL,
                              amount INT NOT NULL,
                              FOREIGN KEY (source_id) REFERENCES wallets(id) ON DELETE SET NULL,
                              FOREIGN KEY (destination_id) REFERENCES wallets(id) ON DELETE SET NULL
);

-- Create indexes for performance
CREATE INDEX idx_source_type_created ON transactions (source_id, type, created_at);
CREATE INDEX idx_destination_type_created ON transactions (destination_id, type, created_at);