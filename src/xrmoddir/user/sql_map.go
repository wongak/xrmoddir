package user

import (
	"database/sql"
	"fmt"
	"time"
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
SELECT
	id
FROM %s
WHERE
	username = ?
	`
	stmts["selectIdByUsername"] = fmt.Sprintf(stmt, TABLE_NAME_USERS)

	stmt = `
SELECT
	u.id
FROM %s AS u
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
	stmts["selectIdByMeta"] = fmt.Sprintf(stmt, TABLE_NAME_USERS, TABLE_NAME_USER_METADATA, TABLE_NAME_USER_METADATA)

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

	stmt = `
SELECT
	u.id
FROM %s AS u
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
	stmts["selectMetadata"] = fmt.Sprintf(stmt, TABLE_NAME_USERS, TABLE_NAME_USER_METADATA, TABLE_NAME_USER_METADATA)
}

func selectUser(row *sql.Row) (*User, error) {
	u := new(User)
	err := row.Scan(&u.Id, &u.Username, &u.Created, &u.password, &u.Email, &u.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("SQL error on select user: %v", err)
	}
	return u, nil
}

func SQLUserById(db *sql.DB, id int64) (*User, error) {
	row := db.QueryRow(stmts["select"]+" WHERE u.id = ?", id)
	return selectUser(row)
}

func SQLUserByUsername(db *sql.DB, username string) (*User, error) {
	row := db.QueryRow(stmts["select"]+" WHERE u.username = ?", username)
	return selectUser(row)
}

func SQLIdByUsername(db *sql.DB, username string) (int64, error) {
	row := db.QueryRow(stmts["selectIdByUsername"], username)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("SQLIdByUsername: SQL Error: %v", err)
	}
	return id, nil
}

func SQLIdByEmail(db *sql.DB, email string) (int64, error) {
	row := db.QueryRow(stmts["selectIdByMeta"]+" WHERE meta.email = ?", email)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("SQLIdByEmail: SQL Error: %v", err)
	}
	return id, nil
}

func SQLInsertUser(tx *sql.Tx, u *User) (int64, error) {
	stmt, err := tx.Prepare(stmts["insertUser"])
	if err != nil {
		return 0, fmt.Errorf("SQLInsertUser: SQL Error: %v", err)
	}
	u.Created = time.Now().UTC()
	res, err := stmt.Exec(u.Username, u.Created)
	if err != nil {
		return 0, fmt.Errorf("SQLInsertUser: SQL Error on exec: %v", err)
	}
	u.Id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("SQLInsertUser: Last insert id error: %v", err)
	}
	return u.Id, nil
}

func SQLInsertPassword(tx *sql.Tx, u *User) error {
	stmt, err := tx.Prepare(stmts["insertPassword"])
	if err != nil {
		return fmt.Errorf("SQLInsertPassword: SQL Error: %v", err)
	}
	ts := time.Now().UTC().UnixNano()
	_, err = stmt.Exec(u.Id, ts, u.password)
	if err != nil {
		return fmt.Errorf("SQLInsertPassword: SQL Error on exec: %v", err)
	}
	return nil
}

func SQLInsertMeta(tx *sql.Tx, u *User) error {
	stmt, err := tx.Prepare(stmts["insertMetadata"])
	if err != nil {
		return fmt.Errorf("SQLInsertMeta: SQL Error: %v", err)
	}
	ts := time.Now().UTC().UnixNano()
	_, err = stmt.Exec(u.Id, ts, u.Email, u.Active)
	if err != nil {
		return fmt.Errorf("SQLInsertMeta: SQL Error on exec: %v", err)
	}
	return nil
}
