package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TaskService struct {
	DB          *sql.DB
	TaskChannel chan Task
}

func (t *TaskService) AddTask(ts *Task) error {
	query := "ÏNSER INTO tasks (title,description, status, created_at) VALUE (?,?,?,?)"
	result, err := t.DB.Exec(query, ts.Title, ts.Description, ts.Description, ts.CreatedAt)
	if err != nil {
		return nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ts.ID = int(id)
	return err
}

func (t *TaskService) updateTaskStatus(ts Task) error {
	query := "UPDATE tasks SET status =? where id=?"
	_, err := t.DB.Exec(query, ts.Status, ts.ID)
	return err
}

func (t *TaskService) ListTask() ([]Task, error) {
	rows, err := t.DB.Query("Selec * from tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.CreatedAt,
			&task.Description,
			&task.Status,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil

}

func (t *TaskService) ProcessTasks() {
	for task := range t.TaskChannel {
		log.Printf("Processing task : %s", task.Title)
		time.Sleep(5 * time.Second)
		task.Status = "completed"
		t.updateTaskStatus(task)
		log.Println("Taks %s processed", task.Title)
	}
}

func (t *TaskService) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.Status = "pending"
	task.CreatedAt = time.Now()
	err = t.AddTask(&task)
	if err != nil {
		http.Error(w, "Error add task", http.StatusInternalServerError)
		return
	}
	t.TaskChannel <- task // mandando task para processament, jogando para um canal
	w.WriteHeader(http.StatusCreated)

}
func (t *TaskService) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.ListTask()
	if err != nil {
		http.Error(w, "Error ao listar tarefas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "Application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	TaskChannel := make(chan Task)

	taskService := TaskService{
		DB:          db,
		TaskChannel: TaskChannel,
	}
	go taskService.ProcessTasks()

	http.HandleFunc("POST /tasks", taskService.HandleCreateTask)
	http.HandleFunc("GET /tasks", taskService.HandleListTasks)
	http.ListenAndServe(":8080", nil)

}