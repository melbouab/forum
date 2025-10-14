package midemware

import (
	"net/http"

	"forum/server/database"
)

type HomeHandler func(http.ResponseWriter, *http.Request, *database.Repo, int)

func HomeMiddleware(next HomeHandler, repo *database.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// userID, isValid := lib.Help.CheckSession(r, repo)
		// if !isValid {
		// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
		// 	return
		// }
		next(w, r, repo, 0)
	}
}
