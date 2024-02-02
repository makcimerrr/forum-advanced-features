package forum

import (
	"database/sql"
	"forum/api"
	"net/http"
	"strconv"
)

func WhereIsTheDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		whereDislike := r.PostFormValue("whereDislike")

		if whereDislike == "discussion" {
			DislikeDiscussion(w, r)
		} else if whereDislike == "comment" {
			DislikeComment(w, r)
		} else {
			http.Error(w, "Erreur whereDislike", http.StatusMethodNotAllowed)
		}
	}

	// Traitez d'autres méthodes HTTP comme nécessaire
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func DislikeDiscussion(w http.ResponseWriter, r *http.Request) {

	var lien []string
	var id string

	id = r.PostFormValue("id")
	lien = r.Form["lien"]

	discussionIDInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Obtenez le nom d'utilisateur à partir du cookie "username"
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		// Gérez l'erreur, par exemple, en redirigeant l'utilisateur vers une page de connexion s'il n'est pas connecté.
		http.Redirect(w, r, "/log_in", http.StatusSeeOther)
		return
	}
	username := usernameCookie.Value

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	//recupérer l'id de l'utilisateur
	idUser, err := api.GetUserByUsername(db, username)
	if err != nil {
		http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
		return
	}

	// Vérifiez si l'utilisateur a déjà disliké cette discussion

	disliked, err := api.GetDislikesFromOneDiscussion(db, discussionIDInt, idUser)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error get dislike from discussion", http.StatusInternalServerError)
		return
	}

	if disliked {
		// Si l'utilisateur a déjà disliké la discussion, supprimez le dislike
		err = api.DeleteDislikeByUserIdForDiscussion(db, discussionIDInt, idUser)
		if err != nil {
			http.Error(w, "Internal Server Error delete dislike by user id for discussion", http.StatusInternalServerError)
			return
		}
	} else {

		messageNotif := "Une personne a pas aimer votre post"

		notif, err := api.CheckIfNotificationNotDouble(db, idUser, discussionIDInt, messageNotif)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error get like", http.StatusInternalServerError)
			return
		}

		var discussion api.Discussion

		discussion, err = api.GetOneDiscussions(db, discussionIDInt)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error get like", http.StatusInternalServerError)
			return
		}

		userCreateur := discussion.Username

		//recupérer l'id de l'utilisateur
		userIDCreateur, err := api.GetUserByUsername(db, userCreateur)
		if err != nil {
			http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
			return
		}

		if notif && userIDCreateur != idUser {

		} else {

			err = api.SetNotification(db, userIDCreateur, idUser, discussionIDInt, messageNotif)
			if err != nil {
				http.Error(w, "Internal Server Error set notif", http.StatusInternalServerError)
				return
			}
		}
		// Si l'utilisateur a déjà aimé la discussion, supprimez le like
		err = api.DeleteLikeByUserIdForDiscussion(db, discussionIDInt, idUser)
		if err != nil {
			http.Error(w, "Internal Server Error delete like by user id for discussion", http.StatusInternalServerError)
			return
		}

		// Ajoutez un dislike
		err = api.SetDisLikesDiscussion(db, discussionIDInt, idUser)
		if err != nil {
			http.Error(w, "Internal Server Error dislike discussion", http.StatusInternalServerError)
			return
		}
	}

	var liens string

	switch lien[0] {
	case "Toutes les categories":
		liens = "/"
		break
	case "discussion":
		liens = "/discussion/" + id
		break
	default:

		liens = "/?categories="
		test := true
		var temp int
		var temp2 string
		for i := 0; i < len(lien); i++ {

			temp, err = api.GetIDCategoryByCategory(db, lien[i])
			if err != nil {
				http.Error(w, "Internal Server Error get id category", http.StatusInternalServerError)
				return
			}
			temp2 = strconv.Itoa(temp)
			if test {
				liens += temp2
				test = false
			} else {
				liens += "," + temp2
			}
		}
	}

	// Redirigez l'utilisateur vers la page d'accueil après la mise à jour du dislike
	http.Redirect(w, r, liens, http.StatusSeeOther)
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {

	var id string
	var discussionId string

	id = r.PostFormValue("id")
	discussionId = r.PostFormValue("discussionID")

	discussionIdInt, err := strconv.Atoi(discussionId)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}
	commentID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Obtenez le nom d'utilisateur à partir du cookie "username"
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		// Gérez l'erreur, par exemple, en redirigeant l'utilisateur vers une page de connexion s'il n'est pas connecté.
		http.Redirect(w, r, "/log_in", http.StatusSeeOther)
		return
	}
	username := usernameCookie.Value

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	//recupérer l'id de l'utilisateur
	idUser, err := api.GetUserByUsername(db, username)
	if err != nil {
		http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
		return
	}

	// Vérifiez si l'utilisateur a déjà aimé ou disliké cette discussion

	disliked, err := api.GetDislikeFromOneComment(db, commentID, idUser)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal Server Error get like", http.StatusInternalServerError)
		return
	}

	if !disliked {
		// Si l'utilisateur a déjà disliké la discussion, supprimez le dislike
		err = api.DeleteLikeByUserIdForComment(db, commentID, idUser)
		if err != nil {
			http.Error(w, "Internal Server Error delete dislike", http.StatusInternalServerError)
			return
		}

		// Ajoutez un like
		err = api.SetDislikeComment(db, discussionIdInt, commentID, idUser)
		if err != nil {
			http.Error(w, "Internal Server Error set like", http.StatusInternalServerError)
			return
		}
	} else {
		// Si l'utilisateur a déjà aimé la discussion, supprimez le like
		err = api.DeleteDislikeByUserIdForComment(db, commentID, idUser)
		if err != nil {
			http.Error(w, "Internal Server Error delete like", http.StatusInternalServerError)
			return
		}
	}

	liens := "/discussion/" + discussionId

	// Redirigez l'utilisateur vers la page d'accueil après la mise à jour du like
	http.Redirect(w, r, liens, http.StatusSeeOther)
}
