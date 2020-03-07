# streamedit - a byte stream editor

streamedit is a byte oriented stream editor, designed for manipulations of byte streams. It doesn't understand UTF8, character sets, or languages, and changes byte even if that breaks higher level encodings. In fact, breaking higher levels is its main purpose.
streamedit is a by-product of the lrzsz package, and has been designed to allow simulations of line errors.

## Usage

	streamedit [options] `actions` ...

## Options:

*  -help
    	show usage.

*  -seed int
    	set the random number seed (-1: by the current nanosecond, which is the default).
	This is used inside the test suite.

*  -verbose
    	be verbose about changes. It prints input position, output position, and what it does. This sometimes helps to understand what is going on. 

*  -version
    	show the version information.
## Actions

`actions ...` is a series of actions, processed in order. If the program reaches the end of the actions, it will continue with the first one, unless `end` or `*` has been encountered.

An action is either:

+ PFX
	The prefix for the next action (only), defaulting to 1.

	PFX is a number. 1023, x3ff and o1777 all mean the same thing. The prefix 
	determines how often the next action will be executed.

	If you append a `@` sign to PFX, then PFX determines up to which position in
	the input stream the next action will be executed (counted from 0). If that
	position has already been passed, the action will not be executed.

	If you append an literal `r` to PFX, the meaning changes to 'a random number of bytes
	below N' (0 to N-1 bytes). **Note**: This may change in the future to mean 0 to N bytes.

+ *
	A special prefix meaning "all following bytes". Alias: all.
	This especially avoids the looping at the end of the action list.
+ -
	deletes PFX bytes.
+ p
	prints PFX bytes.
+ r 
	changes PFX bytes to a random value (which may be the original one)
+ R
	changes PFX bytes to a random value (different from the original one)
+ =BY 
	changes PFX bytes to BY.

+ +BY
	inserts PFX times BY at the current position.

+ +r
	inserts PFX times a random byte at the current position.
+ |BY
	"or"s the current byte with BY.
+ oBY
	"or"s the current byte with BY (an alias for ease of use).
+ &BY
	"and"s the current byte with BY.
+ aBY
	"and"s the current byte with BY.
+ ^BY
	"xor"s the current byte with BY.
+ vBY
	search for the PFXth occurance of BY and print all bytes before it ("v" as in "verbose search").
+ qBY
	two searches for the same byte must not follow each other directly, as this would mean an endless loop. The program in that case bails out with an error.

	search for the PFXth occurance of BY without printing the bytes before it ("q" as in "quiet search").

end
	is the same as `* p` (print all until reaching the end of the input).

### Bytes

`BY` as used above is a bytes value. By default it is an decimal number ("64" stands for '@'),
but you may use 'x' and 'o' as prefixes for hexadecimal and octal numbers, 
and you may use 'c' as prefix for a character used literally. "64", "c@",
"x40" and "o100" all mean the very same thing, the well known character 
separating localpart and domain of an email address.

## Examples

### remove the first byte of input

	echo 123 | streamedit - \* p

	"minus to delete a byte, and `* p` to print all others."

### remove all X

	echo abcXdefXghiXjkl | streamedit vcX -

	"verbose search for character X, then delete, then repeat at the start."

### replace all X with Y

	echo abcXdefXghiXjkl | streamedit vcX =cY

	"verbose search for character X, then set to character Y, then repeat at the start."

### remove X to second X, and third X to fourth X

	echo abcXdefXghiXjklXmno | streamedit vcX - qcX -

	"verbose search for character X, then set to character Y, then repeat at the start."

### randomly remove one of byte four to seven, print all else.

	echo abcXdefXghiXjklXmno | streamedit 4 p 4r p - \* p

	"print 4 bytes, print 0..3 bytes, delete one, print the remainder."

### randomly remove 0 to three bytes starting from position 4, print all else.

	echo abcXdefXghiXjklXmno | streamedit 4 p 4r - \* p

### some random change every 0 to 4 bytes

	echo abcXdefXghiXjklXmno | streamedit 5r ~ | od -t x1
	
	The `od` is there to protect your terminal.


Report bugs to either your distributor, <uwe@ohse.de>,
or <https://github.com/UweOhse/streamedit/issues>.
Homepage: <https://ohse.de/uwe/software/streamedit/streamedit.html>
