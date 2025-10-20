package parser

import "fmt"

func (e ParseError) Error() string {
	return fmt.Sprintf("Error in table '%s', chain '%s' at line %d: %s (content: %+v)",
		e.Table, e.Chain, e.LineNumber, e.Message, e.LineText)
}
