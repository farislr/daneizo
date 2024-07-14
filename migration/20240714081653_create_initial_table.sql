-- +goose Up
CREATE TABLE customers
(
    id         BIGINT PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    type       VARCHAR(100) NOT NULL COMMENT "borrower, investor, employee",
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE loans
(
    id                   BIGINT PRIMARY KEY,
    borrower_id          BIGINT         NOT NULL,
    principal_amount     DECIMAL(10, 2) NOT NULL,
    invested_amount      DECIMAL(10, 2) NOT NULL,
    interest_rate        DECIMAL(5, 2)  NOT NULL COMMENT "Per annum",,
    status               VARCHAR(100)   NOT NULL COMMENT "proposed, approved, invested, disbursed",
    approval_date        TIMESTAMP,
    approval_employee_id BIGINT,
    disbursement_date    TIMESTAMP,
    created_at           TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE loan_investments
(
    id          BIGINT PRIMARY KEY,
    loan_id     BIGINT         NOT NULL,
    investor_id BIGINT         NOT NULL,
    amount      DECIMAL(10, 2) NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS loan_investments;
