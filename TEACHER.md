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

### Get the latest server

To use Code4Brownies as a teacher, first, download the server and run it on the instructor's machine.

- [OSX amd64](https://www.dropbox.com/s/g9xsjgwqhqcdook/c4b_osx_0.45?dl=0)
- [Win amd64](https://www.dropbox.com/s/wb27tnckvmzt0ab/c4b_win_0.45.exe?dl=0)

The server will automatically create a database to store teacher's and students' shared code, brownie points, and other information.

Change the permission of the file to executable.

### Run the server on the instructor's laptop

Students and the instructor communicate by sending messages to a server.  The server should be run on the the instructor's computer.

OSX: run the server in a terminal
```
    ./c4b_osx_0.XX
````

Windows: run the server in a terminal
```
    ./c4b_win_0.XX.exe
````

If you want to run the server with the source code, you need to install Go.  To run the server:
```
    ./go run *.go
````

## Teach with Code4Brownies

To begin a class, run the server.  Once this is done, the instructor can enable a few activities.

(A) Scaffold and share work with students

1. Teacher shares a previously prepared note.  For example, a partially complete piece of code for student to work on as an in-class exercise.

2. Teacher looks at student submissions.

3. Teacher provides feedback to a student, if he or she desires.

4. Teacher gives brownie points to students.

A [teaching assistant](ASSISTANT.md) can help with steps 2, 3 and 4.

(B) Administer quizzes and polls

(C) Have students take attendance by themselves

(D) Share notes with [teaching assistants](ASSISTANT.md)

(E) Monitor submissions and polls

## Uninstallation of Sublime Text plugin

Open Sublime Text, go to View, click Show Console, copy this code, paste to console and hit enter:

```
import os; import shutil; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); shutil.rmtree(package_path)
```
