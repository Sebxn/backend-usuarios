package main

import (
	"backend/api/routes"
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

func main() {
	// Inicializa la configuración de Firebase
	ctx := context.Background()
	opt := option.WithCredentialsFile(".env")
	config := &firebase.Config{ProjectID: "test-5eebf"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	// Inicializa el cliente de Firestore
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// Cerrar el cliente de Firestore cuando ya no se necesite
	defer client.Close()

	// Crear un enrutador utilizando gorilla/mux
	r := mux.NewRouter()
	routes.ConfigureRoutes(r, app)

	log.Fatal(http.ListenAndServe(":3000", r))
}
