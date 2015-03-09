# go-php
Go language wrapper to download PHP runtime and run PHP phars. Just like Java Web Start does for java jars.

What does this do?

1. It downloads a php script (phar archive created by the box-project).
2. Then it checks for PHP runtine (is php.exe in PATH).
3. If not, it downloads the current PHP runtime and runs the phar file.

Notes.

- The current version only redownloads the phar archive every time the go-php.exe wrapper is run but a standard manifest.json should be implemented to check for updates instead.
