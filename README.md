% streamedit(1)
% Uwe Ohse
% 2020-03-07

# NAME

streamedit - a byte stream editor

`streamedit [options] actions ... <input >output`

# DESCRIPTION

streamedit is a byte oriented stream editor, designed for manipulations of byte streams. It has no understanding of any higher level, does not understand UTF8, character sets, encodings (or whatever you fancy), and changes bytes even if that breaks higher level encodings. In fact, breaking higher levels is its main purpose.

# INSTALLATION

     git clone https://github.com/UweOhse/streamedit
     cd streamedit
     make
     `su`
     `make install`

_Warning_ : if a complete stranger can get you `su` and execute `make`, then you do need to worry about security. ___Really___, i mean it. `make install` may silently install a backdoor, a trojan horse or a keylogger on your system.

# USAGE

See ./streamedit.md

# LICENSE

GPLv2, 

# EXAMPLES

## remove the first byte of input

streamedit v6 p v6 - </dev/tty | lsz --xmodem -a somefile.txt >/dev/tty

"delete every second NAK character found in the input stream (this would make binary transfers impossible, and slow down text file transfers to a crawl. oh well, xmoden isn't exactly fast anyway).

# BUGS

None are known at the time of writing, but alas, all useful software has bugs.
If you find one, please report it either to your distributor, <https://github.com/UweOhse/streamedit/issues>, or <uwe@ohse.de>.

# NOTES

streamedit is a by-product of the lrzsz(1) package, and has been designed to allow simulations of line errors.


