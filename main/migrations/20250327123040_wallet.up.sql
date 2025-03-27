-- Migration UP: Create wallets table
CREATE TABLE wallets (
                         id INT AUTO_INCREMENT PRIMARY KEY,
                         user_id INT NOT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         is_deleted BOOLEAN DEFAULT FALSE,
                         deleted_at TIMESTAMP NULL,
                         amount INT DEFAULT 0,
                         currency VARCHAR(3) DEFAULT 'INR'
);