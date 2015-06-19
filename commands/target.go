package commands

import (
	"fmt"
	"os"

	"github.com/backstage/backstage-cli/backstage"
	"github.com/codegangsta/cli"
)

type Target struct {
	backstage.Target
}

func (cmd *Target) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "target-add",
			Usage:       "target-add <label> <endpoint>",
			Description: "Add a new target in the list of targets.",
			Action:      cmd.targetAdd,
		},

		{
			Name:        "target-list",
			Usage:       "",
			Description: "List all targets.",
			Action:      cmd.targetList,
		},

		{
			Name:        "target-remove",
			Usage:       "target-remove <label>",
			Description: "Remove a target from the list of targets.",
			Action:      cmd.targetRemove,
		},

		{
			Name:        "target-set",
			Usage:       "target-set <label>",
			Description: "Set a target as default.",
			Action:      cmd.targetSet,
		},
	}
}

func (cmd *Target) targetAdd(c *cli.Context) {
	defer RecoverStrategy("target-add")()

	targets, err := backstage.LoadTargets()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	args := c.Args()
	label, endpoint := args[0], args[1]
	err = targets.Add(label, endpoint)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Your new target has been added.")
}

func (cmd *Target) targetList(c *cli.Context) {
	defer RecoverStrategy("target-list")()

	targets, err := backstage.LoadTargets()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	current, sortedMapKeys, options := targets.GetOptions()

	table := &Table{
		Content: [][]string{},
		Header:  []string{"Default", "Label", "Backstage Server"},
	}

	for _, label := range sortedMapKeys {
		endpoint := options[label]
		line := []string{""}

		if current == label {
			line[0] = "*"
		}
		line = append(line, label, endpoint)
		table.Content = append(table.Content, line)
	}

	context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
	table.Render(context)
}

func (cmd *Target) targetRemove(c *cli.Context) {
	defer RecoverStrategy("target-remove")()

	context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
	if Confirm(context, "Are you sure you want to remove this target? This action cannot be undone.") != true {
		fmt.Println(backstage.ErrCommandCancelled)
	} else {
		targets, err := backstage.LoadTargets()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		label := c.Args().First()
		err = targets.Remove(label)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("The target `" + label + "` has been remove.")
	}
}

func (cmd *Target) targetSet(c *cli.Context) {
	defer RecoverStrategy("target-set")()

	targets, err := backstage.LoadTargets()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	label := c.Args().First()
	err = targets.SetDefault(label)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("You have a new target as default!")
}
