package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const PORT string = ":8080"

func main() {
	// Multiplexor de peticiones
	mux := http.NewServeMux()

	// Línea añadida para servir archivos estáticos:
	// Esto sirve los archivos dentro de /assets en la ruta /assets/
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Obtiene 4 imágenes al azar.
		images, err := getRandomImages("assets/images", 4)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// El path de cada imagen esté correctamente formateado para el cliente web.
		for i, imgPath := range images {
			images[i] = "/assets/images/" + filepath.Base(imgPath)
		}

		// Obtener el hostname
		hostname, err := os.Hostname()
		if err != nil {
			// Manejar el error, posiblemente configurando hostname a una cadena de reserva
			hostname = "Desconocido"
		}

		tmpl, err := template.ParseFiles("templates/template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Datos a pasar al template
		data := struct {
			Title        string
			Hostname     string
			Theme        string
			Images       []string
			CourseName   string
			Participants string
			Date         string
		}{
			Title:        "Servidor de imágenes",
			Hostname:     hostname,
			Theme:        "Ciudades del Quindío",
			Images:       images,
			CourseName:   "Computación en la Nube",
			Participants: "Juan Pablo Loaiza Nieto",
			Date:         "2024-1",
		}

		// Renderizar el template
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Printf("Servidor corriendo en http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, mux))
}

func getRandomImages(imagesDir string, n int) ([]string, error) {
	var images []string

	// Lee todos los archivos en el directorio y los añade a la lista de imágenes.
	err := filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Crea un generador de números aleatorios local
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Mezcla la lista de imágenes usando el generador local
	rng.Shuffle(len(images), func(i, j int) { images[i], images[j] = images[j], images[i] })

	if len(images) < n {
		return images, nil
	}

	// Selecciona las primeras n imágenes después de mezclar.
	return images[:n], nil
}
