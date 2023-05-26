# Developer notes

This application is primarily developed under Linux, but will run
under Windows. Please let me know if it also works on other
platforms.

The program grew using the good old Copy and Paste method. It is
not designed to be a lighthouse of great architecture, but to 
solve a problem.

The man-page requires the program go-md2man and sed/awk.
The man-page is essentially a patched version of the README.md.

https://github.com/cpuguy83/go-md2man.git

Because it's not assumed that the program is available usually,
the generated man-page is part of the release.

I am going to improve the openssf-scorecard rating and the
code-coverage of the tests.

Also, I might try to get the tool into NixOS.



