package main

import (
	"fmt"
	"net/http"

	"forum/forum"
)

func main() {
	http.HandleFunc("/", forum.Home)
	http.HandleFunc("/home", forum.Home)
	http.HandleFunc("/404", forum.HandleNotFound)
	http.HandleFunc("/500", forum.HandleServerError)
	http.HandleFunc("/400", forum.HandleBadRequest)
	http.HandleFunc("/logorsign", forum.Logorsign)
	http.HandleFunc("/log_in", forum.Log_in)
	http.HandleFunc("/sign_up", forum.Sign_up)
	http.HandleFunc("/logout", forum.Logout)
	http.HandleFunc("/create_discussion", forum.CreateDiscussion)
	http.HandleFunc("/discussion/", forum.ShowDiscussion)
	http.HandleFunc("/add_message/", forum.AddMessage)
	http.HandleFunc("/like/", forum.LikeDiscussion)
	http.HandleFunc("/dislike/", forum.DislikeDiscussion)

	

	// Définir le dossier "static" comme dossier de fichiers statiques
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Voici le lien pour ouvrir la page web http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
