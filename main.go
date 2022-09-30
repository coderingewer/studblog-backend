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
	router.HandleFunc("/api/users/new", middlewares.SetMiddlewareJSON(controllers.CreateUserByAdmin)).Methods("OPTIONS")
	router.HandleFunc("/api/users/confirm/{token}", middlewares.SetMiddlewareJSON(controllers.ConfirmAcoount)).Methods("POST")
	router.HandleFunc("/api/users/getAll", middlewares.SetMiddlewareJSON(controllers.GetUsers)).Methods("GET")
	router.HandleFunc("/api/users/getById/{id}", middlewares.SetMiddlewareJSON(controllers.GetUser)).Methods("GET")
	router.HandleFunc("/api/users/getByToken/", middlewares.SetMiddlewareJSON(controllers.GetUserByToken)).Methods("GET")
	router.HandleFunc("/api/users/delete/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUser)).Methods("DELETE")
	router.HandleFunc("/api/users/deleteByAdmin/{id}", middlewares.SetMiddlewareAuthentication(controllers.DeleteUserByAdmin)).Methods("DELETE")
	router.HandleFunc("/api/users/update/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateUser))).Methods("POST")
	router.HandleFunc("/api/users/updatePassword/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePassword))).Methods("POST")
	router.HandleFunc("/api/users/updatePasswordByAdmin/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdatePasswordByAdmin))).Methods("POST")
	router.HandleFunc("/api/users/login", middlewares.SetMiddlewareJSON(controllers.Login)).Methods("POST")
	router.HandleFunc("/api/users/updateuserimage", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controllers.UpdateUsermage))).Methods("POST")

	//Posts
	router.HandleFunc("/api/posts/new", middlewares.SetMiddlewareJSON(controllers.CreatePost)).Methods("POST")
	router.HandleFunc("/api/posts/getAll", middlewares.SetMiddlewareJSON(controllers.GetPosts)).Methods("GET")
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "*"}), handlers.AllowedMethods([]string{"POST", "PUT", "GET", "DELETE", "PATCH", "HEAD", "OPTIONS", "*"}), handlers.AllowedOrigins([]string{"*", "https://studblog.netlify.app"}))(router)))

}
