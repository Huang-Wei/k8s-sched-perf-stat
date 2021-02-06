package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{
		Use:   "schedperftest oldfile [newfile]",
		Short: "Compare results of scheduler integration performance tests",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("requires at least one argument")
			}
			if len(args) > 2 {
				return errors.New("at most two arguments can be provided")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			t := &Table{
				Keys: make(map[string]*Row),
				Mode: Old,
			}
			t.LoadFile(args[0], Old)
			if len(args) == 2 {
				t.Mode = New
				t.LoadFile(args[1], New)
			}
			FormatText(t)
		},
	}

	c.Execute()
}
