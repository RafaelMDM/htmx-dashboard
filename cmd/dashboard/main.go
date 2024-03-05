package main

import (
	"context"
	"os"
)

func main() {
	component := hello("Rafael")
	component.Render(context.Background(), os.Stdout)
}
