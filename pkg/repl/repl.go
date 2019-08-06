package repl

import (
	"fmt"
	"github.com/chzyer/readline"
	"go.smartmachine.io/cumulus/pkg/compiler"
	"go.smartmachine.io/cumulus/pkg/lexer"
	"go.smartmachine.io/cumulus/pkg/object"
	"go.smartmachine.io/cumulus/pkg/parser"
	"go.smartmachine.io/cumulus/pkg/vm"
	"io"
	"io/ioutil"
)

const PROMPT = ">> "
const RESULT = "== "

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

		comp := compiler.New()
		err = comp.Compile(program)
		if err != nil {
			write(out, "Whoops! Compilation failed:\n %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			write(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		write(out, "%s%s\n", RESULT, lastPopped.String())
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

func write(out io.Writer, message string, a ...interface{}) {
	num, err := fmt.Fprintf(out, message, a...)
	if num != len(message) {
		fmt.Printf("unable to write all %d bytes to output channel (%d written)\n", len(message), num)
	}
	if err != nil {
		fmt.Printf("error writing to output channel: %v", err)
	}
}