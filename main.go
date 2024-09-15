package main

import (
	"encoding/json"
	"fmt"
	"io" // para controlar las entradas/salidas
	"log"
	"net/http"
	"strconv" // conversor de strings

	"github.com/gorilla/mux"
)

type task struct {
	Id      int    `json:"Id"`   // respondo un json con el campo "Id"
	Name    string `json:"Name"` // LOS NOMBRES TIENEN QUE TENER LA PRIMER LETRA MAYUSCULA
	Content string `json:"Content"`
}

type tasksList []task

var tasks = tasksList{
	{
		Id:      1,
		Name:    "Task 1",
		Content: "content 1",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) { // w -> es lo que mi SV responde al user y r -> soli del user (lo que el user manda a mi SV)
	fmt.Fprintf(w, "welcome to API")
}

func showTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // creo un header por c/ respuesta indicando el formato de respuesta
	json.NewEncoder(w).Encode(tasks)                   // devuelve las tareas
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task // variable vacia al principio de tipo "task" por lo que tendra Id, Name y Content

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid task!")
	}

	json.Unmarshal(reqBody, &newTask) // lo que venga en el reqBody lo asigno a la variable newTask con "&"

	newTask.Id = len(tasks) + 1    // asigno un Id personalizado. Usando el índice del elemento en el array
	tasks = append(tasks, newTask) // agrego la nueva task al array de tasks

	w.Header().Set("Content-Type", "application/json") // creo un header por c/ respuesta indicando el formato de respuesta
	w.WriteHeader(http.StatusCreated)                  // mando un status de "operacion correcta" o algo similar

	json.NewEncoder(w).Encode(newTask) // le mando al cliente la tarea creada
}

func selectTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // mux.Vars(<request>) retorna las variables existentes en la URL (en este caso solo la variable "id" VER CUANDO DEFINO LA RUTA ABAJO)

	taskId, err := strconv.Atoi(vars["id"]) // strconv.Atoi() recibe un string y lo convierte a un int ademas de devolver un error si es necesario.

	if err != nil { // si el err NO ES NULO (hay error)
		fmt.Fprintf(w, "Invalid ID") // imprimo: Invalid ID
		return
	}

	for _, task := range tasks { // recorro todas las tasks actuales
		if task.Id == taskId { // si el task.id (de la task actual) es igual al taskId que me manda el cliente (POR LA URL VER ARRIBA)...
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "INVALID ID")
		return
	}

	for index, task := range tasks {
		if task.Id == taskId {
			tasks = append(tasks[:index], tasks[index+1:]...) // copia los elementos que estan antes del indice con tasks[:index] y los que estan despues del indice con tasks[index+1:]... esto hace eliminar el task corresp al index
			fmt.Fprintf(w, "Task with id %v removed", taskId)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "INVALID ID")
		return
	}

	reqBody, err := io.ReadAll(r.Body) // leo el body de la request
	if err != nil {
		fmt.Fprintf(w, "Invalid Data")
		return
	}

	json.Unmarshal(reqBody, &updatedTask) // almaceno el body de reqBody en la var "updatedTask"

	for i, task := range tasks {
		if task.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...) // elimino la task
			updatedTask.Id = taskId                   // vuelvo a configurar el id
			tasks = append(tasks, updatedTask)        // la agrego nuevamente

			fmt.Fprintf(w, "Task with ID %v updated!", taskId)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true) // creo un enrutador. Con .StrictSlash(true) hago que la ruta sea estrica. Es decir, si la url es /tasks tiene que ser SI O SI así sino genera error
	router.HandleFunc("/", indexRoute)          // cuando se visite la url "/",

	// <----- RUTAS ----->
	router.HandleFunc("/tasks", showTasks).Methods("GET")        // cuando llega a la url /tasks llama a la func showTasks que devuelve las tareas
	router.HandleFunc("/createTask", createTask).Methods("POST") // ruta que solo funciona con el metodo POST
	router.HandleFunc("/select/{id}", selectTask).Methods("GET")
	router.HandleFunc("/del/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/update/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router)) // escucha el sv en puerto 3000 y si hay algun error ejecuta el log.Fatal() | Fatal es equivalente a un Print() y termina el programa
}
