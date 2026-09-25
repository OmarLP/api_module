package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/OmarLP/api_module/internal/middleware"
	"github.com/OmarLP/api_module/internal/user"
	"github.com/OmarLP/api_module/pkg/authorization"
	"github.com/OmarLP/api_module/pkg/bootstrap"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

// embed de la carpeta web completa
//
//go:embed web/*
var webFS embed.FS

func main() {
	// cargando los cartificados
	err := authorization.LoadFiles("cmd/certificates/private.pem", "cmd/certificates/public.pem")
	if err != nil {
		log.Fatalf("Error al cargar los certificados: %v", err)
	}

	// cargar las variables de entorno desde el archivo .env
	_ = godotenv.Load()

	// inicializamos el logger para mensajes de log
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	// extraer el subdirectorio web para no exponer la ruta frontend en la url
	webContent, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal("error loading static files")
	}

	// servidor de archivos estáticos
	fileServer := http.FileServer(http.FS(webContent))

	// registrar la ruta
	http.Handle("/", fileServer)

	userRepo := user.NewRepository(l, db)
	userService := user.NewService(l, userRepo)
	userEndpoint := user.MakeEndpoints(userService)

	//----
	// manejar unicamente assets css, js, imagenes en /static
	staticContent, err := fs.Sub(webContent, "static")
	if err != nil {
		log.Fatal("Error al sub-extraer la carpeta static:", err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticContent))))
	// mapear rutas limpias
	http.HandleFunc("GET /login", serveHTML(webContent, "login.html"))
	http.HandleFunc("GET /profile", serveHTML(webContent, "profile.html"))
	http.HandleFunc("GET /register", serveHTML(webContent, "register.html"))
	http.HandleFunc("GET /reset-password", serveHTML(webContent, "reset-password.html"))
	http.HandleFunc("GET /update-password", serveHTML(webContent, "update-password.html"))
	// mapear la raiz a login
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	//--
	// limitar /login a 5 peticiones por minuto por IP
	loginLimiter := middleware.RateLimitMiddleware(rate.Every(12*time.Second), 5)

	http.HandleFunc("POST /api/login", loginLimiter(http.HandlerFunc(userEndpoint.Login)))
	http.HandleFunc("POST /api/CreateUsers", userEndpoint.CreateUser)
	http.HandleFunc("PATCH /api/UpdatePassword", userEndpoint.UpdatePassword)
	http.HandleFunc("PATCH /api/ResetPassword", userEndpoint.ResetPassword)
	http.HandleFunc("POST /api/refresh", userEndpoint.RefreshToken)
	http.HandleFunc("GET /api/profile", middleware.AuthMiddleware(http.HandlerFunc(userEndpoint.Profile)))
	http.HandleFunc("POST /api/logout", middleware.AuthMiddleware(http.HandlerFunc(userEndpoint.Logout)))

	// puedo mostrar un mensaje de conección exitosa con el nombre de la base de datos
	fmt.Printf("Conected to database: %v - %s\n", db, db.Migrator().CurrentDatabase())

	port := ":8080"
	address := "" + port
	fmt.Printf("Server is running on http://localhost%s\n", address)

	log.Fatal(http.ListenAndServe(address, nil))

}

// helper para servir archivos específicos
func serveHTML(fsys fs.FS, filePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := fsys.Open(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		http.ServeContent(w, r, filePath, time.Now(), file.(io.ReadSeeker))
	}
}
