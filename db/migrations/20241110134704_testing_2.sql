-- migrate:up

CREATE TABLE employees_2 (
    id INT PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(100),
    phone_number VARCHAR(15),
    hire_date DATE,
    salary DECIMAL(10, 2)
);
-- migrate:down

