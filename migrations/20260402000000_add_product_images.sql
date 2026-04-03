-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Drop old image_url column from products (optional, can keep for backward compatibility)
-- ALTER TABLE products DROP COLUMN IF EXISTS image_url;

CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id);