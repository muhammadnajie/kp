package users

import (
	"database/sql"
	database "github.com/muhammadnajie/kp/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (u *User) Authenticate() bool {
	stmt, err := database.Db.Prepare("select password from users where username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(u.Username)
	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckPasswordHash(u.Password, hashedPassword)
}

func (u *User) Create() {
	stmt, err := database.Db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	print(stmt)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(u.Password)
	_, err = stmt.Exec(u.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select id from users WHERE username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
