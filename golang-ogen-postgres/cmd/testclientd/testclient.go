package main

import (
	"context"
	"fmt"
	"os"

	"example-server/internal/openapi/ogen"

	"github.com/fatih/color"
)

func randId() int64 {
	return int64(os.Getpid()) + int64(os.Getuid()) + int64(os.Geteuid())
}

func run(ctx context.Context) error {
	client, err := ogen.NewClient("http://localhost:8000")
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	if err := testPing(ctx, client); err != nil {
		return err
	}
	if err := testCreateItem(ctx, client); err != nil {
		return err
	}
	if err := testGetItem(ctx, client); err != nil {
		return err
	}
	if err := testUpdateItem(ctx, client); err != nil {
		return err
	}
	return nil
}

func testPing(ctx context.Context, client *ogen.Client) error {
	resp, err := client.Ping(ctx)
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp)
	return nil
}

func testCreateItem(ctx context.Context, client *ogen.Client) error {
	req := &ogen.ItemCreateRequest{
		Data: ogen.ItemIn{
			Name:  fmt.Sprintf("Item-%d", 1000+randId()),
			Price: 19.99,
		},
	}
	resp, err := client.CreateItem(ctx, req)
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp)
	return nil
}

func testGetItem(ctx context.Context, client *ogen.Client) error {
	resp, err := client.GetItem(ctx, ogen.GetItemParams{
		ItemId: 1,
	})
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp)
	return nil
}

func testUpdateItem(ctx context.Context, client *ogen.Client) error {
	req := &ogen.ItemUpdateRequest{
		Data: ogen.ItemIn{
			Name:  fmt.Sprintf("Updated Item-%d", 1000+randId()),
			Price: 29.99,
		},
	}
	resp, err := client.UpdateItem(
		ctx,
		req,
		ogen.UpdateItemParams{
			ItemId: 1,
		},
	)
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp)
	return nil
}

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		color.New(color.FgRed).Println(err)
		os.Exit(2)
	}
}
