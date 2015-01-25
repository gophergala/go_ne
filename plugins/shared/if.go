package shared

import "fmt"

// If wraps the current command in an if statement with the given
// condition.
func (c *Command) If(condition string) {
	c.Name = fmt.Sprintf("if [ %v ]; then %v; fi", condition, c.Name)
}

// Unless wraps the current command in an if not statement with the given
// condition.
func (c *Command) Unless(condition string) {
	c.Name = fmt.Sprintf("if [ ! %v ]; then %v; fi", condition, c.Name)
}
