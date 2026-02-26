package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/HushaamShah/gator/internal/config"
	"github.com/HushaamShah/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	dbconfig *config.Config
	queries  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *commands) run(s *state, cmd command) error {
	if val, ok := c.commands[cmd.name]; ok {
		return val(s, cmd)
	} else {
		return fmt.Errorf("Command Not Found")
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func main() {
	commands := commands{commands: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handleAggregate)

	dbconfig, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	var st state
	st.dbconfig = &dbconfig
	db, err := sql.Open("postgres", st.dbconfig.DBURL)
	st.queries = database.New(db)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("too few args")
		os.Exit(1)
	}
	var cmd command
	cmd.name = args[1]
	cmd.args = args[2:]

	err = commands.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login command requires a username")
	}
	_, err := s.queries.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User not Found")
		}
		return err
	}

	s.dbconfig.SetUser(cmd.args[0])
	fmt.Println("User has been set.")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	var args database.CreateUserParams
	args.ID = uuid.New()
	args.CreatedAt = time.Now()
	args.UpdatedAt = time.Now()
	args.Name = cmd.args[0]

	user, err := s.queries.CreateUser(context.Background(), args)
	if err != nil {
		return err
	}
	fmt.Println("USER: ")
	fmt.Println(user)
	s.dbconfig.SetUser(cmd.args[0])
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.queries.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("All User Deleted!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.queries.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.dbconfig.User {
			fmt.Printf("* %s (current) \n", user)
			continue
		}
		fmt.Printf("* %s \n", user)
	}
	return nil
}

func handleAggregate(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	xmlBody := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("user-agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err1 := xml.Unmarshal(body, &xmlBody)
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}

	xmlBody.Channel.Title = html.UnescapeString(xmlBody.Channel.Title)
	xmlBody.Channel.Description = html.UnescapeString(xmlBody.Channel.Description)

	for i := range xmlBody.Channel.Item {
		xmlBody.Channel.Item[i].Title = html.UnescapeString(xmlBody.Channel.Item[i].Title)
		xmlBody.Channel.Item[i].Description = html.UnescapeString(xmlBody.Channel.Item[i].Description)
	}

	return &xmlBody, nil
}
