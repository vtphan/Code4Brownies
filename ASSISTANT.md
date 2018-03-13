To use this software to share code in class, you will need to (1) install [Sublime Text 3](https://www.sublimetext.com/3) and (2) install a specific plug in for Sublime Text.

To install Sublime Text 3, [go here.](https://www.sublimetext.com/3)

To install the necessary Sublime Text plug in, follow these steps:

+ Open Sublime Text
+ Click Show Console in the View menu.
+ Copy this code:
```
import os; package_path = os.path.join(sublime.packages_path(), "C4BAssistant"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4BrowniesAssistant.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); c4b_version = os.path.join(package_path, "VERSION"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BAssistant/Code4BrowniesAssistant.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BAssistant/Main.sublime-menu", c4b_menu); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
```
+ Paste copied code to Console and hit enter.

## Uninstallation of Sublime Text plugin

Open Sublime Text, go to View, click Show Console, copy this code, paste to console and hit enter:

```
import os; import shutil; package_path = os.path.join(sublime.packages_path(), "C4BAssistant"); shutil.rmtree(package_path)

## Using Code4Brownies as a teaching assistant

(1) The main workflow of a TA is as follows:

+ Get student submissions.  One at a time; a few at a time; or all at once.
+ Optionally, comment on it to give feedback to a student.
+ Give brownie points to each submission.  Or dequeue or requeue or ungrade it.

(2) Get and share notes with the teacher.

(3) Track how many submissions are on the queue.

```
