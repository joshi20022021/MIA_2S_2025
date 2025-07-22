package main

// Importamos los paquetes necesarios
import (
	"net/http" // Provee constantes y funciones para manejar respuestas HTTP
	"strconv"  // Para convertir de string a int o viceversa

	"github.com/gin-gonic/gin" // Framework web Gin para construir APIs REST en Go
)

// Definimos una estructura (struct) llamada album, que modela un álbum musical
type album struct {
	ID     int    `json:"id"`     // ID del álbum (clave única), se serializa como "id" en JSON
	Title  string `json:"title"`  // Título del álbum
	Artist string `json:"artist"` // Nombre del artista
	Year   int    `json:"year"`   // Año de lanzamiento
}

// Creamos una lista (slice) de álbumes de ejemplo que servirá como nuestra "base de datos" en memoria
var albums = []album{
	{ID: 1, Title: "Abbey Road", Artist: "The Beatles", Year: 1969},
	{ID: 2, Title: "The Dark Side of the Moon", Artist: "Pink Floyd", Year: 1973},
	{ID: 3, Title: "Thriller", Artist: "Michael Jackson", Year: 1982},
	{ID: 4, Title: "Back in Black", Artist: "AC/DC", Year: 1980},
}

// Handler para GET /albums
// Devuelve todos los álbumes en formato JSON
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums) // Responde con status 200 y la lista de álbumes
}

// Handler para POST /albums
// Recibe un nuevo álbum en formato JSON, lo agrega a la lista y lo retorna
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Intentamos hacer un bind del JSON recibido al struct album
	if err := c.BindJSON(&newAlbum); err != nil {
		// Si ocurre un error, respondemos con status 400 y un mensaje de error
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Agregamos el nuevo álbum a la lista
	albums = append(albums, newAlbum)

	// Respondemos con status 201 (creado) y el álbum insertado
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// Handler para GET /albums/:id
// Busca un álbum por su ID y lo devuelve si lo encuentra
func getAlbumByID(c *gin.Context) {
	id := c.Param("id") // Obtenemos el parámetro de la URL como string

	// Recorremos todos los álbumes
	for _, a := range albums {
		// Comparamos el ID recibido con el ID del álbum convertido a string
		if id == strconv.Itoa(a.ID) {
			c.IndentedJSON(http.StatusOK, a) // Si coincide, respondemos con el álbum
			return
		}
	}

	// Si no se encuentra el álbum, respondemos con 404 y mensaje de error
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// Función principal: configura el servidor y las rutas
func main() {
	router := gin.Default() // Crea un nuevo router con middleware por defecto (logging y recuperación)

	// Definimos las rutas y sus funciones
	router.GET("/albums", getAlbums)        // Ruta para obtener todos los álbumes
	router.POST("/albums", postAlbums)      // Ruta para agregar un nuevo álbum
	router.GET("/albums/:id", getAlbumByID) // Ruta para obtener un álbum por ID

	// Iniciamos el servidor en localhost:8080
	router.Run("localhost:8080")
}
