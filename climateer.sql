CREATE TABLE Users (
    id INT PRIMARY KEY AUTO_INCREMENT,                -- Primary key, auto-incremented
    first_name VARCHAR(50) NOT NULL,                  -- First name of the user (removed UNIQUE constraint)
    last_name VARCHAR(50) NOT NULL,                   -- Last name of the user (removed UNIQUE constraint)
    email VARCHAR(255) UNIQUE NOT NULL,               -- Email must be unique
    password_hash VARCHAR(255) NOT NULL,              -- Storing the hashed password
    phone VARCHAR(255) NOT NULL,                      -- Phone number of the user
    edu_institute VARCHAR(255) NOT NULL,              -- Educational institute (removed UNIQUE constraint)
    session_key VARCHAR(50) UNIQUE NOT NULL,          -- Unique session key for tracking user sessions
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,   -- Timestamp when the record was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for last update
);


-- Create UserMeasurements table
CREATE TABLE UserMeasurements (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    country_id INT,
    indicator_id INT,
    year INT,
    value FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (country_id) REFERENCES Countries(id),
    FOREIGN KEY (indicator_id) REFERENCES Indicators(id)
);

-- Create Countries table
CREATE TABLE Countries (
    id INT PRIMARY KEY,
    country_code VARCHAR(3) UNIQUE,
    country_name VARCHAR(255)
);

-- Create Indicators table
CREATE TABLE Indicators (
    id INT PRIMARY KEY,
    Indicator VARCHAR(255) UNIQUE
);

-- Create Measurements table
CREATE TABLE Measurements (
    id INT PRIMARY KEY,
    country_id INT,
    indicator_id INT,
    year INT,
    value FLOAT,
    FOREIGN KEY (country_id) REFERENCES Countries(id),
    FOREIGN KEY (indicator_id) REFERENCES Indicators(id)
);