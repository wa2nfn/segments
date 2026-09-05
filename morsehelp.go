package main

import (
	"fmt"
	"strings"
)

func morseHelp() {

	ClearScreen()
	screen.Printf(`

Introduction
============

The purpose of the application is to give a wood turner a mechanized way to design a
ring - decorated or embellished by creating a layout of 2-3 different colors of wood. This
single ring can be added to a standard wood-blank bowl, or it can replace a ring within a
multi-ring segmented bowl.  The mathematics to develop the amount of segments, as well
their size and cut angle, require the bowl to be circular. All segments of the given ring
are of the same width and height regardless of color.

Before we go into too much about segments, let me say there will be a questions asking if
this is for Basket Illusions, thats a design on a bowl typicalls done with wood burning and
manually coloring segments instead of cutting them. Therefore this option will bypass the 
need to supply a number of measurments needed for cut segments, as such this option can be
useful in several ways:

1- it can develop the pattern for the basket illusion as the name implies
2- it can be a means to a quick look at segmented wood pattern without the specifics to make
   the physical ring
3- it gives a means to create the pattern for use on flat work such as a cutting board

While the word singular means one ring, a user can add or replace any number of rings in
a bowls' design, but each execution of the application only does a single ring. That single
ring has: a fixed diameter, wall width (or thickness), and a padding size (to give 
flexibility in wall width).

The precision used in length and width sizes is 32nds of an inch. (See Miscellaneous)

General Information
===================

You can give a quick read to get the high level concepts, then try the app and return here
as needed.

There will be a lot of numbers in the following explanation; they are for understanding.
The app uses built-in values as appropriate. When you have a input choice to make, the screen 
will guide you with allowable values. Entries of lengths, widths, etc. are in inches; they nay be 
entered several formats: 7, 7.0, 7-1/2, 7 1/2, 0.25.

A circular ring is composed of a number of segments. The segments are grouped into patterns and
patterns are separated by spaces (called inter-character spaces, inter-word spaces, and ending spaces).
Ending spaces could be the inter-character space (in 'D' design pattern mode), the word space (in 'M' or
'E' mode), either may have extra spaces, as will be described.  We can then 
model a general ring like: <pattern><inter-char. space>(both of those repeated 1 or more 
times)<ending space>. Following is a slightly more specific example model with:

   - 3 patterns
   - inter-char.(or pattern) space of 3
   - ending space of 7

   It's shown within quotes, to show the start/end, and see the end spaces.

   "<pattern A><   ><pattern B><   ><pattern C><      >". We'll refine the pattern details shortly. 
   Note: you'll see that I use "inter-character spaces" and yet they appear between patterns, that's
   because an entire pattern is indicated by typing a single character.

The approach taken in the design of a ring differs from the typical segmented bowl approach.
Typically, a ring's diameter is determined by the shape of the bowl, with that diameter, a 
number of segments of wood (very often the same as every other ring). All that's left 
to the designer's creativity is to decide what woods to use to make some geometric pattern
(if any). 

This app changes that order a little. With some constraints, it gives the designer more 
latitude in the aesthetic look of the ring. The app will collect some of the same 
characteristics: diameter, wall width, and a padding factor (that may have been 
automatically applied by the existing bowls design). The big difference is that the 
number of segments is NOT given, its derived from completed ring design pattern.

The ring is designed by selecting one or more of the 50 built-in patterns. The patterns are
constructed from segments of DARK wood (shown as 'D' or a dark square) and LIGHT wood 
(shown as 'L' or a light square). We can generally think of the LIGHT wood as the background
color. DARK and LIGHT are subjective labels. In a decorative pattern the choice of woods is
purely aesthetic.  <See Miscellaneous for color viewing options)

Pattern Details
===============

1 - All segments are the SAME length (segment edge length or SEL).

2 - A pattern can be as short as 1 segment, or as long as 17 (details to follow).

3 - patterns START and END with a DARK segment (e.g. DLD, or DDDLDDD).

4 - DARK portions will be 1 or 3 segments long (not changeable).

5 - DARK portions will be separated by 1 LIGHT segment (not changeable).

6 - A pattern uses at most 2 colors. 

The rules above are for all 50 patterns and are invariant . The app offers a way to reverse the colors of the
entire ring, if desired.

Ring Designs From Pattern(s)
===========================

The pattern is the fundamental building block of the ring, but now we need to flesh out the
ring with required spacings which separate patterns.

An entire ring can be designed to be one or more of the built-in patterns. The number and
choice is only constrained by the final segment count.

If the ring consists of more than 1 pattern, then patterns MUST be separated.  The pattern
separation, and whether you have some flexibility, is determined by mode of construction.

Mode
====

There are 3 modes to contruct the ring: D for decorative Pattern; M for Morse Code; and E for
Enhanced Morse Code. Morse Code? Don't panic! You don't have to learn it, and your bowl
certainly is not going to make audible sounds.

All 3 modes allow you to combine any number of the built-in patterns. The choice of mode
determines the flexibility in spacing (LIGHT, or OTHER colored segments) in between 
patterns. D is the most flexible since it is just decorative art; M is the least, because 
Morse Code is in effect a language.  So E is a hybrid for those that want to convey a 
message in Morse Code by combing patterns, but want a little more control over some of the
inter-character/pattern spacing and the appearance it makes on the bowl.

You may be feeling this is very abstract. Let me show you how comfortable you can be with this
without even knowing it. Let's think about the sentence that you just read. Each letter is
in effect a "pattern or symbol"; those patterns have a small sliver of space between each
other so that they don't overlap (this is built into the font).  And in between words, we
clearly have the word separator space. All of these have direct correlation to the spacing
the app does.

In mode E, some of the strict spacing requirements of Morse Code spacing are relaxed, this can 
reduce the total segment count.

In addition, the app will auto adjust some of the possible end of message spacing segments
so they are uniformly distributed for a more balanced design. (Extra spaces will be discussed.)

  ** Why ?. because Morse is audible, and ears are not as discriminating as our eyes. So 7
  is very distinctive compared to the 3 between letters, but on a bowl the word space that
  is at least 2 segments longer than the other should be OK.

Each of the 50 patterns has a simple 1 character name which is what that pattern represents
in Morse Code - simpler to say or type 'R' rather than, pattern number 33.

(This can be skipped on initial reading,) In M mode, a prompt asks if you want EXTRA end spaces.
This may make sense for a unique situation where you have a message on two rings, AND the ring 
diameters are similar, AND the number of segments is very DIFFERENT, AND you want SELs to be 
similar. Example: if the two rings were the same diameter, and the upper ring had "LUV" and
the ring below it had "SEGMENTED BOWLS", clearly the first line has many fewer segments;  
therefore, they have a longer SEL, than the line below. The layout would be unbalanced .
What can be done is run the app with the larger number of segments first, then run the app for
the smaller one. Do the mental math to see how many segments to add to the small segment case.
This and/or playing with end spacing can bring the SELs of the two rings into closer agreement.

Below are three of the 50 patterns to give a better visual explaination.

`)
	A := colorize("DLDDD")
	P := colorize("DLDDDLDDDLD")
	T := colorize("DDD")
	screen.Printf(" 'A' .-     DLDDD        %-10s         4   1   5\n\n", A)
	screen.Printf(" 'P' .--.   DLDDDLDDDLD  %-10s   8   3  11\n\n", P)
	screen.Printf(" 'T' -      DDD          %-10s           3   0   3\n\n", T)

	screen.Printf(`

Why Different Modes
===================

The decorative pattern is like a bracelet. There is no start or end, like "D  D  D  D  ".
But Morse Code is meant to create a "message" (to those that can read it) like the inside of
a wedding band for example. In this case the start and end ARE important, and this is 
indicated by a longer blank space.

For example, in the message "Pat and Bill     ", we can have the final inter-word space 
longer, or the "message" (when on a circular ring) could be misread like " and Bill Pat " 
since there is no obvious starting point.

When you enter the first character, after seeing the message/pattern prompt, there will be
a <> followed by a (). Within the <> is the current number of segments that the message 
translates to. This allows you to modify your message if the count gets too high. In 
addition, the () includes the number of additional segments that occur AFTER the message.
For example: you type: CAT, the screen will show: CAT <23> (+7>, meaning you need at least 30
segments for that message (why at least will be explained).

Running The App
=============== 

A few questions ask for a 'Y' for yes or 'M' for Morse Code. You can use upper or 
lowercase. RETURN gets you the default choice.

Your input message/pattern can use the characters listed on the screen.  Unsupported 
characters will be silently ignored. The space ' ' is special: it's the only invisable 
character supported. They are ignored if at the beginning or end (they are accounted for 
else where), and ignored if there are more than one in sequence.  So, if you enter 
"    hello     there    ", it's reduced to "hello there". (ending space is app applied)

When you enter the first character at the message/pattern prompt, the line format changes
to "" <0> (+X). Within the quotes will be your first character; within the diamond will be 
the current running total of SEGMENTS for the pattern (not characters), the (+X) will be
static. It's the number of additional spacing segments that are added after you hit RETURN
(its the number you just entered in the previous prompt. For example, let's assume you used the
decorative mode 'P' and you gave 2 as the inter character space, then you enter 'A'. The line
changes to: "A" <5> (+2). Because A=5 segments from the pattern chart. When you enter a 
second 'A', the line chages to: "AA" <12> (+2). Why 12, not 10? Because there are 2
segments between the 'A' patterns as inter character spaces. You can Backspace to modify
your input. When RETURN is hit to end input, the value in () is added to the value in <> 
for a total (you won't see it on screen). Note: a space will increment <> by a previous input
value as well.

When the app uses your input to create the ring pattern, it will show a data under a heading
"Parsing Message". A data line for each character is just a confirmation of what the ring
design is composed of. It's really important if you don't know Morse Code, but you are making
a bowl with a message. A bad time for a wooden typo!

Diameter: You are asked for the finished ring diameter (in inches) as a mixed number or 
decimal value.  This determines the rings length, and with the calculated segments determines
the SEL.  Using your segmented bowl software or knowledge of your physical bowl and the wood
dimensions, the app will increase this to an appropriate value when you enter a "padding" 
thickness to both the inside and outside.  Consideration needs to be given to the impact of
the steepness of the bowl sides, thickness of the segments, and even turning. The more flat-sloped 
the bowl is (platter like) or wide-walled, at the ring's location, the padding becomes more
significant.  (See Miscellaneous)

Extra End Spacing and Miter Angle
=================================

If your input message/pattern creates a segment count that divides evenly into 180 degrees, 
then you will get a practical angle to cut your segments, like 22 degrees. Note: commercial 
wedges are usually marked with a SEGMENT COUNT and CENTRAL ANGLE. The CENTRAL ANGLE is 
TWICE the MITER ANGLE! Both are shown in the Ring Summary when the app is run. Be CAREFUL!

FYI: if segments are added to your message segment total (above in the <>) that will be
shown in a screen note.

It's possible that the math will result in an unacceptable angle, like 21.38 degrees. The
app will add more spaces to the end space and test for an acceptable angle.  This will 
only be done if the angle is not based on the segments of a commercial wedge: the 
built-in values are: 12, 16, 18, and 24 segments. If you have more or less than those 
available, then create a simple text file called "wedges.txt", in the current directory, 
with the segments you do have. One value per line, like:

10
16
40

In summary, if you see an asterisk next to the segment count, like 16*, this means that
the designs' angle matches an existing wedge; therefore, you don't have to be concerned about
a calculated angle.  Obviously, an non-asterisked segment count means you will need to 
create a wedge or use a protractor to set your fence for the cut angle.

A few simple examples:

In the table, column 1 is the "name" of the pattern (easy to say or type), column 2 is what 
the name in column 1 would "sound" like in Morse Code ('A' sounds like didah, where a di is
1/3 the length of a dah); column 3 is a schematic representation of the pattern as it would
be in wood segments; 4 is the same as 3 with color replacing the Ds and Ls. Finally, the last
3 columns are the count of DARK, LIGHT segments and the Total number of segments to create
the pattern.

Use your creativity to choose patterns and their order - remembering they are separated
with spacing segments (1-6) which changes the visual, just as: A A A vs A   A   A   A,
A     A, etc.

The space pattern ' ' is special in that it is longer than the 1-6 we already saw. In real 
Morse Code it is 7 times the length of the di(dit), which is the length of space used INSIDE
of a pattern.

One more limited option exists in terms of color. This option only applies in mode 'P' for 
creating a decorative ring. The choice to use it will only be presented when it applies.
When you run the app, if you are asked about a third color, we will refer to this one as 
OTHER or 'O', since it could be darker or lighter than the other two. It's just gives a 
little more variety.  If you elect to use it, it will ONLY be used as replacement for the
LIGHT segment(s) IN BETWEEN patterns (patterns themselves are always Ds and Ls).  A third
color "Other", would be inserted as follows: "<Pattern A>OOOO<Pattern B><End Spaces>".
This is shown in the last example in contrast to the example before it using the same patterns.

Examples
========

Pattern             # spaces between      Colorized Ring
                        patterns
=======             ================      ==============
`)
	screen.Printf("EEEEEE                    1              %s\n\n", colorize("DLDLDLDLDLDL"))
	screen.Printf("TTTTT                     1              %s\n\n", colorize("DDDLDDDLDDDLDDDLDDD"))
	screen.Printf("TTT                       3              %s\n\n", colorize("DDDLDDDLDDD"))
	screen.Printf(".                         1              %s\n\n", colorize("DLDDDLDLDDDLDLDDD"))
	screen.Printf("SOS                       1              %s\n\n", colorize("DLDLDLDDDLDDDLDDDLDLDLD"))
	screen.Printf("MIMI                      3              %s\n\n", colorize("DDDLDDDLLLDLDLLLDDDLDDDLLLDLDLLL"))
	screen.Printf("MIMI                      3              %s\n\n", colorize("DDDLDDDOOODLDOOODDDLDDDOOODLDOOO"))

	screen.Print(`
Miscellaneous
=============

    Wedges.txt File
    ===============

       To minimise editing of values in this file you can make the app ignore any line
       by starting the line with "#" or "//", rather than deleting a line for a temporary
       purpose. In addition, you can force the app to exit immediately rather than reading
       any additional line, this is done with starting a line with ANY of these 
       words: end, stop, quit, abort.

       The list of the wedge segments in this file, are one per line. The order does not 
       matter.  If the file exists, it will replace all of the internally stored values. 

       Wedges 
       ------
       If you have wedges, list them - not listing a segment number (like popular value 16)
       will cause the app to require many more segments to construct a ring than is 
       necessary. Why? - because the app will not specifiy a segment that doesn't have
       a whole degree angle, like 20, 22 since that would be difficult to measure.  A 16 
       segment wedge has a CENTRAL ANGLE of 22.5 degrees, and a MITER ANGLE of 11.25; but 
       the wedge manufacturer has done the accurate measuring.

       *** Putting a segment count for a wedge that you do NOT have, will simply return 
       the calculated angle (180/segments) and can produce a non-practical value like 10.58
       degrees. Common wisdom says an angle that is not within 0.1 degrees of accuracy will 
       give you gaps or a non circular ring.

       Precision
       ---------
       The wedge file can also have a non-wedge related input. A line with 
       precision=16, 32(default) or 64, will set the precision of length measurements to
       the value given. The actual lengths will be shown in reduced form, that is 1-1/4,
       not 1-16/64s. So the default is 32nds of an inch; if given in the file. 

       Color
       -----
       If you have issues seeing Brown/Tan/Beige used to color segments, you can use 
       color=blue, green, or red in the wedges.txt file. In each case, there is clearly
       a substitue for Dark and Light.  while Other (only available with Design Morse, 
       is spectrally in between.)

       Kerf
       ====
       Overrides the standard 1/8" kerf. Only makes impact to the raw board length needed.

       The wedges.txt file must be in the same directory as the segments.exe app file.
       (Linux: maybe called segments.)

    Padding
    =======

    Two padding parameters, Exterior and Interior, allow you to speficy some amount of 
    leeway in the thickness of the segments. This allowance can be for many reasons:
    accuracy of cutting, sanding, gluing, etc. Typically I'd expect the same value for
    both, but for maximum flexibilty they are collected separately. These work hand-in-hand
    with the ring width that they are added to to determine the board width that you will
    rip to cut your segments from. Since this just becomes a summation of values you
    can decide which variable you use. For example a ring width of 1/2" and, both interior
    and exterior padding of 1/4" makes a board with of 1"; the same as saying both padding
    values are 0 and the ring width is 1".

    Software
    ========

       Why over 2Mb? The app is written in the Go language, the file is completely self 
       contained; no programming environment, no setup, no install, no configuration file,
       no changes to your PC.

       Does it have a GUI or mouse support? No, this is a command line app, run it from any
       Terminal (command) window or Powershell.

       Removal: there is no remove script, use your standard OS remove command. If you save 
       output to a file when prompted the file is written to the current directory. The 
       file will have the suffix ".html". Use a sensible name so you can find and delete 
       them when you no longer need them.
       `)

	screen.Printf("\n\n%sSAFETY IS NUMBER ONE PRIORITY! Cut only segments that can be held safely and securely.%s\n\n", ANSI_RED, ANSI_RESET)
	screen.Printf("Typos, bugs, requested clarifications can be sent to: %swa2nfn@gmail.com%s with subject \"segments app\".", ANSI_RED, ANSI_RESET)
	screen.Printf("\n\nDisclaimer: This software is provided for free, for hobbyists to experiment with.\nIt is made available \"as is\" and without any warranty, expressed or implied.\nUse it at your own risk.\n")

	fmt.Printf("\n%sSCROLL TO THE TOP OF INFORMATION SECTION %s or \n", ANSI_GREEN, ANSI_RESET)
	fmt.Print("Press P to print or RETURN to continue: ")
	var input string
	fmt.Scanln(&input)

	if strings.ToUpper(input) == "P" {
		screen.PrintScreen()
	}
}
