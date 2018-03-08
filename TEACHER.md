To use this software to share code in class, you will need to (1) install [Sublime Text 3](https://www.sublimetext.com/3) and (2) install a specific plug in for Sublime Text.

To install Sublime Text 3, [go here.](https://www.sublimetext.com/3)

To install the necessary Sublime Text plug in, follow these steps:

+ Open Sublime Text
+ Click Show Console in the View menu.
+ Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4BrowniesInstructor.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); c4b_version = os.path.join(package_path, "VERSION"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Code4BrowniesInstructor.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Main.sublime-menu", c4b_menu); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
```
+ Paste copied code to Console and hit enter.

### Code sharing during class

First, download the server and run it on the instructor's machine.

- [OSX amd64](https://umdrive.memphis.edu/vphan/public/C4B/c4b_osx_0.43).
- [Win amd64](https://umdrive.memphis.edu/vphan/public/C4B/c4b_win_0.43.exe).

The server will automatically create a database to store teacher's and students' shared code, brownie points, and other information.

Change the permission of the file to executable.

### Running the server on the instructor's laptop

Students and the instructor communicate by sending messages to a server.  The server should be run on the the instructor's computer.

OSX: run the server in a terminal
```
    ./c4b_osx_0.43
````

Windows: run the server in a terminal
```
    ./c4b_win_0.43.exe
````

If you want to run the server with the source code, you need to install Go.  To run the server:
```
    ./go run *.go
````


## Uninstallation of Sublime Text plugin

Open Sublime Text, go to View, click Show Console, copy this code, paste to console and hit enter:

```
import os; import shutil; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); shutil.rmtree(package_path)
```
