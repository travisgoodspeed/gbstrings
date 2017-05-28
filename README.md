Howdy y'all,

This is a quick little tool that I tossed together one night for
finding GB2312 Chinese strings from the memory of an imported ham
radio.  You might find it handy when translating old video games, as
well.  (GB2312 is not Unicode, and far better tools exist for locating
Chinese Unicode strings.)

I cannot speak Chinese, so it's quite likely that you can improve upon
this tool.  Pull requests are welcome.

Some brief notes on GB2312 follow.  This textfile, and all files of
the project except binaries, are encoded in UTF8.

--Travis


## Software Notes

This is a standard Golang project, which you can install by `go get
github.com/travisgoodspeed/gbstrings` after setting your `$GOPATH`.
It works well on Linux, but will likely fail on Windows.

## GB2312 Notes

GB2312 (国家标准2312) is Chinese national standard character encoding
that stores Chinese symbols as byte pairs that may freely mingle with
ASCII.

In GB2313, a glyph follows these rules.

1. The first byte is from 0xA1 to 0xF7.
2. The second byte is from 0xA1 to 0xFE.

Glyphs and ASCII may intermix, but a glyph will always be two bytes
long and there must be no multi-byte strings outside the allowed
range.

Because most glyphs do not sit in the outer corners of the range, the
range differences are not very helpful for determining alignment.  A
better technique is to look at the start and the end of a potential
string, ensuring that neither side overlaps an illegal byte or an
ASCII byte.

Libiconv is very powerful for identifying that a string is valid,
legitimate GB2312, but it is very slow at validating or refusing a
string.  For that reason, this tool first attempts to approximate a
string's legitimacy with these loose rules, then goes back to confirm
the validity with libiconv before exporting the string.

Font support in terminals can vary dramatically.  I've found that
`uxterm` works best with the Large font size.



## Related Projects

For Unicode strings, you're far better off with Radare2's `rabin2 -zzq
foo.bin` technique.

The IDA Pro [Translator
Plugin](https://github.com/kyrus/ida-translator/wiki/Introducing-the-IDA-Pro-Translator-Plugin)
can be helpful when reverse engineering foreign strings.  They provide specific examples for
numerous regional text standards.

[APT Friend Finder](https://www.aptfriendfinder.com/) is a Perl script
that finds and translates false positives of Chinese text within
binaries.  It was written to joke about China being blamed for all
malware, back in the days before Russia was to blame for everything.

