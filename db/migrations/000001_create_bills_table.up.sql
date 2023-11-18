CREATE TABLE bills (
    congress VARCHAR(255),
    bill_number INT,
    origin_chamber VARCHAR(10),
    origin_chamber_code VARCHAR(10),
    title VARCHAR(8000),
    type VARCHAR(650),
    url VARCHAR(3000),
    latest_action_date TIMESTAMP,
    latest_action_text VARCHAR(1024),
    update_date TIMESTAMP,
    update_including_text DATE,
    PRIMARY KEY (type, bill_number),
    created TIMESTAMP
);
