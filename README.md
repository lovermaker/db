# DAO (Data Access Object) Client for gorm library
## Quick Start
### install library
```shell
go get github.com/lovemaker/db
```
### write code

if you need to access database using gorm library, you can do it with following steps:

1. write your DAO object
```go
// DAO data access object
type DAO struct {
	user  *UserDAO
	class *ClassDAO
}

// NewDAO new a dao
func NewDAO(db *gorm.DB) *DAO {
	return &DAO{
		user:  &UserDAO{db: db},
		class: &ClassDAO{db: db},
	}
}

// User returns user dao
func (d *DAO) User() *UserDAO {
	return d.user
}

// Class returns class dao
func (d *DAO) Class() *ClassDAO {
	return d.class
}

// UserDAO user data access object
type UserDAO struct {
	db *gorm.DB
}

// CreateUser create a user
func (d *UserDAO) CreateUser() error {
	return nil
}

// ClassDAO class data access object
type ClassDAO struct {
	db *gorm.DB
}

// CreateClass create a class
func (d *ClassDAO) CreateClass() error {
	return nil
}
```
2. use this library to access db
```go

var config db.Config

func init() {
	config.Host = "127.0.0.1"
	config.Port = 3306
	config.UserName = "root"
	config.Password = "12345678"
}

func main() {
    // new db client
    client, err := db.NewClient(NewDAO, &config)
    if err != nil {
		return
    }
    // create a user
    err = client.DAO().User().CreateUser()
    if err != nil {
		return
    }
    // create a class
    err = client.DAO().Class().CreateClass()
    if err != nil {
        return
    }
    // start a transaction
    tx := client.StartTransaction()
    err = tx.DAO().User().CreateUser()
    if err != nil {
        tx.Rollback()
        return
    }
    err = tx.DAO().Class().CreateClass()
    if err != nil {
        tx.Rollback()
        return
    }
    tx.Commit()
}
```
