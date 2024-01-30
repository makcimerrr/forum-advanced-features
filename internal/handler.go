package forum

import (
	"fmt"
	"forum/api"
	"html/template"
	"net/http"
	//"net/url"
	"strconv"
	"strings"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	
	tmpl := template.Must(template.ParseFiles("./web/templates/404.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error erreur 404", http.StatusInternalServerError)
		return
	}
}

func HandleServerError(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./web/templates/500.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error 500", http.StatusInternalServerError)
		return
	}
}

func HandleBadRequest(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./web/templates/400.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error 500", http.StatusInternalServerError)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifiez la validité de la session
    /* validSession, errMsg := isSessionValid(w, r)
    if !validSession {
        clearSessionCookies(w)
        // La session n'est pas valide, redirigez l'utilisateur vers la page de connexion ou effectuez d'autres actions
        http.Redirect(w, r, "/log_in?error="+url.QueryEscape(errMsg), http.StatusSeeOther)
        return
    } */

	// Récupérer l'URL de la requête
	url := r.URL.Path

	// Vérification du chemin d'accès
	if url != "/" && url != "/logOrSign" && url != "/log_in" && url != "/sign_up" && url != "/home" {
		// Appeler codeErreur pour la redirection
		codeErreur(w, r, url, "/", "./index.html")
		return
	}

	// Récupérer le nom d'utilisateur à partir du cookie "username"
	usernameCookie, err := r.Cookie("username")
	var username string

	if err == nil {
		username = usernameCookie.Value
	}

	var category []string
	var discussions []api.Discussion
	var categoryTitle []string

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
			return
		}
		category = r.Form["categories"]
	} else {
		var category2 []string
		tempCategory := r.URL.Query().Get("categories")
		category2 = strings.Split(tempCategory, ",")
		if category2[0] != "" {
			category = category2
		}
	}

	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error Open BDD", http.StatusInternalServerError)
		return
	}

	if len(category) == 0 {
		// Récupérer toutes les discussions à partir de la base de données

		discussions, err = api.GetAllDiscussions(db)
		if err != nil {
			http.Error(w, "Internal Server Error All Discussion", http.StatusInternalServerError)
			return
		}

		categoryTitle = append(categoryTitle, "Toutes les categories")

	} else {
		//recupe discussion via category
		var tempTabInt []int
		var TabInt []int

		for i := 0; i < len(category); i++ {

			categoryID, _ := strconv.Atoi(category[i])

			tempTabInt, err = api.GetDiscussionIDByCategoryID(db, categoryID)
			if err != nil {
				http.Error(w, "Internal Server Error All Discussion", http.StatusInternalServerError)
				return
			}

			TabInt = append(TabInt, tempTabInt...)
		}

		TabInt = supprimerDoublons(TabInt)

		var discussion api.Discussion
		var tempBoleen bool

		for i := 0; i < len(TabInt); i++ {

			for j := 0; j < len(category); j++ {

				categoryID, _ := strconv.Atoi(category[j])

				tempBoleen, err = api.CheckIfDiscussionCategoryOk(db, categoryID, TabInt[i])
				if err != nil {
					http.Error(w, "Internal Server Error check if la discussion a bien toute les categories", http.StatusInternalServerError)
					return
				}
				if !tempBoleen {
					break
				}
			}

			if tempBoleen {
				discussion, err = api.GetOneDiscussions(db, TabInt[i])
				if err != nil {
					http.Error(w, "Internal Server Error one discussion", http.StatusInternalServerError)
					return
				}
				discussions = append(discussions, discussion)
			}
		}

		var tempCategory string

		for i := 0; i < len(category); i++ {

			categoryID, _ := strconv.Atoi(category[i])

			tempCategory, err = api.GetCategoryByID(db, categoryID)
			if err != nil {
				http.Error(w, "Internal Server Error get category by id", http.StatusInternalServerError)
				return
			}
			categoryTitle = append(categoryTitle, tempCategory)
		}

	}

	// Récupérer les catégories uniques
	categories, err := api.GetCategory(db)
	if err != nil {
		http.Error(w, "Internal Server Error category", http.StatusInternalServerError)
		return
	}

	if username != "" {
		idUser, err := api.GetUserByUsername(db, username)
		if err != nil {
			http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
			return
		}

		for i := range discussions {
			liked, err := api.CheckIfUserLikedDiscussion(db, idUser, discussions[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check like", http.StatusInternalServerError)
				return
			}
			discussions[i].Liked = liked

			// Pour chaque discussion, vérifiez si l'utilisateur l'a pas aimée
			disliked, err := api.CheckIfUserDislikedDiscussion(db, idUser, discussions[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check dislike", http.StatusInternalServerError)
				return
			}
			discussions[i].Disliked = disliked

		}

	}

	// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
	for i := range discussions {

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
		numberLike, err := api.CheckNumberOfLikesForDiscussion(db, discussions[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
			return
		}

		discussions[i].NumberLike = numberLike

		numberDislike, err := api.CheckNumberOfDislikesForDiscussion(db, discussions[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
			return
		}

		discussions[i].NumberDislike = numberDislike

		numberComment, err := api.CheckNumberOfCommentForDiscussion(db, discussions[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of comment", http.StatusInternalServerError)
			return
		}

		discussions[i].NumberComment = numberComment

		categoriesID, err := api.GetCategoryIDByDiscussionID(db, discussions[i].ID)
		if err != nil {
			http.Error(w, "Error fetching get id category", http.StatusInternalServerError)
			return
		}

		var categories []string
		for i := 0; i < len(categoriesID); i++ {
			value, err := api.GetCategoryByID(db, categoriesID[i])
			if err != nil {
				http.Error(w, "Error fetching get category", http.StatusInternalServerError)
				return
			}

			categories = append(categories, value)
		}

		discussions[i].Category = categories
	}

	// Créer une structure de données pour passer les informations au modèle
	data := struct {
		Username      string
		Discussions   []api.Discussion
		Categories    []api.Categories
		CategoryTitle []string
	}{
		Username:      username,
		Discussions:   discussions,
		Categories:    categories,
		CategoryTitle: categoryTitle,
	}

	tmpl := template.Must(template.ParseFiles("./web/templates/index.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error template index", http.StatusInternalServerError)
		return
	}
}

func CreateDiscussionHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifiez la validité de la session
    /* validSession, errMsg := isSessionValid(w, r)
    if !validSession {
        clearSessionCookies(w)
        // La session n'est pas valide, redirigez l'utilisateur vers la page de connexion ou effectuez d'autres actions
        http.Redirect(w, r, "/log_in?error="+url.QueryEscape(errMsg), http.StatusSeeOther)
        return
    } */

	// Obtenez le nom d'utilisateur à partir du cookie "username"
	usernameCookie, err := r.Cookie("username")
	if err != nil {
		// Gérer l'erreur ici, par exemple, en redirigeant l'utilisateur vers une page de connexion s'il n'est pas connecté.
		http.Redirect(w, r, "/log_in", http.StatusSeeOther)
		return
	}
	username := usernameCookie.Value

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		message := r.FormValue("message")
		categories := r.Form["categories"]

		if len(categories) == 0 {
			errors := "At least one category must be selected"

			db, err := api.OpenBDD()
			if err != nil {
				http.Error(w, "Internal Server Error Open BDD", http.StatusInternalServerError)
				return
			}

			category, err := api.GetCategory(db)
			if err != nil {
				http.Error(w, "Internal Server Error category", http.StatusInternalServerError)
				return
			}

			data := struct {
				Username string
				Category []api.Categories
				Error    string
				Title    string
				Message  string
			}{
				Username: username,
				Category: category,
				Error:    errors,
				Title:    title,
				Message:  message,
			}

			// Affichez la page pour écrire une discussion (write_discussion.html) avec le menu déroulant
			tmpl := template.Must(template.ParseFiles("./web/templates/write_discussion.html"))
			tmpl.Execute(w, data)
		} else {

			// Ouverture de la connexion à la base de données
			db, err := api.OpenBDD()
			if err != nil {
				http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
				return
			}

			//convertir username par id

			idUser, err := api.GetUserByUsername(db, username)
			if err != nil {
				http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
				return
			}

			// Insérez la nouvelle discussion dans la base de données, y compris la catégorie

			err = api.SetDiscussion(db, idUser, title, message)
			if err != nil {
				http.Error(w, "Internal Server Error set discussion", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

			// Récupérez l'ID de la discussion nouvellement créée
			var discussionID int
			err = db.QueryRow("SELECT last_insert_rowid()").Scan(&discussionID)
			if err != nil {
				http.Error(w, "Internal Server Error last id discussion", http.StatusInternalServerError)
				return
			}

			var categoryID int

			for i := 0; i < len(categories); i++ {

				categoryID, err = strconv.Atoi(categories[i])
				err = api.SetDiscussionCategory(db, discussionID, categoryID)
				if err != nil {
					http.Error(w, "Internal Server Error set discussion_category", http.StatusInternalServerError)
					fmt.Println(err)
					return
				}
			}

			// Redirigez l'utilisateur vers la page de la discussion
			http.Redirect(w, r, fmt.Sprintf("/discussion/%d", discussionID), http.StatusSeeOther)
			return
		}
	} else {

		db, err := api.OpenBDD()
		if err != nil {
			http.Error(w, "Internal Server Error Open BDD", http.StatusInternalServerError)
			return
		}

		category, err := api.GetCategory(db)
		if err != nil {
			http.Error(w, "Internal Server Error category", http.StatusInternalServerError)
			return
		}

		data := struct {
			Username string
			Category []api.Categories
			Error    string
			Title    string
			Message  string
		}{
			Username: username,
			Category: category,
			Error:    "",
			Title:    "",
			Message:  "",
		}

		// Affichez la page pour écrire une discussion (write_discussion.html) avec le menu déroulant
		tmpl := template.Must(template.ParseFiles("./web/templates/write_discussion.html"))
		tmpl.Execute(w, data)
	}
}

func ShowDiscussionHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifiez la validité de la session
    /* validSession, errMsg := isSessionValid(w, r)
    if !validSession {
        clearSessionCookies(w)
        // La session n'est pas valide, redirigez l'utilisateur vers la page de connexion ou effectuez d'autres actions
        http.Redirect(w, r, "/log_in?error="+url.QueryEscape(errMsg), http.StatusSeeOther)
        return
    } */


	// Récupérer le nom d'utilisateur à partir du cookie "username"
	usernameCookie, err := r.Cookie("username")
	var username string

	if err == nil {
		username = usernameCookie.Value
	}
	// Récupérez l'ID de la discussion à partir de l'URL
	discussionID := r.URL.Path[len("/discussion/"):]
	// Convertissez l'ID de la discussion en un entier
	discussionIDInt, err := strconv.Atoi(discussionID)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Ouverture de la connexion à la base de données
	db, err := api.OpenBDD()
	if err != nil {
		http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
		return
	}

	var discussions api.Discussion

	// Effectuez une requête SQL pour récupérer les détails de la discussion en fonction de l'ID
	discussions, err = api.GetOneDiscussions(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Discussion not found", http.StatusNotFound)
		return
	}

	// Effectuez une autre requête SQL pour récupérer les commentaires associés à cette discussion
	comment, err := api.GetCommentsFromDiscussion(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	categoriesID, err := api.GetCategoryIDByDiscussionID(db, discussionIDInt)
	if err != nil {
		http.Error(w, "Error fetching get id category", http.StatusInternalServerError)
		return
	}

	var categories []string
	for i := 0; i < len(categoriesID); i++ {
		value, err := api.GetCategoryByID(db, categoriesID[i])
		if err != nil {
			http.Error(w, "Error fetching get category", http.StatusInternalServerError)
			return
		}

		categories = append(categories, value)
	}

	discussions.Category = categories

	if username != "" {
		idUser, err := api.GetUserByUsername(db, username)
		if err != nil {
			http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
			return
		}

		liked, err := api.CheckIfUserLikedDiscussion(db, idUser, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check like", http.StatusInternalServerError)
			return
		}
		discussions.Liked = liked

		// Pour chaque discussion, vérifiez si l'utilisateur l'a pas aimée
		disliked, err := api.CheckIfUserDislikedDiscussion(db, idUser, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check dislike", http.StatusInternalServerError)
			return
		}
		discussions.Disliked = disliked

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
		numberLike, err := api.CheckNumberOfLikesForDiscussion(db, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
			return
		}

		discussions.NumberLike = numberLike

		numberDislike, err := api.CheckNumberOfDislikesForDiscussion(db, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
			return
		}

		discussions.NumberDislike = numberDislike

		numberComment, err := api.CheckNumberOfCommentForDiscussion(db, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of comment", http.StatusInternalServerError)
			return
		}

		discussions.NumberComment = numberComment

		for i := range comment {
			liked, err := api.CheckIfUserLikedComment(db, idUser, comment[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check like", http.StatusInternalServerError)
				return
			}
			comment[i].Liked = liked

			// Pour chaque comment, vérifiez si l'utilisateur l'a pas aimée
			disliked, err := api.CheckIfUserDislikeComment(db, idUser, comment[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check dislike", http.StatusInternalServerError)
				return
			}
			comment[i].Disliked = disliked

			// Pour chaque comment, vérifiez si l'utilisateur l'a aimée
			numberLike, err := api.CheckNumberOfLikesForComment(db, comment[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
				return
			}

			comment[i].NumberLike = numberLike

			numberDislike, err := api.CheckNumberOfDislikeForComment(db, comment[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
				return
			}

			comment[i].NumberDislike = numberDislike
		}

	} else {

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
		numberLike, err := api.CheckNumberOfLikesForDiscussion(db, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
			return
		}

		discussions.NumberLike = numberLike

		numberDislike, err := api.CheckNumberOfDislikesForDiscussion(db, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
			return
		}

		discussions.NumberDislike = numberDislike

		numberComment, err := api.CheckNumberOfCommentForDiscussion(db, discussions.ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of comment", http.StatusInternalServerError)
			return
		}

		discussions.NumberComment = numberComment

		// Pour chaque coment, vérifiez si l'utilisateur l'a aimée
		for i := range comment {

			// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
			numberLike, err := api.CheckNumberOfLikesForComment(db, comment[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
				return
			}

			comment[i].NumberLike = numberLike

			numberDislike, err := api.CheckNumberOfDislikeForComment(db, comment[i].ID)
			if err != nil {
				http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
				return
			}

			comment[i].NumberDislike = numberDislike

		}
	}

	// Créez une structure de données pour stocker les détails de la discussion et les commentaires
	data := struct {
		Username   string
		Discussion api.Discussion
		Comments   []api.Comment
	}{
		Username:   username,
		Discussion: discussions,
		Comments:   comment,
	}

	/* 	var filter sql.NullString
	   	err = db.QueryRow("SELECT COALESCE(filter, '') FROM discussion_user WHERE id = ?", discussionIDInt).Scan(&filter)
	   	if err != nil {
	   		http.Error(w, "Filter not found", http.StatusNotFound)
	   		return
	   	}

	   	var filterValue string
	   	if filter.Valid {
	   		filterValue = filter.String
	   	}

	   	data.Filter = &filterValue // Permet l'affiche du filtre {{ .Filter}}
	*/

	// Affichez les détails de la discussion et les commentaires dans un modèle HTML
	tmpl := template.Must(template.ParseFiles("./web/templates/show_discussion.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error vers show discussion", http.StatusInternalServerError)
		return
	}
}

func LogOrSignHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifiez la validité de la session
    /* validSession, errMsg := isSessionValid(w, r)
    if !validSession {
        clearSessionCookies(w)
        // La session n'est pas valide, redirigez l'utilisateur vers la page de connexion ou effectuez d'autres actions
        http.Redirect(w, r, "/log_in?error="+url.QueryEscape(errMsg), http.StatusSeeOther)
        return
    } */

	tmpl := template.Must(template.ParseFiles("./web/templates/logOrSign.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error vers logOrSign", http.StatusInternalServerError)
		return
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifiez la validité de la session
    /* validSession, errMsg := isSessionValid(w, r)
    if !validSession {
        clearSessionCookies(w)
        // La session n'est pas valide, redirigez l'utilisateur vers la page de connexion ou effectuez d'autres actions
        http.Redirect(w, r, "/log_in?error="+url.QueryEscape(errMsg), http.StatusSeeOther)
        return
    } */

	var formError []string

	if r.Method == http.MethodPost {
		// Récupération des informations du formulaire
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hashpass := hash(password)

		db, err := api.OpenBDD()
		if err != nil {
			http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
			return
		}

		// Vérification si le nom d'utilisateur est déjà utilisé
		var existingUsername string
		err = db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&existingUsername)
		if err == nil {
			formError = append(formError, "This Username Is Already Use !! ")
		}

		// Vérification si l'e-mail est déjà utilisé
		var existingEmail string
		err = db.QueryRow("SELECT email FROM user WHERE email = ?", email).Scan(&existingEmail)
		if err == nil {
			formError = append(formError, "This Email Is Already Use !!")
		}

		if formError == nil {

			err = api.SetUser(db, username, email, hashpass)
			if err != nil {
				http.Error(w, "Internal Server Error set user", http.StatusInternalServerError)
				return
			}

			err := CreateAndSetSessionCookies(w, username)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Rediriger l'utilisateur vers la page "/home" après l'enregistrement
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("./web/templates/sign_up.html"))
	err := tmpl.Execute(w, formError)
	if err != nil {
		http.Error(w, "Internal Server Error vers sign_up", http.StatusInternalServerError)
		return
	}
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifiez la validité de la session
    /* validSession, errMsg := isSessionValid(w, r)
    if !validSession {
        clearSessionCookies(w)
        // La session n'est pas valide, redirigez l'utilisateur vers la page de connexion ou effectuez d'autres actions
        http.Redirect(w, r, "/log_in?error="+url.QueryEscape(errMsg), http.StatusSeeOther)
        return
    } */

	var formError []string

	if r.Method == http.MethodPost {
		loginemail := r.FormValue("loginemail")
		loginpassword := r.FormValue("loginpassword")

		// Ouverture de la connexion à la base de données
		db, err := api.OpenBDD()
		if err != nil {
			http.Error(w, "Internal Server Error open BDD", http.StatusInternalServerError)
			return
		}

		var trueemail string
		var truepassword uint32
		var username string
		err = db.QueryRow("SELECT username, email, password FROM user WHERE email = ?", loginemail).Scan(&username, &trueemail, &truepassword)
		if err != nil {
			formError = append(formError, "Email Doesn't exist .")
		} else {
			hashloginpassword := hash(loginpassword)

			// Vérifier le mot de passe
			if hashloginpassword != truepassword {
				formError = append(formError, "Password Failed.")
			} else {
				// L'utilisateur est connecté avec succès
				err := CreateAndSetSessionCookies(w, username)
				if err != nil {
					fmt.Println(err)

					return
				}

				// Redirigez l'utilisateur vers la page "/"
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
	}

	tmpl := template.Must(template.ParseFiles("./web/templates/login.html"))
	err := tmpl.Execute(w, formError)
	if err != nil {
		http.Error(w, "Internal Server Error vers login", http.StatusInternalServerError)
		return
	}
}

func ProfilHandler(w http.ResponseWriter, r *http.Request) {

	// Vérifiez la validité de la session
    /* validSession, errMsg := isSessionValid(w, r)
    if !validSession {
        clearSessionCookies(w)
        // La session n'est pas valide, redirigez l'utilisateur vers la page de connexion ou effectuez d'autres actions
        http.Redirect(w, r, "/log_in?error="+url.QueryEscape(errMsg), http.StatusSeeOther)
        return
    } */


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
		http.Error(w, "Internal Server Error Open BDD", http.StatusInternalServerError)
		return
	}

	idUser, err := api.GetUserByUsername(db, username)
	if err != nil {
		http.Error(w, "Internal Server Error get id by username", http.StatusInternalServerError)
		return
	}

	var discussionsCreated []api.Discussion

	discussionsCreated, err = api.GetAllDiscussionsForOneUser(db, idUser)
	if err != nil {
		http.Error(w, "Internal Server Error All Discussion", http.StatusInternalServerError)
		return
	}

	for i := range discussionsCreated {
		/* liked, err := api.CheckIfUserLikedDiscussion(db, idUser, discussionsCreated[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check like", http.StatusInternalServerError)
			return
		}
		discussionsCreated[i].Liked = liked

		// Pour chaque discussion, vérifiez si l'utilisateur l'a pas aimée
		disliked, err := api.CheckIfUserDislikedDiscussion(db, idUser, discussionsCreated[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check dislike", http.StatusInternalServerError)
			return
		}
		discussionsCreated[i].Disliked = disliked */

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
		numberLike, err := api.CheckNumberOfLikesForDiscussion(db, discussionsCreated[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
			return
		}

		discussionsCreated[i].NumberLike = numberLike

		numberDislike, err := api.CheckNumberOfDislikesForDiscussion(db, discussionsCreated[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
			return
		}

		discussionsCreated[i].NumberDislike = numberDislike

		numberComment, err := api.CheckNumberOfCommentForDiscussion(db, discussionsCreated[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of comment", http.StatusInternalServerError)
			return
		}

		discussionsCreated[i].NumberComment = numberComment

		categoriesID, err := api.GetCategoryIDByDiscussionID(db, discussionsCreated[i].ID)
		if err != nil {
			http.Error(w, "Error fetching get id category", http.StatusInternalServerError)
			return
		}

		var categories []string
		for i := 0; i < len(categoriesID); i++ {
			value, err := api.GetCategoryByID(db, categoriesID[i])
			if err != nil {
				http.Error(w, "Error fetching get category", http.StatusInternalServerError)
				return
			}

			categories = append(categories, value)
		}

		discussionsCreated[i].Category = categories

	}

	var commentCreated []api.Comment

	commentCreated, err = api.GetAllCommentForOneUser(db, idUser)
	if err != nil {
		http.Error(w, "Internal Server Error All Discussion", http.StatusInternalServerError)
		return
	}

	for i := range commentCreated {

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
		numberLike, err := api.CheckNumberOfLikesForComment(db, commentCreated[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
			return
		}

		commentCreated[i].NumberLike = numberLike

		numberDislike, err := api.CheckNumberOfDislikeForComment(db, commentCreated[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
			return
		}

		commentCreated[i].NumberDislike = numberDislike

	}

	var discussionsLiked []api.Discussion

	var tempId []int
	//recupérer les id de discussion où utilisater a liker
	tempId, err = api.GetDiscussionIdByLikeForOneUser(db, idUser)
	if err != nil {
		http.Error(w, "Internal Server Error der id Discussion by like for ine user", http.StatusInternalServerError)
		return
	}

	//recupérer les discussion en grace au id
	var discussion api.Discussion
	for i := 0; i < len(tempId); i++ {
		discussion, err = api.GetOneDiscussions(db, tempId[i])
		fmt.Println(err)
		if err != nil {
			http.Error(w, "Internal Server Error get one discussion for like", http.StatusInternalServerError)
			return
		}
		discussionsLiked = append(discussionsLiked, discussion)
	}

	for i := range discussionsLiked {
		/* liked, err := api.CheckIfUserLikedDiscussion(db, idUser, discussionsLiked[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check like", http.StatusInternalServerError)
			return
		}
		discussionsLiked[i].Liked = liked

		// Pour chaque discussion, vérifiez si l'utilisateur l'a pas aimée
		disliked, err := api.CheckIfUserDislikedDiscussion(db, idUser, discussionsLiked[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check dislike", http.StatusInternalServerError)
			return
		}
		DiscussionsLiked[i].Disliked = disliked */

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
		numberLike, err := api.CheckNumberOfLikesForDiscussion(db, discussionsLiked[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
			return
		}

		discussionsLiked[i].NumberLike = numberLike

		numberDislike, err := api.CheckNumberOfDislikesForDiscussion(db, discussionsLiked[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
			return
		}

		discussionsLiked[i].NumberDislike = numberDislike

		numberComment, err := api.CheckNumberOfCommentForDiscussion(db, discussionsLiked[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of comment", http.StatusInternalServerError)
			return
		}

		discussionsLiked[i].NumberComment = numberComment

		categoriesID, err := api.GetCategoryIDByDiscussionID(db, discussionsLiked[i].ID)
		if err != nil {
			http.Error(w, "Error fetching get id category", http.StatusInternalServerError)
			return
		}

		var categories []string
		for i := 0; i < len(categoriesID); i++ {
			value, err := api.GetCategoryByID(db, categoriesID[i])
			if err != nil {
				http.Error(w, "Error fetching get category", http.StatusInternalServerError)
				return
			}

			categories = append(categories, value)
		}

		discussionsCreated[i].Category = categories

	}

	var commentLiked []api.Comment

	//recupérer les id de discussion où utilisater a liker
	tempId, err = api.GetCommentIdByLikeForOneUser(db, idUser)
	if err != nil {
		http.Error(w, "Internal Server Error der id Discussion by like for ine user", http.StatusInternalServerError)
		return
	}

	//recupérer les discussion en grace au id
	var comment api.Comment
	for i := 0; i < len(tempId); i++ {
		comment, err = api.GetOneCommentById(db, tempId[i])
		if err != nil {
			http.Error(w, "Internal Server Error get one discussion", http.StatusInternalServerError)
			return
		}
		commentLiked = append(commentLiked, comment)
	}

	for i := range commentLiked {

		// Pour chaque discussion, vérifiez si l'utilisateur l'a aimée
		numberLike, err := api.CheckNumberOfLikesForComment(db, commentLiked[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of like", http.StatusInternalServerError)
			return
		}

		commentLiked[i].NumberLike = numberLike

		numberDislike, err := api.CheckNumberOfDislikeForComment(db, commentLiked[i].ID)
		if err != nil {
			http.Error(w, "Internal Server Error check number of dislike", http.StatusInternalServerError)
			return
		}

		commentLiked[i].NumberDislike = numberDislike

	}

	// Créer une structure de données pour passer les informations au modèle
	data := struct {
		Username            string
		DiscussionsCreated  []api.Discussion
		DiscussionsLiked    []api.Discussion
		DiscussionsDisliked []api.Discussion
		CommentCreated      []api.Comment
		CommentLiked        []api.Comment
		CommentDisliked     []api.Comment
	}{
		Username:           username,
		DiscussionsCreated: discussionsCreated,
		DiscussionsLiked:   discussionsLiked,
		CommentCreated:     commentCreated,
		CommentLiked:       commentLiked,
	}

	tmpl := template.Must(template.ParseFiles("./web/templates/profil.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error template index", http.StatusInternalServerError)
		return
	}
}


