package user

import (
	"fmt"
)

const (
	TABLE_NAME_USERS          = "users"
	TABLE_NAME_USER_PASSWORDS = "user_passwords"
	TABLE_NAME_USER_METADATA  = "user_metadata"
)

var stmts map[string]string

func init() {
	stmts = make(map[string]string)
	stmt := `
SELECT
	u.id,
	u.username,
	u.created,
	pw.password,
	meta.email,
	meta.active
FROM %s AS u
INNER JOIN %s AS pw ON
	pw.user_id = u.id
	AND
	pw.timestamp = (
		SELECT MAX(timestamp)
		FROM %s
		WHERE
		user_id = pw.user_id
	)
INNER JOIN %s AS meta ON
	meta.user_id = u.id
	AND
	meta.timestamp = (
		SELECT MAX(timestamp)
		FROM %s
		WHERE
		user_id = meta.user_id
	)
	`
	stmts["select"] = fmt.Sprintf(stmt, TABLE_NAME_USERS, TABLE_NAME_USER_PASSWORDS, TABLE_NAME_USER_PASSWORDS, TABLE_NAME_USER_METADATA, TABLE_NAME_USER_METADATA)

	stmt = `
INSERT INTO %s
(username, created)
VALUES
(?, ?)
	`
	stmts["insertUser"] = fmt.Sprintf(stmt, TABLE_NAME_USERS)

	stmt = `
INSERT INTO %s
(user_id, timestamp, password)
VALUES
(?, ?, ?)
	`
	stmts["insertPassword"] = fmt.Sprintf(stmt, TABLE_NAME_USER_PASSWORDS)

	stmt = `
INSERT INTO %s
(user_id, timestamp, email, active)
VALUES
(?, ?, ?, ?)
	`
	stmts["insertMetadata"] = fmt.Sprintf(stmt, TABLE_NAME_USER_METADATA)
}
