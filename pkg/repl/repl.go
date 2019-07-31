package repl

import (
	"fmt"
	"github.com/chzyer/readline"
	"go.smartmachine.io/cumulus/pkg/evaluator"
	"go.smartmachine.io/cumulus/pkg/lexer"
	"go.smartmachine.io/cumulus/pkg/object"
	"go.smartmachine.io/cumulus/pkg/parser"
	"io"
	"io/ioutil"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {

	rl, err := readline.NewEx(&readline.Config{
		Stdout: out,
		Stdin: ioutil.NopCloser(in),
		Prompt: PROMPT,
		StdinWriter: out,
		Stderr: out,
		HistoryFile: ".cumulus.hist",
		HistoryLimit: 1000,
	})

	if err != nil {
		fmt.Printf("unable to start readline: %v", err)
		return
	}
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		line, err := rl.Readline()
		if err != nil {
			fmt.Printf("unable to read line: %v", err)
			return
		}
		if line == "" {
			continue
		}

		switch line {
		case `\q`:
			write(out, "Bye!\n")
			return
		case `\env`:
			write(out, fmt.Sprintf("Environment: %s\n", env))
			continue
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			write(out, evaluated.String())
			write(out, "\n")
		}
	}
}

const SHROOM = `        --_--
     (  -_    _).
   ( ~       )   )
 (( )  (    )  ()  )
  (.   )) (       )` +
"\n    ``..     ..``\n" +
`         | |
       (=| |=)
         | |
     (../( )\.))
`

func printParserErrors(out io.Writer, errors []string) {
	write(out, SHROOM)
	write(out, "Woops! Mushroom cloud!\n")
	write(out, " parser errors:\n")
	for _, msg := range errors {
		write(out, "\t"+msg+"\n")
	}
}

func write(out io.Writer, message string) {
	num, err := io.WriteString(out, message)
	if num != len(message) {
		fmt.Printf("unable to write all %d bytes to output channel (%d written)\n", len(message), num)
	}
	if err != nil {
		fmt.Printf("error writing to output channel: %v", err)
	}
}