package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // hashed
	Email    string `json:"email"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

const dbFile = "users.json"

var users = make(map[string]User)
var currentUser *User

func hashPassword(pw string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(pw)))
}

func saveUsers() {
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(dbFile, data, 0644)
}

func loadUsers() {
	data, err := os.ReadFile(dbFile)
	if err == nil {
		json.Unmarshal(data, &users)
	}
}

func register() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if _, exists := users[username]; exists {
		fmt.Println("❌ Username already exists.")
		return
	}

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Full Name: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("Age: ")
	var age int
	fmt.Scan(&age)

	users[username] = User{
		Username: username,
		Password: hashPassword(strings.TrimSpace(password)),
		Email:    strings.TrimSpace(email),
		Name:     strings.TrimSpace(name),
		Age:      age,
	}
	saveUsers()
	fmt.Println("✅ Registration complete.")
}

func login() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	hashed := hashPassword(password)

	user, exists := users[username]
	if !exists || user.Password != hashed {
		fmt.Println("❌ Invalid username or password.")
		return
	}

	currentUser = &user
	fmt.Println("✅ Logged in as", currentUser.Username)
	dashboard()
}

func dashboard() {
	if currentUser == nil {
		fmt.Println("⚠️ No user logged in.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n--- DASHBOARD ---")
		fmt.Printf("Username: %s\nName: %s\nEmail: %s\nAge: %d\n", currentUser.Username, currentUser.Name, currentUser.Email, currentUser.Age)
		fmt.Println("1. Edit Username")
		fmt.Println("2. Edit Password")
		fmt.Println("3. Edit Email")
		fmt.Println("4. Edit Name")
		fmt.Println("5. Edit Age")
		fmt.Println("6. Logout")
		fmt.Print("Select option: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Print("New Username: ")
			newUsername, _ := reader.ReadString('\n')
			newUsername = strings.TrimSpace(newUsername)
			delete(users, currentUser.Username)
			currentUser.Username = newUsername
			users[newUsername] = *currentUser
		case 2:
			fmt.Print("New Password: ")
			newPass, _ := reader.ReadString('\n')
			currentUser.Password = hashPassword(strings.TrimSpace(newPass))
		case 3:
			fmt.Print("New Email: ")
			newEmail, _ := reader.ReadString('\n')
			currentUser.Email = strings.TrimSpace(newEmail)
		case 4:
			fmt.Print("New Name: ")
			newName, _ := reader.ReadString('\n')
			currentUser.Name = strings.TrimSpace(newName)
		case 5:
			fmt.Print("New Age: ")
			var newAge int
			fmt.Scan(&newAge)
			currentUser.Age = newAge
		case 6:
			currentUser = nil
			saveUsers()
			fmt.Println("🔒 Logged out.")
			return
		default:
			fmt.Println("Invalid option.")
		}
		saveUsers()
	}
}

func main() {
	loadUsers()
	for {
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Exit")
		fmt.Print("Choose: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			register()
		case 2:
			login()
		case 3:
			fmt.Println("👋 Bye.")
			return
		default:
			fmt.Println("❌ Invalid choice.")
		}
	}
}
