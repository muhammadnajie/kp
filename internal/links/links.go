package links

import (
	"fmt"
	database "github.com/muhammadnajie/kp/internal/pkg/db/mysql"
	"github.com/muhammadnajie/kp/internal/users"
	"log"
)

type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

func (l Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO links(title, address, userid) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(l.Title, l.Address, l.User.ID)
	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Row inserted!")
	return id
}

func (l *Link) Update() (int64, error) {
	stmt, err := database.Db.Prepare("UPDATE links SET title = ?, address = ? where id = ?")
	if err != nil {
		log.Println("error sini kah?")
		log.Fatal(err)
	}
	res, err := stmt.Exec(l.Title, l.Address, l.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("updated...")
	return res.RowsAffected()
}

func (l *Link) Delete() (int64, error) {
	stmt, err := database.Db.Prepare("DELETE FROM links where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(l.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("deleted")
	return res.RowsAffected()
}

func GetByTitle(title string, userID string) ([]Link, error) {
	query := `SELECT l.id, l.title, l.address, l.userid, u.username
		from links l 
		inner join users u
		on u.id = l.userid
	where l.userid = ?`
	query = fmt.Sprintf(`%s AND title like '%%%s%%'`, query, title)
	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	var links []Link
	for rows.Next() {
		var link Link
		var id string
		var username string
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username)
		if err != nil {
			return nil, err
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

func GetByID(id int, userID int) (Link, error) {
	var link Link
	var username string
	query := `SELECT l.id, l.title, l.address, u.username
		FROM links l 
			INNER JOIN users u
			ON l.userid = ?
		where l.id = ?`
	err := database.Db.QueryRow(query, userID, id).
		Scan(&link.ID, &link.Title, &link.Address, &username)
	link.User = &users.User{
		ID:       string(rune(userID)),
		Username: username,
	}
	if err != nil {
		fmt.Println(err)
		return Link{}, err
	}
	return link, nil
}

func GetAll(userID string) ([]Link, error) {
	query := `SELECT l.id, l.title, l.address, l.userid, u.username
		from links l 
			inner join users u
			on u.id = ?
		where l.userid = ?`
	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(userID, userID)
	if err != nil {
		return nil, err
	}

	var links []Link
	for rows.Next() {
		var link Link
		var id string
		var username string
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username)
		if err != nil {
			return nil, err
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}
