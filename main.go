package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var actions = []string{
	"log in",
	"log out",
	"view",
	"create",
	"update",
	"delete",
	"like",
	"comment",
}

// var r rand.Source

type LogItem struct {
	action string
	time   time.Time
}

type User struct {
	id    int
	email string
	logs  []LogItem
}

func NewUser(id int, email string, logs []LogItem) *User {
	return &User{id, email, logs}
}

func (u *User) getActivityInfo() string {
	activityInfo := fmt.Sprintf("ID: %d | Email: %s\nActivity log:\n", u.id, u.email)
	for i, log := range u.logs {
		activityInfo += fmt.Sprintf("%d. [%s] at %s\n", i+1, log.action, log.time.Format(time.RFC3339))
	}

	return activityInfo
}

func generateUser(count int) []User {
	users := make([]User, count)

	//replace with slice of pointers and try to "for i, user := range"
	for i := 0; i < count; i++ {
		users[i] = User{
			id:    i + 1,
			email: fmt.Sprintf("user_%d@gmail.com", i+1),
			logs:  generateLogs(500 + rand.Intn(1000)),
		}
	}

	return users
}

func generateLogs(count int) []LogItem {
	logs := make([]LogItem, count)

	for i := 0; i < count; i++ {
		logs[i] = LogItem{
			action: actions[rand.Intn(len(actions)-1)],
			time:   time.Now(),
		}
	}

	return logs
}

func writeUserLogs(user User, wg *sync.WaitGroup) error {
	time.Sleep(time.Millisecond * 10)
	fmt.Printf("Writing logs for user: %d\n", user.id)

	filename := fmt.Sprintf("logs/user_%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	_, err = file.WriteString(user.getActivityInfo())
	if err != nil {
		return err
	}

	wg.Done()

	return nil
}

// func init() {
// 	r = rand.NewSource(time.Now().Unix())
// }

func main() {
	rand.Seed(time.Now().Unix())

	t := time.Now()

	wg := &sync.WaitGroup{}

	users := generateUser(1000)

	for _, user := range users {
		wg.Add(1)
		go writeUserLogs(user, wg)
	}

	wg.Wait()

	fmt.Println("Execution time: %s", time.Since(t).String())
}
