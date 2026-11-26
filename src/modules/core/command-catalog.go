package core

import (
	"fmt"
	"sort"
)

type commandCatalog struct {
	byName map[string]Command
}

// NewCommandCatalog creates an empty catalog of commands.
func NewCommandCatalog() CommandCatalog {
	return &commandCatalog{byName: make(map[string]Command)}
}

func (c *commandCatalog) Add(command Command) error {
	if command.Name == "" {
		return ErrCommandNameEmpty
	}
	if _, exists := c.byName[command.Name]; exists {
		return fmt.Errorf("%w: %s", ErrCommandAlreadyRegistered, command.Name)
	}
	c.byName[command.Name] = command
	return nil
}

func (c *commandCatalog) List() []Command {
	names := make([]string, 0, len(c.byName))
	for name := range c.byName {
		names = append(names, name)
	}
	sort.Strings(names)
	commands := make([]Command, 0, len(names))
	for _, name := range names {
		commands = append(commands, c.byName[name])
	}
	return commands
}
