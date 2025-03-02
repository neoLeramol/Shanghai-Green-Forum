package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

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

    return db, nil
}

// 插入测试帖子
func insertTestPosts() {
    db, err := initDB()
    if err != nil {
        log.Fatalf("数据库初始化失败: %v", err)
    }
    defer db.Close()

    // 三段自然描写选段
    texts := []string{
        "清晨，第一缕阳光轻柔地洒在翠绿的草地上，露珠在草叶上闪烁着晶莹的光芒，宛如珍珠般璀璨。微风拂过，带来了泥土的芬芳和花朵的清香，让人仿佛置身于一个梦幻的世界。",
        "午后，湛蓝的天空中飘着几朵洁白的云朵，像棉花糖一样轻盈。远处的山峦连绵起伏，与蓝天白云相映成趣。一条清澈的小溪潺潺流淌，溪水在阳光的照耀下波光粼粼，仿佛是一条流动的丝带。",
        "傍晚，夕阳的余晖将整个天空染成了橙红色，绚丽多彩。森林里，鸟儿归巢，它们欢快的歌声在林间回荡。落叶在地上堆积，踩上去发出沙沙的声响，仿佛是大自然演奏的美妙乐章。",
    }

    location := "31.201556,121.46126"
    image := "placeholder.jpg"
    date := time.Now()

    for _, text := range texts {
        result, err := db.Exec("INSERT INTO posts (text, image, date, location) VALUES (?,?,?,?)", text, image, date, location)
        if err != nil {
            log.Printf("插入帖子失败: %v", err)
            continue
        }

        id, err := result.LastInsertId()
        if err != nil {
            log.Printf("获取插入帖子的 ID 失败: %v", err)
            continue
        }

        fmt.Printf("成功插入帖子，ID: %d\n", id)
    }
}

func main() {
    insertTestPosts()
}