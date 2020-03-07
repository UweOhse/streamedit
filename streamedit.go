package main

import(
	"io"
	"fmt"
	"os"
	"time"
	"math"
	"flag"
	"strconv"
	"bufio"
	"math/rand"
)

var program string = "streamedit"
var helpOption bool
var versionOption bool
var verboseOption bool
var seedOption int = -1 // this is for the self check.
func init() {
	flag.BoolVar(&helpOption, "help", false, "show usage");
	flag.BoolVar(&versionOption, "version", false, "show version");
	flag.BoolVar(&verboseOption, "verbose", false, "be verbose about changes");
	flag.IntVar(&seedOption, "seed", -1, "set the random number seed (-1: by the nanosecond)");
}

var outBuf *bufio.Writer
const (
	opPrint = iota
	opDelete
	opInsert
	opInsertRandom
	opSet
	opSetRandom
	opSetRandomForced
	opBitOr
	opBitAnd
	opBitXor
	opSearch
	opPaste
	opDoSomethingRandom
)
type opType int
func (op opType) String() string {
	return []string{"p","-","r","R","=","r","R","|","&","^","/","P","~"}[op]
}

type spec struct {
	end                  bool
        prefix               uint64
	randomPrefix         bool // operation shall work [0,prefix) times.
	jumpToPos            bool

	op		     opType

	theByte              byte
	searchSilent         bool

}

var specs []spec
var specidx int
var curspec *spec
var togo uint64
var inputPosition uint
var outputPosition uint
var lastDeleted byte

func verboseByte(c,d byte) {
	if verboseOption {
		fmt.Fprintf(os.Stderr, "%d/%d: x%x -> x%x\n",inputPosition, outputPosition, c, d)
	}
}
func verboseByteOp(c byte, s string) {
	if verboseOption {
		fmt.Fprintf(os.Stderr, "%d/%d: x%x -> %s\n",inputPosition, outputPosition, c, s)
	}
}
func verboseOp(s string) {
	if verboseOption {
		fmt.Fprintf(os.Stderr, "%d/%d: %s\n",inputPosition, outputPosition, s)
	}
}
func logicError(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr,"internal error: "+s+"\n\n",args...)
	fmt.Fprintf(os.Stderr,"Please report this to %s or %s.",
		PackageIssues, PackageBugreport)
	os.Exit(1)
}
func simpleOps(c byte) byte {
	if curspec.op==opSetRandom {
		t := byte(rand.Intn(256))
		verboseByte(c, t)
		return t
	}
	if curspec.op==opSetRandomForced {
		for {
			d:=byte(rand.Intn(256))
			if c!=d {
				verboseByte(c, d)
				return d
			}
		}
	}
	if curspec.op==opBitOr {
		t := c|curspec.theByte
		verboseByte(c, t)
		return t
	}
	if curspec.op==opBitAnd {
		t := c&curspec.theByte
		verboseByte(c, t)
		return t
	}
	if curspec.op==opBitXor {
		t := c ^ curspec.theByte
		verboseByte(c, t)
		return t
	}
	if curspec.op==opSet {
		t := curspec.theByte
		verboseByte(c, t)
		return t
	}
	logicError("invalid action specification in inner loop: %#v",curspec)
	return byte(0) // dummy
}

func writeOut(inbuf []byte) {
	n, err := outBuf.Write(inbuf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write: %v\n",err)
		os.Exit(1)
	}
	outputPosition+=uint(n)
	return
}

func nextAction() {
	var wasSearch bool
	var wasByte byte
	if curspec!=nil { // special case for the first use.
		specidx++
		if curspec.op==opSearch {
			wasSearch=true
			wasByte=curspec.theByte
		}
	}
	if specidx>=len(specs) {
		specidx=0
	}
	curspec=&specs[specidx]
	togo=curspec.prefix
	if curspec.randomPrefix {
		togo=uint64(rand.Intn(int(togo)))
	}
	if curspec.op==opSearch && wasSearch && wasByte==curspec.theByte {
		fmt.Fprintf(os.Stderr, "error: invalid search. two searches directly following each other must not search for the same byte.\n")
		usageError()
	}
}

func editorLoop(inbuf []byte) {
	for {
// _ = outBuf.Flush()
// fmt.Printf("buf %v, op=%s, spec %#v, togo=%d\n",string(inbuf), curspec.op.String(), curspec, togo)
		// n:=uint64(len(inbuf))
		if togo==0 && !curspec.jumpToPos {
			nextAction()
			continue
		}

		if len(inbuf)==0 {
			// fmt.Fprintf(os.Stderr, "l=0, cs=%#v\n",curspec)
			return
		}

		if curspec.end && curspec.op==opPrint { // * p needs to be fast.
			writeOut(inbuf)
			return
		}
		if curspec.end && curspec.op==opDelete { // * - needs to be fast, too.
			verboseOp("delete all")
			return
		}

		if curspec.op==opSearch {
//			fmt.Printf("s %v against %v, togo=%d\n",curspec.theByte, inbuf[0], togo)
			if curspec.theByte==inbuf[0] {
				togo--
			}
			if curspec.theByte!=inbuf[0] || togo>0 {
				if !curspec.searchSilent {
					writeOut(inbuf[:1])
				} else {
					verboseByteOp(inbuf[0], "deleted")
				}
				inbuf=inbuf[1:]
				inputPosition++
			}
			continue
//			fmt.Printf(" ==> %v with togo=%d\n",string(inbuf),togo)
		}
// fmt.Fprintf(os.Stderr,"i=%d jtp=%v c.p=%d\n",inputPosition,curspec.jumpToPos, curspec.prefix)
		if !curspec.jumpToPos {
			togo--
		} else {
			if uint64(inputPosition) > curspec.prefix {
				nextAction()
				continue
			}
		}

		if curspec.op==opPaste {
			writeOut([]byte{lastDeleted})
			verboseByteOp(lastDeleted, "inserted (last deleted)")
			continue
		}

		if curspec.op==opDelete {
			verboseByteOp(inbuf[0], "deleted")
			lastDeleted=inbuf[0]
			inbuf=inbuf[1:]
			inputPosition++
			continue
		}
		if curspec.op==opPrint {
			writeOut(inbuf[:1])
			inbuf=inbuf[1:]
			inputPosition++
			continue
		}

		if curspec.op==opInsert {
			verboseByteOp(curspec.theByte, "inserted")
			writeOut([]byte{curspec.theByte})
			continue
		}
		if curspec.op==opInsertRandom {
			t:=byte(rand.Intn(256))
			verboseByteOp(t, "inserted")
			writeOut([]byte{t})
			continue
		}
		if curspec.op==opDoSomethingRandom {
			t:=rand.Intn(5)
			if 0==t {
				verboseByteOp(inbuf[0], "deleted (~ random)")
				lastDeleted=inbuf[0]
				inbuf=inbuf[1:]
				inputPosition++
			} else if 1==t {
				t:=byte(rand.Intn(256))
				verboseByteOp(t, "inserted (~ random)")
				writeOut([]byte{t})
			} else if 2==t {
				t:=inbuf[0] ^ 255;
				verboseByteOp(t, "flipped (~ random)")
				writeOut([]byte{t})
				inbuf=inbuf[1:]
				inputPosition++
			} else if 3==t {
				t:=byte(rand.Intn(256))
				verboseByteOp(t, "set (~ random)")
				writeOut([]byte{t})
				inbuf=inbuf[1:]
				inputPosition++
			} else if 4==t {
				t:=byte(1<<rand.Intn(8))
				t=inbuf[0] ^ t
				verboseByteOp(t, "bitflipped (~ random)")
				writeOut([]byte{t})
				inbuf=inbuf[1:]
				inputPosition++
			}
			continue
		}

		c:=simpleOps(inbuf[0])
		writeOut([]byte{c})
		inbuf=inbuf[1:]
		inputPosition++
	}
}

func parseByte(s, orig string) byte {
	var t uint64
	var err error
	if len(s)==0 {
		fmt.Fprintf(os.Stderr, "bad specification %s: empty byte specification\n",orig)
		usageError()
	} else if s[0]=='o' {
		t, err=strconv.ParseUint(s[1:],8,8)
	} else if (s[0]=='x') {
		t, err=strconv.ParseUint(s[1:],16,8)
	} else if (s[0]=='c') {
		if len(s)!=2 {
			fmt.Fprintf(os.Stderr, "bad specification %s: not exactly one character after the c.\n",orig)
			usageError()
		}
		t=uint64(s[1])
	} else {
		t, err=strconv.ParseUint(s,10,8)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad specification %s: %v\n",orig, err)
		usageError()
	}
	return byte(t)
}
func parseNumber(s, orig string) uint64 {
	var t uint64
	var err error
	if (s[0]=='x') {
		t, err=strconv.ParseUint(s[1:],16,64)
	} else if (s[0]=='o') {
		t, err=strconv.ParseUint(s[1:],8,64)
	} else if (s[0]=='c') {
		if len(s)!=2 {
			fmt.Fprintf(os.Stderr, "bad specification %s: not exactly one character after the c.\n",orig)
			usageError()
		}
		t=uint64(s[1])
	} else {
		t, err=strconv.ParseUint(s,10,64)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad specification %s: %v\n",orig, err)
		usageError()
	}
	return t
}

func catchEnd(s string, haveEnd bool) {
	if !haveEnd {
		return
	}
	fmt.Fprintf(os.Stderr, "bad specification: %s follows an end action.\n",s)
	usageError()
}

/*
  bON : bit operation:
        b^1 switch bit 0.
        b|192 set bit 6 and 7 to 1 ("64" + "128").
        b&4 set bit 2 to 0 ("4").
*/
func parseSpecs(args []string) {
	specs=make([]spec,0)
	haveRandom := false
	haveJump := false
	prefix := uint64(1)
	haveEnd := false
	haveWildcard := false
	for _, a := range args {
		var S spec
		S.prefix=prefix
		if haveRandom {
			S.randomPrefix=true
		}
		if haveJump {
			S.jumpToPos=true
		}
		if haveWildcard {
			S.prefix=math.MaxUint64
			S.end=true
		}
		prefix=1
		haveRandom=false
		haveJump=false
		haveWildcard=false
		if len(a) == 0 {
			fmt.Fprint(os.Stderr, "bad specification: empty string.\n",a)
			usageError()
		}

		if a == "end" {
			S.end=true
			S.op=opPrint
			specs=append(specs,S)
			catchEnd(a, haveEnd)
			haveEnd = true
			continue
		}
		if a == "all" {
			haveWildcard=true
			continue
		}
		if a == "*" {
			haveWildcard=true
			continue
		}
		if a == "p" {
			catchEnd(a, haveEnd)
			S.op=opPrint
			specs=append(specs,S)
			continue
		}
		if a == "~" {
			catchEnd(a, haveEnd)
			S.op=opDoSomethingRandom
			specs=append(specs,S)
			continue
		}
		if a == "P" {
			S.op=opPaste
			specs=append(specs,S)
			continue
		}
		if a == "r" {
			catchEnd(a, haveEnd)
			S.op=opSetRandom
			specs=append(specs,S)
			continue
		}
		if a == "R" {
			catchEnd(a, haveEnd)
			S.op=opSetRandomForced
			specs=append(specs,S)
			continue
		}
		if a == "-" {
			catchEnd(a, haveEnd)
			S.op=opDelete
			specs=append(specs,S)
			continue
		}
		if a == "+r" {
			S.op=opInsertRandom
			specs=append(specs,S)
			continue
		}
		if a[0] == '+' {
			S.theByte=parseByte(a[1:], a)
			S.op=opInsert
			specs=append(specs,S)
			continue
		}
		if a[0] == '=' {
			catchEnd(a, haveEnd)
			S.theByte=parseByte(a[1:], a)
			S.op=opSet
			specs=append(specs,S)
			continue
		}
		if a[0]=='|' || a[0]=='O' || a[0]=='&' || a[0]=='A' || a[0]=='^' {
			catchEnd(a, haveEnd)
			op:=a[0]
			S.theByte=parseByte(a[1:], a)
			if op=='|' || op=='O' {
				S.op=opBitOr
				specs=append(specs,S)
				continue;
			}
			if op=='&' || op=='A' {
				S.op=opBitAnd
				specs=append(specs,S)
				continue;
			}
			if op=='^' {
				S.op=opBitXor
				specs=append(specs,S)
				continue;
			}
			fmt.Fprintf(os.Stderr, "bad specification %s: invalid operation.\n",a)
			usageError()
		}
		if a[0]=='q' || a[0]=='v'  {
			catchEnd(a, haveEnd)
			if len(a)<2 {
				fmt.Fprintf(os.Stderr, "bad specification %s: expected / and number.\n",a)
				usageError()
			}
			S.searchSilent = a[0]=='q'
			S.theByte=parseByte(a[1:],a)
			S.op=opSearch
			specs=append(specs,S)
			continue;
		}
		l:=len(a)
		if l>1 && a[l-1]=='r' {
			prefix=parseNumber(a[:l-1],a)
			haveRandom=true
			if prefix < 1 { // 0 panic()s, 1 is just useless.
				fmt.Fprintf(os.Stderr, "bad specification %s: invalid prefix.\n",a)
				usageError()
			}
			continue
		}
		if l>1 && a[l-1]=='@' {
			catchEnd(a, haveEnd)
			prefix=parseNumber(a[:l-1],a)
			haveJump=true
			continue
		}

		// NN - a prefix bytes.
		prefix=parseNumber(a,a)
		continue;
	}
	if len(specs)==0 {
		usageError()
	}
	nextAction()

}
func usageError() {
	fmt.Fprintf(os.Stderr, "Try streamedit --help for more information.\n")
	os.Exit(2)
}
func showHelp() {
	fmt.Print(
`Usage: streamedit [options] actions ...
A byte-stream editor, made to modify bytes, not characters.

Options:
`)
	flag.CommandLine.SetOutput(os.Stdout)
	flag.PrintDefaults()

	fmt.Printf(`
actions... is a series of actions, processed in order. If the program
reaches the end of the actions, it will continue with the first one.

An action is either:
  PFX : a prefix for the next action, defaulting to 1. 
  *   : a special prefix meaning "all following bytes". Alias: all.
        This especially avoids the looping at the end of the action list.
  -   : deletes PFX bytes.
  p   : prints PFX bytes.
  r   : changes PFX bytes to a random value (which may be the original one)
  R   : changes PFX bytes to a random value (different from the original one)
  =BY : changes PFX bytes to BY.
  +BY : inserts PFX times BY at the current position.
  +r  : inserts PFX times a random byte at the current position.
  |BY : "or"s the current byte with BY.
  oBY : "or"s the current byte with BY (an alias for ease of use).
  &BY : "and"s the current byte with BY.
  aBY : "and"s the current byte with BY.
  ^BY : "xor"s the current byte with BY.
  vBY : search for the PFXth occurance of BY and print all bytes before it
        ("v" as in "verbose search").
  qBY : search for the PFXth occurance of BY without printing the bytes
        before it. ("q" as in "quiet search").
  ~   : does something random to PFX input bytes. There is an equal change
        that a byte set to a random value, XORed with 255, deleted, that 
	a single bit is flipped, or a random byte is inserted.
  end : is the same as '* p' (print all until reaching the end of the input).

BY is a bytes value. By default it's an decimal number ("64" stands for '@'),
but you may use 'x' and 'o' as prefixes for hexadecimal and octal numbers, 
and you may use 'c' as prefix for a character used literally. "64", "c@",
"x40" and "o100" all mean the very same thing, an '@' character.

PFX is a number. 1023, x3ff and o1777 all mean the same thing. The prefix
determines how often the next action will be executed.
If you append a @ sign to PFX, then PFX determines up to which position in
the input stream the next action will be executed (counted from 0). If that
position has already been passed, the action will not be executed.
If you append an r to PFX, the meaning changes to 'a random number of bytes
below N' (0 to N-1 bytes).

Report bugs to either your distributor, <%s>,
or <%s>.
Homepage: <%s>
`, PackageBugreport, PackageIssues, PackageURL)
}
func showVersion() {
	fmt.Printf("%s (%s) %s [git %s]\n",program,PackageName,PackageVersion,Build)
	fmt.Printf(
`(C) 2020+ Uwe Ohse, <uwe@ohse.de>.
License: GNU GPL version 2.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
`)
}

func endHandling() {
	if !curspec.end {
		return
	}
	for {
		specidx++
		if specidx>=len(specs) {
			return
		}
		curspec=&specs[specidx]
		togo=curspec.prefix
		if curspec.randomPrefix {
			togo=uint64(rand.Intn(int(togo)))
		}
		for {
			if togo==0 {
				break
			}
			togo--
			if curspec.op==opInsert {
				writeOut([]byte{curspec.theByte})
				continue
			}
			if curspec.op==opInsertRandom {
				t:=byte(rand.Intn(256))
				writeOut([]byte{t})
				continue
			}
			if curspec.op==opPaste {
				writeOut([]byte{lastDeleted})
				verboseByteOp(lastDeleted, "inserted (last deleted)")
				continue
			}
			logicError("invalid action in end loop: %#v",curspec)
		}

	}
}
func mainLoop() {
	bufsize:=65536
	buf:=make([]byte,bufsize,bufsize);

	for {
		buf=buf[:bufsize]
		n, err := os.Stdin.Read(buf)
		buf=buf[:n]
		if n!=0 || true {
			editorLoop(buf)

		}
		// fmt.Fprintf(os.Stderr, "read returned %d %v\n", n, err);
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr,"failed to read: %s\n", err);
				os.Exit(1);
			}
			break;
		}
	}
	endHandling()
}
func main() {
	flag.Parse()
	if versionOption {
		showVersion()
		return
	}
	if helpOption {
		showHelp()
		return
	}
	outBuf = bufio.NewWriter(os.Stdout)


	if seedOption<0 {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(int64(seedOption))
	}

	parseSpecs(flag.Args())

	mainLoop()

	err := outBuf.Flush() // Don't forget to flush!
	if err != nil {
		fmt.Fprintf(os.Stderr,"failed to write: %s\n", err);
		os.Exit(1);
	}
}
