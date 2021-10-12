package main

import (
    "encoding/json"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "net/http"
)

type Post struct {
    Id uint `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
    Comments []Comment `json:"comments" gorm:"-" default:"[]"`
}

type Comment struct {
    Id uint `json:"id"`
    PostId uint `json:"post_id"`
    Text string `json:"text"`
}

func main() {
    db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/posts_ms"), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    db.AutoMigrate(Post{})

    app := fiber.New()
    app.Use(cors.New())

    app.Get("/api/posts", func(c *fiber.Ctx) error {
        var posts []Post
        db.Find(&posts)

        for i, post := range posts {
            response, err := http.Get(fmt.Sprintf("http://localhost:8001/api/posts/%d/comments", post.Id))
            if err != nil {
                return err
            }
            var comments []Comment
            json.NewDecoder(response.Body).Decode(&comments)
            posts[i].Comments = comments
        }

        return c.JSON(posts)
        //return c.SendString("Hello, World 👋!")
    })

    app.Post("/api/posts", func(ctx *fiber.Ctx) error {
        var post Post
        if err := ctx.BodyParser(&post); err != nil{
            return err
        }
        db.Create(&post)
        //ctx.BodyParser(&posts)
        //db.Find(&posts)
        return ctx.JSON(post)
    })
    app.Listen(":8000")
}
