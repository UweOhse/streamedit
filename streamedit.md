% streamedit(1)
% Uwe Ohse
% 2020-03-07

# NAME

streamedit - a byte stream editor

# SYNOPSIS

streamedit [options] actions ... <input >output

# DESCRIPTION

streamedit is a byte oriented stream editor, designed for manipulations of byte streams. It has no understanding of any higher level, does not understand UTF8, character sets, encodings (or whatever you fancy), and changes bytes even if that breaks higher level encodings. In fact, breaking higher levels is its main purpose.

# OPTIONS

#### --help, -h

show usage information.

#### --seed int

set the random number seed (-1: by the current nanosecond, which is the default, anyway). This is used inside the test suite.

#### --verbose

be verbose about changes. It prints input position, output position, and what it does. This sometimes helps to understand what is going on. 

#### --version

show the version information.

# ACTIONS

`actions ...` is a series of actions, processed in order. If the program reaches the end of the actions, it will continue with the first one, unless `end` or `*` has been encountered.

An action is either:

#### PREFIX

The prefix for the next action (only), defaulting to 1.

PREFIX is a number. 1023, x3ff and o1777 all mean the same thing. The prefix determines how often the next action will be executed.

If you append a `@` sign to PREFIX, then PREFIX determines up to which position in
the input stream the next action will be executed (counted from 0). If that
position has already been passed, the action will not be executed.

If you append an literal `r` to PREFIX, the meaning changes to 'a random number of bytes
below N' (0 to N-1 bytes). _Note_: This may change in the future to mean 0 to N bytes.

#### *
A special prefix meaning "all following bytes". Alias: all.
This especially avoids the looping at the end of the action list.

#### -
Deletes PREFIX bytes.
#### p
Prints PREFIX bytes.
#### r 
Changes PREFIX bytes to a random value (which may be the original one)
#### R
Changes PREFIX bytes to a random value (different from the original one)
#### =BYTE 
Changes PREFIX bytes to BYTE.

#### +BYTE
Inserts PREFIX times BYTE at the current position.

#### +r
Inserts PREFIX times a random byte at the current position.
#### |BYTE
"Or"s the current byte with BYTE.
#### oBYTE
"Or"s the current byte with BYTE (an alias for ease of use).
#### &BYTE
"And"s the current byte with BYTE.
#### aBYTE
"And"s the current byte with BYTE.
#### ^BYTE
"Xor"s the current byte with BYTE.
#### vBYTE
Search for the PREFIX occurance of BYTE and print all bytes before it ("v" as in "verbose search").

Two searches for the same byte must not follow each other directly, as this would mean an endless loop. The program in that case bails out with an error.
#### qBYTE

Search for the PREFIX occurance of BYTE without printing the bytes before it ("q" as in "quiet search").

#### end
is the same as `* p` (print all until reaching the end of the input).

# BYTES

`BYTE` as used above is a bytes value. By default it is an decimal number ("64" stands for '@'),
but you may use 'x' and 'o' as prefixes for hexadecimal and octal numbers, 
and you may use 'c' as prefix for a character used literally.

"64", "c@",
"x40" and "o100" all mean the very same thing, the well known character 
separating localpart and domain of an email address.

# EXAMPLES

## remove the first byte of input

streamedit v6 p v6 - </dev/tty | lsz --xmodem -a somefile.txt >/dev/tty

"delete every second NAK character found in the input stream (this would make binary transfers impossible, and slow down text file transfers to a crawl. oh well, xmoden isn't exactly fast anyway).

## remove all X

echo abcXdefXghiXjkl | streamedit vcX -

"verbose search for character X, then delete, then repeat at the start."

## replace all X with Y

echo abcXdefXghiXjkl | streamedit vcX =cY

"verbose search for character X, then set to character Y, then repeat at the start."

## remove X to second X, and third X to fourth X

echo abcXdefXghiXjklXmno | streamedit vcX - qcX -

"verbose search for character X, then set to character Y, then repeat at the start."

## randomly remove one of byte four to seven, print all else.

echo abcXdefXghiXjklXmno | streamedit 4 p 4r p - \* p

"print 4 bytes, print 0..3 bytes, delete one, print the remainder."

## randomly remove 0 to three bytes starting from position 4, print all else.

echo abcXdefXghiXjklXmno | streamedit 4 p 4r - \* p

## some random change every 0 to 4 bytes

echo abcXdefXghiXjklXmno | streamedit 5r ~ | od -t x1
	
The `od` is there to protect your terminal.

# BUGS

None are known at the time of writing, but alas, all useful software has bugs.
If you find one, please report it either to your distributor, <https://github.com/UweOhse/streamedit/issues>, or <uwe@ohse.de>.

# NOTES

streamedit is a by-product of the lrzsz(1) package, and has been designed to allow simulations of line errors.


