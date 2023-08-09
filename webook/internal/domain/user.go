package domain

// User 领域对象，是 DDD 中的 entity
// BO(business object)
type User struct {
	Id           int64
	Email        string
	Password     string
	Ctime        int64
	Nickname     string
	Birthday     string
	Introduction string
}

//type Address struct {
//}
