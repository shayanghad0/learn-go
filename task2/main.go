package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
)

type User struct {
    Name   string   `json:"name"`
    Age    int      `json:"age"`
    Skills []string `json:"skills"`
    Active bool     `json:"active"`
}

func main() {
    file, err := os.Open("data.json")
    if err != nil {
        fmt.Println("❌ Error opening file:", err)
        return
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println("❌ Error reading file:", err)
        return
    }

    var user User
    err = json.Unmarshal(data, &user)
    if err != nil {
        fmt.Println("❌ Error parsing JSON:", err)
        return
    }

    fmt.Println("✅ Parsed JSON Data:")
    fmt.Println("Name:", user.Name)
    fmt.Println("Age:", user.Age)
    fmt.Println("Skills:", user.Skills)
    fmt.Println("Active:", user.Active)
}
