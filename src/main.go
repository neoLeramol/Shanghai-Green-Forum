package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Post 定义帖子结构体
type Post struct {
	ID       int       `json:"id"`
	Text     string    `json:"text"`
	Image    string    `json:"image"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

// 初始化数据库
func initDB() (*sql.DB, error) {
	// 请根据实际情况修改数据库连接信息
	dsn := "user:password@tcp(127.0.0.1:3306)/your_database_name"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// 创建 posts 表
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id INT AUTO_INCREMENT PRIMARY KEY,
            text TEXT NOT NULL,
            image VARCHAR(255) NOT NULL,
            date DATETIME NOT NULL,
            location TEXT
        )
    `)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 获取所有帖子
func getPosts(w http.ResponseWriter, r *http.Request) {
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, text, image, date, location FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Text, &post.Image, &post.Date, &post.Location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// 搜索帖子
func searchPosts(w http.ResponseWriter, r *http.Request) {
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	idStr := r.URL.Query().Get("id")
	date := r.URL.Query().Get("date")

	var query string
	var args []interface{}
	if idStr != "" && date != "" {
		query = "SELECT id, text, image, date, location FROM posts WHERE id = ? AND date = ?"
		id, _ := strconv.Atoi(idStr)
		args = []interface{}{id, date}
	} else if idStr != "" {
		query = "SELECT id, text, image, date, location FROM posts WHERE id = ?"
		id, _ := strconv.Atoi(idStr)
		args = []interface{}{id}
	} else if date != "" {
		query = "SELECT id, text, image, date, location FROM posts WHERE date = ?"
		args = []interface{}{date}
	} else {
		query = "SELECT id, text, image, date, location FROM posts"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Text, &post.Image, &post.Date, &post.Location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// 发布帖子
func publishPost(w http.ResponseWriter, r *http.Request) {
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	text := r.FormValue("text")
	location := r.FormValue("location")
	date := time.Now()

	// 处理图片上传（这里只是简单示例，实际需要更完善的处理）
	image := "placeholder.jpg"

	result, err := db.Exec("INSERT INTO posts (text, image, date, location) VALUES (?,?,?,?)", text, image, date, location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "帖子发布成功", "id": strconv.FormatInt(id, 10)})
}

// 删除帖子
func deletePost(w http.ResponseWriter, r *http.Request) {
	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "帖子删除成功"})
}

func main() {
	http.HandleFunc("/posts", getPosts)
	http.HandleFunc("/search", searchPosts)
	http.HandleFunc("/publish", publishPost)
	http.HandleFunc("/delete", deletePost)

	log.Println("服务器启动，监听端口 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
