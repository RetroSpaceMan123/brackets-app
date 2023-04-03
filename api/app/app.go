package app

import (
	"log"
	"net/http"

	"example.com/api/app/handler"
	"example.com/api/app/model"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
func (a *App) Initialize() {

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// Set all required routers
func (a *App) setRouters() {
	// Routing for handling the projects

	a.Post("/users/signup", a.SignUp)
	a.Put("/users/signin", a.SignIn) // Put operation because sessions
	a.Post("/users/googlesignup", a.GoogleSignUp)
	a.Put("/users/googlesignin", a.GoogleSignIn)
	a.Get("/users", a.GetAllUsers)
	a.Get("/users/{userid}", a.GetUser) // Changed from username to UUID as identifier
	a.Put("/users/{userid}", a.UpdateUser)
	a.Delete("/users/{userid}", a.DeleteUser)
	a.Get("/users/{userid}/brackets", a.GetUserBrackets)
	a.Post("/users/{userid}/brackets", a.CreateBracket)
	a.Get("/brackets", a.GetAllBrackets)
	a.Put("/users/{userid}/{bracketid}", a.UpdateBracket)
	a.Get("/users/{userid}/{bracketid}", a.GetBracket)
	a.Delete("/users/{userid}/{bracketid}", a.DeleteBracket)
	a.Post("/users/{userid}/{bracketid}/teams", a.AddTeam)
	a.Get("/users/{userid}/{bracketid}/teams", a.GetAllTeams)
	a.Get("/users/{userid}/{bracketid}/teams/{index}", a.GetTeam)
	a.Put("/users/{userid}/{bracketid}/teams/{index}", a.UpdateTeam)
	a.Delete("/users/{userid}/{bracketid}/teams/{index}", a.DeleteTeam)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Handlers for user data
func (a *App) SignUp(w http.ResponseWriter, r *http.Request) {
	handler.SignUp(a.DB, w, r)
}
func (a *App) SignIn(w http.ResponseWriter, r *http.Request) {
	handler.SignIn(a.DB, w, r)
}
func (a *App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	handler.GetAllUsers(a.DB, w, r)
}
func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	handler.GetUser(a.DB, w, r)
}
func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUser(a.DB, w, r)
}
func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	handler.DeleteUser(a.DB, w, r)
}

// Handlers for Google sign-in
func (a *App) GoogleSignUp(w http.ResponseWriter, r *http.Request) {
	handler.GoogleSignUp(a.DB, w, r)
}
func (a *App) GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	handler.GoogleSignIn(a.DB, w, r)
}

/*func (a *App) SignUp(w http.ResponseWriter, r *http.Request) {
	handler.SignUp(a.DB, w, r)
}*/

// Handlers for the bracket functions
func (a *App) GetBracket(w http.ResponseWriter, r *http.Request) {
	handler.GetBracket(a.DB, w, r)
}
func (a *App) GetUserBrackets(w http.ResponseWriter, r *http.Request) {
	handler.GetUserBrackets(a.DB, w, r)
}
func (a *App) GetAllBrackets(w http.ResponseWriter, r *http.Request) {
	handler.GetAllBrackets(a.DB, w, r)
}
func (a *App) CreateBracket(w http.ResponseWriter, r *http.Request) {
	handler.CreateBracket(a.DB, w, r)
}
func (a *App) UpdateBracket(w http.ResponseWriter, r *http.Request) {
	handler.UpdateBracket(a.DB, w, r)
}
func (a *App) DeleteBracket(w http.ResponseWriter, r *http.Request) {
	handler.DeleteBracket(a.DB, w, r)
}

// Handlers for user functions
func (a *App) AddTeam(w http.ResponseWriter, r *http.Request) {
	handler.AddTeam(a.DB, w, r)
}
func (a *App) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	handler.GetAllTeams(a.DB, w, r)
}
func (a *App) GetTeam(w http.ResponseWriter, r *http.Request) {
	handler.GetTeam(a.DB, w, r)
}
func (a *App) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	handler.UpdateTeam(a.DB, w, r)
}
func (a *App) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	handler.DeleteTeam(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
