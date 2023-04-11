# repoimpl
Repositories, migrations and tests with model


We write the necessary model.
The function writes table queries, repositories and tests for each model in the migrations folder.

## Install 
```
go get github.com/realtemirov/repoimpl
```

## Repository
```
type User struct {
    Username string
    Password string
}

err := repoimpl.NewRepository(User{})
if err != nil {
    panic(err)
}
```
Repository implement CRUD to *postgres* database with sql package and creates **interface storage**

```
go run main.go
```
./repository
    ./postgres/
        ./user.go - CRUD methods
    ./storage.go - CRUD interfaces

## Migrations
```
type User struct {
    Username string
    Password string
}

err := repoimpl.NewMigration(User{})
if err != nil {
    panic(err)
}
```
Creates **migration_user.sql** in the **migrations** folder

```
go run main.go
```

```
CREATE TABLE IF NOT EXITS "users" (
    "username" TEXT,
    "password" TEXT
);
```

If **db** is written **tag** it will be written by tag

```
type User struct {
    Username string `db:"user_name"`
    Password string `db:"pass_word"`
}

err := repoImpl.NewDBTable(User{})
if err != nil {
    panic(err)
}
```
**migration_user.sql**
```
CREATE TABLE IF NOT EXITS "users" (
    "user_name" TEXT,
    "pass_word" TEXT
);
```
