package main

import (
	"context"
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
	if err := testCreateItem(ctx, client); err != nil {
		return err
	}
	if err := testItemsGet(ctx, client); err != nil {
		return err
	}
	if err := testItemsIDGet(ctx, client); err != nil {
		return err
	}
	if err := testItemsAllGet(ctx, client); err != nil {
		return err
	}
	return nil
}

func testPing(ctx context.Context, client *ogen.Client) error {
	resp, err := client.PingGet(ctx)
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp)
	return nil
}

func testCreateItem(ctx context.Context, client *ogen.Client) error {
	req := &ogen.CreateItemRequest{
		Data: ogen.ItemIn{
			Name: fmt.Sprintf(
				"Item-%d",
				1000+int64(os.Getpid())+int64(os.Getuid())+int64(os.Geteuid()),
			),
			Price: 19.99,
		},
	}
	resp, err := client.ItemsPost(ctx, req)
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp.Data)
	return nil
}

func testItemsGet(ctx context.Context, client *ogen.Client) error {
	resp, err := client.ItemsGet(ctx, ogen.ItemsGetParams{
		ItemIds: []int{1, 2, 3},
	})
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp.Data)
	return nil
}

func testItemsIDGet(ctx context.Context, client *ogen.Client) error {
	resp, err := client.ItemsIDGet(ctx, ogen.ItemsIDGetParams{
		ID: 2,
	})
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp.Data)
	return nil
}

func testItemsAllGet(ctx context.Context, client *ogen.Client) error {
	resp, err := client.ItemsAllGet(ctx, ogen.ItemsAllGetParams{
		ChunkSize: 10,
		Offset:    0,
	})
	if err != nil {
		color.New(color.FgRed).Println(err)
		return err
	}
	color.New(color.FgGreen).Println(resp.Data)
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
