package forum

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/api"
	"hash/fnv"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	_ "modernc.org/sqlite"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load("data.env"); err != nil {
		log.Fatal("No .env file found")
	}
}

func Jsp() {
	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New("our-google-client-id", "our-google-client-secret", "http://localhost:8080/auth/google/callback", "email", "profile"),
	)

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})

	if err := godotenv.Load("data.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func codeErreur(w http.ResponseWriter, r *http.Request, url string, route string, html string) {
	if url != route {
		http.Redirect(w, r, "/404", http.StatusFound)
	}
	_, err := template.ParseFiles(html)
	if err != nil {
		http.Redirect(w, r, "/500", http.StatusFound)
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func generateSessionToken() (string, error) {
	token := make([]byte, 32) // Crée un slice de bytes de 32 octets

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(token), nil
}

/* // isSessionValid vérifie si le token de session dans les cookies correspond à celui dans la base de données
func isSessionValid(w http.ResponseWriter, r *http.Request) (bool, string) {
	// Obtenez le jeton de session à partir des cookies
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		return true, "No session cookie"
	}
	sessionToken := sessionCookie.Value

	// Obtenez le nom d'utilisateur à partir des cookies
	userCookie, err := r.Cookie("username")
	if err != nil {
		return false, "No username cookie"
	}
	username := userCookie.Value

	db, err := api.OpenBDD()
			if err != nil {
				http.Error(w, "Internal Server Error Open BDD", http.StatusInternalServerError)
			}

	// Vérifiez si le nom d'utilisateur est vide
	if username == "" {
		return false, "No username in cookies"
	}

	fmt.Println("Username:", username)

	// Vérifiez si le jeton de session correspond à celui dans la base de données
	var dbSessionToken string
	dbSessionToken, err = api.GetTokenByUser(db, username)
	if err != nil {
		clearSessionCookies(w)
	}

	fmt.Println("DB Session Token:", dbSessionToken)
	fmt.Println("Session Token:", sessionToken)

	return sessionToken == dbSessionToken, "You have been disconnected"
} */

func clearSessionCookies(w http.ResponseWriter) {
	// Créer un cookie avec une date d'expiration antérieure pour effacer le cookie
	clearCookie := http.Cookie{
		Name:     "username",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &clearCookie)

	clearCookie = http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &clearCookie)
}

func CreateAndSetSessionCookies(w http.ResponseWriter, username string) error {
	// Générer un jeton de session unique
	sessionToken, err := generateSessionToken()
	if err != nil {
		return err
	}

	// Ouverture de la connexion à la base de données
	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return err
	}

	// Enregistrez le jeton de session dans la base de données
	err = api.SetTokenSession(db, username, sessionToken)
	if err != nil {
		http.Error(w, "Internal Server Error insert session user", http.StatusInternalServerError)
		return err
	}

	// Créer un cookie contenant le nom d'utilisateur
	userCookie := http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &userCookie)

	// Créer un cookie contenant le jeton de session
	sessionCookie := http.Cookie{
		Name:     "session",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &sessionCookie)

	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) {

	// Obtenez le nom d'utilisateur à partir du cookie "username"
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		// Gérez l'erreur, par exemple, en redirigeant l'utilisateur vers une page de connexion s'il n'est pas connecté.
		http.Redirect(w, r, "/log_in", http.StatusSeeOther)
		return
	}
	username := usernameCookie.Value

	//delete session_user

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	err = api.DeleteTokenSession(db, username)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	// Supprimez le cookie "username"
	userCookie := http.Cookie{
		Name:     "username",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &userCookie)

	// Supprimez le cookie "session" (s'il existe)
	sessionCookie, err := r.Cookie("session")
	if err == nil {
		sessionCookie.Value = ""
		sessionCookie.Expires = time.Unix(0, 0)
		sessionCookie.MaxAge = -1
		http.SetCookie(w, sessionCookie)
	}

	// Redirigez l'utilisateur vers une page HTML de déconnexion réussie
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func supprimerDoublons(tableau []int) []int {
	tableauSansDoublons := make([]int, 0)
	tableauMap := make(map[int]bool)

	for _, valeur := range tableau {
		if _, existe := tableauMap[valeur]; !existe {
			tableauSansDoublons = append(tableauSansDoublons, valeur)
			tableauMap[valeur] = true
		}
	}

	return tableauSansDoublons
}

// fonction permettant d'ajouter un commentaire a un post et se redirige vers ce post
func AddMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Récupérez l'ID de la discussion à partir de l'URL
		discussionID := r.URL.Path[len("/add_message/"):]
		// Convertissez l'ID de la discussion en un entier
		discussionIDInt, err := strconv.Atoi(discussionID)
		if err != nil {
			http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
			return
		}

		message := r.FormValue("message")
		// Obtenez le nom d'utilisateur à partir du cookie "username"
		usernameCookie, err := r.Cookie("username")
		if err != nil {
			// Gérer l'erreur ici, par exemple, en redirigeant l'utilisateur vers une page de connexion s'il n'est pas connecté.
			http.Redirect(w, r, "/log_in", http.StatusSeeOther)
			return
		}
		username := usernameCookie.Value

		db, err := api.OpenBDD()
		if err != nil {
			http.Error(w, "Internal Server Error open bdd", http.StatusInternalServerError)
			return
		}

		idUser, err := api.GetUserByUsername(db, username)
		if err != nil {
			http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
			return
		}

		err = api.SetComments(db, discussionIDInt, message, idUser)
		if err != nil {
			log.Printf("Erreur lors de l'insertion du message : %v", err)
			http.Error(w, "Internal Server Error set comments", http.StatusInternalServerError)
			return
		}

		// Redirigez l'utilisateur vers la page de discussion
		http.Redirect(w, r, fmt.Sprintf("/discussion/%d", discussionIDInt), http.StatusSeeOther)
		return
	}

	// Affichez la page pour écrire une discussion (write_discussion.html)
	tmpl := template.Must(template.ParseFiles("./web/templates/show_discussion.html"))
	tmpl.Execute(w, nil)
}

func EditPost(w http.ResponseWriter, r *http.Request) {

	id := r.PostFormValue("id")
	title := r.PostFormValue("title")
	message := r.PostFormValue("message")

	discussionIDInt, _ := strconv.Atoi(id)

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	err = api.EditDiscussion(db, title, message, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	// Redirigez l'utilisateur vers la page de la discussion
	http.Redirect(w, r, fmt.Sprintf("/discussion/%d", discussionIDInt), http.StatusSeeOther)

}

func EditComment(w http.ResponseWriter, r *http.Request) {

	id := r.PostFormValue("id")
	discussionId := r.PostFormValue("discussionID")
	message := r.PostFormValue("message")

	DiscussionIDInt, _ := strconv.Atoi(discussionId)
	commentIDInt, _ := strconv.Atoi(id)

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	err = api.EditComment(db, message, commentIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	// Redirigez l'utilisateur vers la page de la discussion
	http.Redirect(w, r, fmt.Sprintf("/discussion/%d", DiscussionIDInt), http.StatusSeeOther)

}

func DeletePost(w http.ResponseWriter, r *http.Request) {

	id := r.PostFormValue("id")

	discussionIDInt, _ := strconv.Atoi(id)

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	err = api.DeleteDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteLikeFromDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteDislikeFromDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteCommentFromDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteLikeCommentFromDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteDislikeCommentFromDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteDiscussionCategoryFromDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {

	id := r.PostFormValue("id")

	commentIDInt, _ := strconv.Atoi(id)

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	err = api.DeleteCommentFromID(db, commentIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteLikeCommentFromID(db, commentIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}

	err = api.DeleteDislikeCommentFromID(db, commentIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error delete token", http.StatusInternalServerError)
		return
	}


	http.Redirect(w, r, "/", http.StatusSeeOther)
}
