package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/HushaamShah/gator/internal/config"
	"github.com/HushaamShah/gator/internal/database"
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
	commands.register("addfeed", middlewareLoggedIn(handleAddFeed))
	commands.register("feeds", handleFeeds)
	commands.register("follow", middlewareLoggedIn(handleFollow))
	commands.register("following", middlewareLoggedIn(handleFollowing))
	commands.register("unfollow", middlewareLoggedIn(handleUnfollow))

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

func handlerReset(s *state, cmd command) error {
	err := s.queries.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("All User Deleted!")
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		userDetails, err := s.queries.GetUser(context.Background(), s.dbconfig.User)
		if err != nil {
			return err
		}
		return handler(s, c, userDetails)
	}
}
