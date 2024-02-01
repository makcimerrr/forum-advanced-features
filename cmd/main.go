package main

import (
	"fmt"
	"forum/api"
	forum "forum/internal"
	"forum/pkg/facebook"
	"forum/pkg/github"
	"net/http"

)

func main() {

	forum.Jsp()

	api.CreateBDD()

	// Login route
	http.HandleFunc("/login/github/", github.GithubLoginHandler)

	// Github callback
	http.HandleFunc("/login/github/callback", github.GithubCallbackHandler)

	// Route where the authenticated user is redirected to
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		github.LoggedinHandler(w, r, "")
	})

	http.HandleFunc("/login/facebook", facebook.HandleFacebookLogin)
	http.HandleFunc("/oauth2callback", facebook.HandleFacebookCallback)

	

	http.HandleFunc("/", forum.HomeHandler)
	http.HandleFunc("/404", forum.HandleNotFound)
	http.HandleFunc("/500", forum.HandleServerError)
	http.HandleFunc("/400", forum.HandleBadRequest)
	http.HandleFunc("/logOrSign", forum.LogOrSignHandler)
	http.HandleFunc("/log_in", forum.LogInHandler)
	http.HandleFunc("/sign_up", forum.SignUpHandler)
	http.HandleFunc("/create_discussion", forum.CreateDiscussionHandler)
	http.HandleFunc("/edit_discussion", forum.EditDiscussionHandler)
	http.HandleFunc("/edit_comment", forum.EditCommentHandler)
	http.HandleFunc("/discussion/", forum.ShowDiscussionHandler)
	http.HandleFunc("/profil", forum.ProfilHandler)
	

	//fake page
	http.HandleFunc("/add_message/", forum.AddMessage)
	http.HandleFunc("/logout", forum.Logout)
	http.HandleFunc("/like/", forum.WhereIsTheLike)
	http.HandleFunc("/dislike/", forum.WhereIsTheDislike)
	http.HandleFunc("/editPost", forum.EditPost)
	http.HandleFunc("/editComment", forum.EditComment)
	http.HandleFunc("/deletePost", forum.DeletePost)
	http.HandleFunc("/deleteComment", forum.DeleteComment)
	

	// Définir le dossier "static" comme dossier de fichiers statiques
	fs := http.FileServer(http.Dir("./web/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))


	fmt.Println("Voici le lien pour ouvrir la page web http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
