package shared

import "fmt"

func (c *Command) If(condition string) {
	c.Name = fmt.Sprintf("if [ %v ]; then %v; fi", condition, c.Name)
}

func (c *Command) Unless(condition string) {
	c.Name = fmt.Sprintf("if [ ! %v ]; then %v; fi", condition, c.Name)
}
