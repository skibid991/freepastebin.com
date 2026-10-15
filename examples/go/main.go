package main

import (
	"context"
	"fmt"
	clientpkg "github.com/skibid991/freepastebin.com/sdk/go"
	"log"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	ctx := context.Background()
	client := clientpkg.New()
	r, err := client.Create(ctx, "Hello from Go", map[string]any{"expires_in": 300}, "POST")
	must(err)
	var paste struct {
		ID    string `json:"id"`
		Token string `json:"delete_token"`
	}
	must(r.JSON(&paste))
	defer func() { must(client.Delete(ctx, paste.ID, paste.Token)) }()
	r, err = client.Raw(ctx, paste.ID, "")
	must(err)
	fmt.Println(r.Text())
}
