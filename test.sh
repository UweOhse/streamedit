#! /bin/sh
docover()
{
	if test "$COVER" = 1 ; then
		gocovmerge cover.out cover1.out >cover.t.out 2>/dev/null
		mv cover.t.out cover.out
	fi
}
exec 2>&1 
S=./streamedit
if test "$COVER" = "1" ; then
	S="./streamedit-covering -test.coverprofile=cover1.out "
	echo | $S --version
	mv cover1.out cover.out
	echo | $S --help
	docover
fi

echo "# abuse"
printf "abcdef" | $S garbage ; echo $? ; docover
printf "abcdef" | $S "" ; echo $? ; docover
printf "abcdef" | $S ; echo $? ; docover
echo

echo "# r works" # parse and operation
printf "abcdhj" | $S --seed 1 r | od -t a ; docover
printf "abcdhj" | $S --seed 2 r | od -t a ; docover
echo

echo "# R works (no output is good)" # parse and operation
printf "abcdhj" | $S --seed 2 r | od -t a ; docover
printf "abcdhj" | $S --seed 2 R | od -t a ; docover
echo

echo "# - works" # parse and operation
printf "abcdef" | $S - end ; echo ; docover
echo

echo "# N works" # parse and operation
printf "abcdef" | $S 0 - end ; echo ; docover # delete 0 elements
printf "abcdef" | $S 1 - end ; echo ; docover # delete 1
printf "abcdef" | $S 4 - end ; echo ; docover # delete 4
printf "abcdef" | $S o3 - end ; echo ; docover
printf "abcdef" | $S o2 - end ; echo ; docover
printf "abcdef" | $S x01 - end ; echo ; docover
printf "abcdef" | $S -- -1 - end ; echo ; docover
echo

echo "# + works" # parse and operation
printf "abcdef" | $S 1 +c@ end ; echo ; docover
printf "abcdef" | $S 2 +64 end ; echo ; docover
printf "abcdef" | $S 3 +x40 end ; echo ; docover
printf "abcdef" | $S 4 +o100 end ; echo ; docover
printf "abcdef" | $S 2 p 1 +c@ end ; echo ; docover
printf "abcdef" | $S 2 p 2 +64 end ; echo ; docover
printf "abcdef" | $S 2 p 3 +x40 end ; echo ; docover
printf "abcdef" | $S 2 p 4 +o100 end ; echo ; docover
printf "abcdef" | $S +o end ; echo ; docover
printf "abcdef" | $S +o400 end ; echo ; docover
printf "abcdef" | $S +x end ; echo ; docover
printf "abcdef" | $S +x100 end ; echo ; docover
printf "abcdef" | $S +c end ; echo ; docover
printf "abcdef" | $S +caa end ; echo ; docover
printf "abcdef" | $S + end ; echo ; docover
echo

echo "# + does not append to line end" # parse and operation
printf "abcdef" | $S 6 p 1 +c@ end ; echo ; docover
printf "abcdef" | $S 5 p 1 +c@ end ; echo ; docover
echo

echo "# +r works"
printf "abcdef" | $S --seed 1 2 p 2 +r end | od -t a ; docover
printf "abcdef" | $S --seed 2 2 p 2 +r end | od -t a  ; docover
echo

echo "# = works" # parse and operation
printf "abcdef" | $S =c@ end ; echo ; docover
printf "abcdef" | $S 1 =64 end ; echo ; docover
printf "abcdef" | $S 2 =x40 end ; echo ; docover
printf "abcdef" | $S 3 =o100 end ; echo ; docover
printf "abcdef" | $S 1 = end ; echo ; docover
echo

echo "# Nr works" # parse and operation
printf "abcdef" | $S --seed 1 '3r' - end ; echo ; docover
printf "abcdef" | $S --seed 2 '3r' - end ; echo ; docover
printf "abcdef" | $S --seed 2 '0r' - end ; echo ; docover
printf "abcdef" | $S --seed 2 '1r' - end ; echo ; docover
echo

echo "# | works" # parse and operation
printf "abcdef" | $S "|2" end ; echo ; docover
printf "abcdef" | $S 1 "|cc" end ; echo ; docover
printf "abcdef" | $S 2 "|x64" end ; echo ; docover
printf "abcdef" | $S 3 "|o155" end ; echo ; docover
echo "# O works" # parse and operation
printf "abcdef" | $S "O2" end ; echo ; docover
printf "abcdef" | $S 1 "Occ" end ; echo ; docover
printf "abcdef" | $S 2 "Ox64" end ; echo ; docover
printf "abcdef" | $S 3 "Oo155" end ; echo ; docover
echo

echo "# & works" # parse and operation
printf "abcdef" | $S "&95" end ; echo ; docover
printf "abcdef" | $S 1 "&c_" end ; echo ; docover
printf "abcdef" | $S 2 "&x5f" end ; echo ; docover
printf "abcdef" | $S 3 "&o137" end ; echo ; docover
echo

echo "# A works" # parse and operation
printf "abcdef" | $S "A95" end ; echo ; docover
printf "abcdef" | $S 1 "Ac_" end ; echo ; docover
printf "abcdef" | $S 2 "Ax5f" end ; echo ; docover
printf "abcdef" | $S 3 "Ao137" end ; echo ; docover
echo

echo "# ^ works" # parse and operation
printf "abcdef" | $S "^32" end ; echo ; docover
printf "abcdef" | $S 1 "^c " end ; echo ; docover
printf "abcdef" | $S 2 "^x20" end ; echo ; docover
printf "abcdef" | $S 3 "^o40" end ; echo ; docover
echo

echo "# | works with high bit" # parse and operation
printf "abcdef" | $S "|128" | od -t x1 ; echo ; docover
printf "abcdef" | $S "|o200" | od -t x1 ; echo ; docover
printf "abcdef" | $S "|x80" | od -t x1 ; echo ; docover
echo

echo "# v works" 
printf "abcdef" | $S vcd 2 - end ; echo ; docover
printf "abcdef" | $S v98 3 =c@ end ; echo ; docover
printf "abcdef" | $S vx61 3 =c@ end ; echo ; docover
printf "abcdef" | $S vx61 3 =c@ vce - end ; echo ; docover
printf "abcdef" | $S v 3 =c@ vce - end ; echo ; docover
echo

echo "# v works with pfx" 
printf "abcdefabcdef" | $S 2 vcd 2 - end ; echo ; docover
printf "abcdefabcdef" | $S 2 v98 3 =c@ end ; echo ; docover
printf "abcdefabcdef" | $S 2 vx61 3 =c@ end ; echo ; docover
printf "abcdefabcdef" | $S 2 vx61 3 =c@ vce - end ; echo ; docover
printf "abcdefabcdef" | $S 2 v 3 =c@ vce - end ; echo ; docover
echo

echo "# q works" 
printf "abcdef" | $S qcd 2 - end ; echo ; docover
printf "abcdef" | $S qcf p end ; echo ; docover
printf "abcdef" | $S q98 3 =c@ end ; echo ; docover
printf "abcdef" | $S qx61 3 =c@ end ; echo ; docover
printf "abcdef" | $S qx61 3 =c@ qce - end ; echo ; docover
printf "abcdef" | $S qx61 3 =c@ vce - end ; echo ; docover
printf "abcdef" | $S q 3 =c@ qce - end ; echo ; docover
echo

echo "# p works" 
printf "abcdef" | $S 2 p - ; echo ; docover
printf "abcdef" | $S p p - p end ; echo ; docover
printf "abcdef" | $S p p - end ; echo ; docover
printf "abcdef" | $S - 3 p - p ; echo ; docover
echo

echo "# binary data works" 
printf "\000\001\002\003\004\005\006\007\010\011\012\013\014\015\016\017" \
	| $S "|64" ; echo ; docover
printf "\300\301\302\303\304\305\306\307\310\311\312\313\314\315\316\317" \
	| $S "^x80" > test.tmp ; docover
	$S "|o40" <test.tmp ; echo ; docover
echo 

echo "# @ works"
printf "abcdef" | $S 0@ - end ; echo ; docover
printf "abcdef" | $S 1@ - end ; echo ; docover
printf "abcdef" | $S 3@ - end ; echo ; docover
printf "abcdef" | $S 4@ - end ; echo ; docover
printf "abcdef" | $S p p 4@ - end ; echo ; docover
printf "abcdef" | $S p p 4@ - ; echo ; docover
echo 

echo "# loops work"
printf "abcdef" | $S p - ; echo ; docover
printf "abcdef" | $S - p ; echo ; docover
printf "abcdefabcdef" | $S 2 p 2 - ; echo ; docover
echo 

echo "# append to end works"
printf "abcdef" | $S end +64 ; echo ; docover
printf "abcdef" | $S end +64 +64 ; echo ; docover
printf "abcdef" | $S end 5 +64 ; echo ; docover
printf "abcdef" | $S --seed 1 end 2 +r | od -t x1 ; docover
printf "abcdef" | $S --seed 2 end 2 +r | od -t x1 ; docover
# append a random number of random bytes.
printf "abcdef" | $S --seed 3 end 7r +r | od -t x1 ; docover
echo 

echo "# anything else after end fails"
printf "abcdef" | $S end - ; echo ; docover
printf "abcdef" | $S end p ; echo ; docover
printf "abcdef" | $S end r ; echo ; docover
printf "abcdef" | $S end R ; echo ; docover
printf "abcdef" | $S end "|1" ; echo ; docover
printf "abcdef" | $S end "O1" ; echo ; docover
printf "abcdef" | $S end "&1" ; echo ; docover
printf "abcdef" | $S end "A1" ; echo ; docover
printf "abcdef" | $S end "^1" ; echo ; docover
printf "abcdef" | $S end =1 ; echo ; docover
printf "abcdef" | $S end vca ; echo ; docover
printf "abcdef" | $S end qca ; echo ; docover
printf "abcdef" | $S end 6@ ; echo ; docover
echo

echo "# wildcards work"
printf "abcdef" | $S end ; echo ; docover
printf "abcdef" | $S all - ; echo ; docover
printf "abcdef" | $S \* A95 ; echo ; docover

echo "# paste works"
printf "abcd" | $S p - p P p  ; echo ; docover
printf "abcd" | $S p - \* p +64 P ; echo ; docover
printf "abcd" | $S p - p end P +64 ; echo ; docover
printf "abcd" | $S p - p p end +64 P ; echo ; docover

echo "# the X...X removal example"
echo abcXdefXghiXjklXmno | $S vcX - qcX - ; echo ; docover

echo "# two searches must not directly follow each other"
echo abcXdefXghiXjklXmno | $S vcX - qcX ; echo ; docover

echo "# ~ handles delete (0)"
printf "abcdef" | $S --seed 11  p "~" "*" p | od -t x1 ; docover
printf "abcdef" | $S --seed 15  p "~" "*" p | od -t x1 ; docover
echo "# ~ handles insert (1)"
printf "abcdef" | $S --seed 1  p "~" "*" p | od -t x1 ; docover
printf "abcdef" | $S --seed 2  p "~" "*" p | od -t x1 ; docover
echo "# ~ handles flip (2)"
printf "abcdef" | $S --seed 45  p "~" "*" p | od -t x1 ; docover
echo "# ~ handles set (3)"
printf "abcdef" | $S --seed 3  p "~" "*" p | od -t x1 ; docover
printf "abcdef" | $S --seed 6  p "~" "*" p | od -t x1 ; docover
echo "# ~ handles bitflip (4)"
printf "abcdef" | $S --seed 4 p "~" "*" p | od -t x1 ; docover
printf "abcdef" | $S --seed 10 p "~" "*" p | od -t x1 ; docover

