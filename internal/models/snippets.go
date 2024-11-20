package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct{
	ID 			int
	Title 		string
	Content 	string
	Create		time.Time
	Expires 	time.Time
}

type SnippetModel struct{
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int)(int, error){
	stmt := `INSERT INTO snippets (title, content, created, expires)
         VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err:= m.DB.Exec(stmt, title, content, expires)
	if err !=nil{
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}

func (m *SnippetModel) Get(id int)(*Snippet, error){
	stmt := `SELECT id, title, content, created, expires FROM snippets
		 WHERE expires > UTC_TIMESTAMP() AND id = ?`

		row := m.DB.QueryRow(stmt, id)
		s := &Snippet{}


		err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Create, &s.Expires)
		
		if err !=nil{
			if errors.Is(err, sql.ErrNoRows){
				return nil, ErrNoRecord
			}else{
				return nil, err
			}
			
		}
		return s, nil
}

func (m *SnippetModel) Latest()([]*Snippet, error){
	  stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	 row, err := m.DB.Query(stmt)
	 if err !=nil{
		return nil, err
	 }
	  defer row.Close()

	  snippet :=[]*Snippet{}

	  for row.Next(){
		s := &Snippet{}
		err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Create, &s.Expires)
		if err !=nil{
			return nil, err
		}
		snippet = append(snippet, s)
	  }

	  if err = row.Err(); err != nil {
		return nil, err
	  }

	  return snippet, nil


}

