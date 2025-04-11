package database

type Database[T any] struct {
	Map map[int]T
}

type database[T any] interface {
	GetUser(id int)
	GetAllUsers()
	DeleteUser(id int)
	GetLen()
	AddUser(User T)
}

// CreateDatabase creates and returns a Database that contains type T
func CreateDatabase[T any]() Database[T] {
	return Database[T]{make(map[int]T, 0)}
}

// GetUser returns User corresponding to id and ok bool value, if ok true, the the User was in the database, if ok false then the User wasn't in the database
func (db Database[T]) GetUser(id int) (T, bool) {
	User, ok := db.Map[id]
	return User, ok
}

// GetAllUsers returns a slices containing all elements of the database, if the Database is empty then it return an empty slice
func (db Database[T]) GetAllUsers() (Users []T) {
	Users = make([]T, 0)
	for i := range db.Map {
		Users = append(Users, db.Map[i])
	}
	return Users
}

// DeleteUser deletes a User from the Database corresponding to the id
func (db Database[T]) DeleteUser(id int) {
	delete(db.Map, id)
}

// GetLen return the amount of elements in the Database
func (db Database[T]) GetLen() int {
	return len(db.Map)
}

// AddUser adds an elements to the database
func (db Database[T]) AddUser(user T) {
	db.Map[db.GetLen()] = user
}
