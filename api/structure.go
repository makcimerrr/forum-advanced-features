package api

type Discussion struct {
	ID            int
	Title         string
	Message       string
	Username      string
	Category      []string
	Liked         bool
	Disliked      bool
	NumberLike    int
	NumberDislike int
	NumberComment int
}

// Ajoutez cette structure pour représenter un message
type Comment struct {
	ID            int
	Username      string
	Message       string
	Discussion_id int
	Liked         bool // Champ pour indiquer si l'utilisateur a aimé cette discussion
	Disliked      bool
	NumberLike    int
	NumberDislike int
}

type Categories struct {
	ID       int
	Category string
}

type Notification struct {
	ID int
	User_id int
	Discussion_id int
	Message string
	vu bool

}
