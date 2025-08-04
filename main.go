package main

import (
	"log"
	"os"
	"strings"

	"context"

	"github.com/urfave/cli/v3"
)

/*
protogen hello.proto --out-file hello.ts
*/

func main() {
	cmd := &cli.Command{
		Name:  "protogen",
		Usage: "Generate TypeScript types from proto files",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "out-file",
				Usage: "Generated results will be written to this file. Defaults to types.ts",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			file := "./example.proto"
			if cmd.NArg() > 0 {
				file = cmd.Args().Get(0)
			}

			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			// generate types
			tokens := analyze(content)
			data, err := parse(tokens)
			if err != nil {
				return err
			}
			output := generate(data)

			// export types to out-file
			outFile := cmd.String("out-file")
			if len(strings.TrimSpace(outFile)) == 0 {
				outFile = "./types.ts"
			}

			err = os.WriteFile(outFile, []byte(output), 0644)
			if err != nil {
				return err
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
