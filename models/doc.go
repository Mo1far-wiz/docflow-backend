package models

import (
	"docflow-backend/db"
	"time"
)

type Doc struct {
	ID          int64
	DocName     string `binding:"required"`
	DateTime    time.Time
	Faculty     string `binding:"required"`
	Specialty   string `binding:"required"`
	YearOfStudy int64  `binding:"required"`
	UserID      int64
}

func (d *Doc) Save() error {
	query := "INSERT INTO docs (docName, dateTime, faculty, specialty, yearOfStudy, user_id) VALUES (?, ?, ?, ?, ?, ?)"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(d.DocName, d.DateTime, d.Faculty, d.Specialty, d.YearOfStudy, d.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	d.ID = id

	return err
}

func GetAllDocsForUser(userId int64) ([]Doc, error) {
	query := "SELECT * FROM docs WHERE user_id = ?"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var docs []Doc
	for rows.Next() {
		var doc Doc
		err := rows.Scan(&doc.ID, &doc.DocName, &doc.DateTime, &doc.Faculty, &doc.Specialty, &doc.YearOfStudy, &doc.UserID)
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

func GetDocByID(id int64) (*Doc, error) {
	query := "SELECT * FROM docs WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var doc Doc
	err := row.Scan(&doc.ID, &doc.DocName, &doc.DateTime, &doc.Faculty, &doc.Specialty, &doc.YearOfStudy, &doc.UserID)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
