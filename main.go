package main

import (
	"log"
	"net/http"
	"os"
	"studapp-blog/api/controllers"
	"studapp-blog/api/utils"
	"studapp-blog/middlewares"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, "Welcome to Studblog API")
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", middlewares.SetMiddlewareJSON(Hello)).Methods("GET")

	//Users
	router.HandleFunc("/api/users/new", middlewares.SetMiddlewareJSON(controllers.CreateUser)).Methods("POST")
	router.HandleFunc("/api/users/confirm/{token}", middlewares.SetMiddlewareJSON(controllers.ConfirmAcoount)).Methods("POST")
	router.HandleFunc("/api/users/getAll", middlewares.SetMiddlewareJSON(controllers.GetUsers)).Methods("GET")
	router.HandleFunc("/api/users/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetUser)).Methods("GET")
	router.HandleFunc("/api/users/getByUserName/{username}", middlewares.SetMiddlewareJSON(controllers.GetUserByUserName)).Methods("GET")
	router.HandleFunc("/api/users/getByToken/", middlewares.SetMiddlewareJSON(controllers.GetUserByToken)).Methods("GET")
	router.HandleFunc("/api/users/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUser)).Methods("DELETE")
	router.HandleFunc("/api/users/deleteByAdmin/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUserByAdmin)).Methods("DELETE")
	router.HandleFunc("/api/users/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateUser))).Methods("POST")
	router.HandleFunc("/api/users/updatePassword/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePassword))).Methods("POST")
	router.HandleFunc("/api/users/updatePasswordByAdmin/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePasswordByAdmin))).Methods("POST")
	router.HandleFunc("/api/users/login", middlewares.SetMiddlewareJSON(controllers.Login)).Methods("POST")
	router.HandleFunc("/api/users/updateuserimage", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateUsermage))).Methods("POST")
	router.HandleFunc("/api/users/updatebyAdmin/{email}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateUserByAdmin))).Methods("POST")

	//Posts
	router.HandleFunc("/api/posts/new", middlewares.SetMiddlewareJSON(controllers.CreatePost)).Methods("POST")
	router.HandleFunc("/api/posts/getAll", middlewares.SetMiddlewareJSON(controllers.GetPosts)).Methods("GET")
	router.HandleFunc("/api/posts/getPopulars", middlewares.SetMiddlewareJSON(controllers.GetPopularPosts)).Methods("GET")
	router.HandleFunc("/api/posts/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetPost)).Methods("GET")
	router.HandleFunc("/api/posts/getByUserId/{userId}", middlewares.SetMiddlewareJSON(controllers.GetPostsByUserID)).Methods("GET")
	router.HandleFunc("/api/posts/getByCategory/{category}", middlewares.SetMiddlewareJSON(controllers.GetPostsByCategory)).Methods("GET")
	router.HandleFunc("/api/posts/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeletePost)).Methods("DELETE")
	router.HandleFunc("/api/posts/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePost))).Methods("POST")
	router.HandleFunc("/api/posts/uploadimg/{postId}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePostImage))).Methods("POST")
	router.HandleFunc("/api/posts/likepost/{id}", middlewares.SetMiddlewareJSON(controllers.LikePost)).Methods("POST")
	router.HandleFunc("/api/posts/unlikepost/{id}", middlewares.SetMiddlewareJSON(controllers.UnLikePost)).Methods("DELETE")
	router.HandleFunc("/api/posts/view/{id}", middlewares.SetMiddlewareJSON(controllers.ViewPost)).Methods("POST")

	//Images
	router.HandleFunc("/api/images/upload", middlewares.SetMiddlewareJSON(controllers.ImgUpload)).Methods("POST")
	router.HandleFunc("/api/images/update/{imageId}", middlewares.SetMiddlewareJSON(controllers.UpdateImage)).Methods("POST")
	router.HandleFunc("/api/images/getAll", middlewares.SetMiddlewareJSON(controllers.GetAllImages)).Methods("GET")
	//fav lists
	router.HandleFunc("/api/favs/new", middlewares.SetMiddlewareJSON(controllers.NewFavsList)).Methods("POST")
	router.HandleFunc("/api/favs/additem", middlewares.SetMiddlewareJSON(controllers.AddItemToList)).Methods("POST")
	router.HandleFunc("/api/favs/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteFavsList)).Methods("DELETE")
	router.HandleFunc("/api/favs/removeitem/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteFavFromList)).Methods("DELETE")
	router.HandleFunc("/api/favs/getByUserId/{userId}", middlewares.SetMiddlewareJSON(controllers.GetFavsByUserId)).Methods("GET")
	router.HandleFunc("/api/favs/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetFavsById)).Methods("GET")
	router.HandleFunc("/api/favs/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateFavsList))).Methods("POST")

	//Comments
	router.HandleFunc("/api/comments/new", middlewares.SetMiddlewareJSON(controllers.CreateComment)).Methods("POST")
	router.HandleFunc("/api/comments/update/{id}", middlewares.SetMiddlewareJSON(controllers.UpdateComment)).Methods("POST")
	router.HandleFunc("/api/comments/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteComment)).Methods("DELETE")
	router.HandleFunc("/api/comments/getByUserID/{userId}", middlewares.SetMiddlewareAuthentication(controllers.GetCommentsByUserID)).Methods("GET")
	router.HandleFunc("/api/comments/getByPostId/{postId}", middlewares.SetMiddlewareJSON(controllers.GetCommentsByPostID)).Methods("GET")
	router.HandleFunc("/api/comments/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetComment)).Methods("GET")
	router.HandleFunc("/api/comments/getAll", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.GetComments))).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	//handler := cors.AllowAll().Handler(router)
	//log.Fatal(http.ListenAndServe(":"+port, handler))

	log.Fatal(http.ListenAndServe(":"+port,
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"POST", "PUT", "GET", "DELETE", "PATCH", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{
				"https://www.studappblog.com", "http://www.studappblog.com",
				"www.studappblog.com", "studappblog.com",
				"https://studblog-demo-2.netlify.app", "studblog-demo-2.netlify.app",
				"localhost:3000", "http://localhost:3000"}))(router)))

}
