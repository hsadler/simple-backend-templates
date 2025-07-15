package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"example-server/internal/openapi/ogen"

	"github.com/fatih/color"
)

func run(ctx context.Context) error {
	client, err := ogen.NewClient("http://localhost:8000")
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	if err := testPing(ctx, client); err != nil {
		return err
	}
	return nil
}

func testPing(ctx context.Context, client *ogen.Client) error {
	resp, err := client.PingGet(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping: %v", err)
	}
	fmt.Println(resp)
	return nil
}

// func testCreateItem(ctx context.Context, client *ogen.Client) error {
// 	req := &ogen.CreateItemRequest{
// 		Data: ogen.ItemIn{
// 			Name:  "Test Item",
// 			Price: 19.99,
// 		},
// 	}
// 	resp, err := client.ItemsPost(ctx, req)
// 	if err != nil {
// 		return fmt.Errorf("failed to create item: %v", err)
// 	}
// 	fmt.Println(resp)
// 	return nil
// }

func main() {
	flag.Parse()
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		color.New(color.FgRed).Println(err)
		os.Exit(2)
	}
}
