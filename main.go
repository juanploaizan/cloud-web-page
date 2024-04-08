package main

import (
	"flag"          // Usado para el parsing de argumentos de línea de comando
	"fmt"           // Operaciones de entrada/salida formateadas (por ejemplo, imprimir en consola)
	"html/template" // Generación de HTML basado en templates
	"log"           // Logging de errores y otra información
	"math/rand"     // Generación de números aleatorios
	"net/http"      // Construcción de servidores HTTP
	"os"            // Interacción con el sistema operativo (OS)
	"path/filepath" // Manipulación y búsqueda de rutas de archivos
	"time"          // Operaciones relacionadas con el tiempo (por ejemplo, generar semillas aleatorias)
)

// Variables globales para los parámetros de línea de comandos, permitiendo configurar el comportamiento de la aplicación al iniciarla.
var (
	port         string // Puerto del servidor web
	title        string // Título de la página web
	theme        string // Tema de las imágenes a mostrar
	courseName   string // Nombre del curso para mostrar en la página
	participants string // Lista de participantes del curso
	date         string // Fecha del curso
)

func init() {
	// Inicializa los parámetros de línea de comando con valores predeterminados y descripciones.
	// Estos valores pueden ser sobrescritos al ejecutar el programa.
	flag.StringVar(&port, "port", "8080", "Puerto para el servidor web")
	flag.StringVar(&title, "title", "Servidor de imágenes", "Título de la página")
	flag.StringVar(&theme, "theme", "Ciudades del Quindío", "Tema de las imágenes")
	flag.StringVar(&courseName, "courseName", "Computación en la Nube", "Nombre del curso")
	flag.StringVar(&participants, "participants", "Juan Pablo Loaiza Nieto", "Participantes del curso")
	flag.StringVar(&date, "date", "2024-1", "Fecha del curso")
}

func main() {
	flag.Parse() // Parsea los argumentos de línea de comando.

	// Configura el multiplexor HTTP para manejar las rutas.
	mux := http.NewServeMux()

	// Manejador para servir archivos estáticos desde el directorio /assets/.
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Manejador para la ruta raíz. Genera y sirve la página principal con imágenes aleatorias.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		images, err := getRandomImages("assets/images", 4)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepara las rutas de las imágenes para ser usadas en la web.
		for i, imgPath := range images {
			images[i] = "/assets/images/" + filepath.Base(imgPath)
		}

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "Desconocido"
		}

		// Parsea el archivo template.html para generar la página web.
		tmpl, err := template.ParseFiles("templates/template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Datos a pasar al template HTML.
		data := struct {
			Title        string
			Hostname     string
			Theme        string
			Images       []string
			CourseName   string
			Participants string
			Date         string
		}{
			Title:        title,
			Hostname:     hostname,
			Theme:        theme,
			Images:       images,
			CourseName:   courseName,
			Participants: participants,
			Date:         date,
		}

		// Ejecuta el template con los datos proporcionados.
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Inicia el servidor y escucha en el puerto configurado.
	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// getRandomImages selecciona un número específico de imágenes aleatorias de un directorio.
func getRandomImages(imagesDir string, n int) ([]string, error) {
	var images []string

	// Recorre todos los archivos en el directorio de imágenes.
	err := filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Genera números aleatorios basados en la hora actual para mezclar las imágenes.
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Mezcla las imágenes de manera aleatoria.
	rng.Shuffle(len(images), func(i, j int) { images[i], images[j] = images[j], images[i] })

	// Selecciona las primeras n imágenes después de mezclar.
	if len(images) < n {
		return images, nil
	}

	return images[:n], nil
}
