package models

import "database/sql"

type Password struct {
	ID    int
	Title string
	Value string
}

type PasswordModel struct {
	DB *sql.DB
}

func (pm *PasswordModel) Add(title, value string) error {
	query := `
	INSERT INTO passwords (Title, Value)
		VALUES (?, ?)
	;`

	_, err := pm.DB.Exec(query, title, value)

	return err
}

func (pm *PasswordModel) Get(title string) (string, error) {
	query := `
	SELECT value FROM passwords
		WHERE Title = ?
	;`

	row := pm.DB.QueryRow(query, title)

	value := ""
	err := row.Scan(&value)
	if err != nil {
		return "", nil
	}

	return value, nil
}

func (pm *PasswordModel) List() ([]*Password, error) {
	query := `
	SELECT * FROM passwords
	;`

	rows, err := pm.DB.Query(query)
	if err != nil {
		return nil, err
	}

	passwords := make([]*Password, 0)
	for rows.Next() {
		password := Password{}

		err = rows.Scan(&password.ID, &password.Title, &password.Value)
		if err != nil {
			return nil, err
		}

		passwords = append(passwords, &password)
	}

	return passwords, nil
}
