-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, name, type)
VALUES (1, 'Alice', 'borrower'),
       (2, 'Bob', 'investor'),
       (3, 'Franda', 'investor'),
       (4, 'Charlie', 'employee');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id IN (1, 2, 3, 4);
-- +goose StatementEnd
