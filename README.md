Code4Brownies
Author: Vinhthuy Phan, 2015

## Update of the Student's plug-in

Assuming that you have properly installed the plug-in.  The quick way is to Show Console (in View)
paste and run this code:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4Brownies.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu", c4b_menu)
```

If this gives an error (e.g. "os not found"), then save [Code4Brownies.py](https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py) and [Main.sublime-menu](https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu) into the Packages/C4BStudent directory.


## Installation of the Student's plug-in

To install the Student's plugin, try the quick installation method first.  If that does not work, try the manual installation method.

### Quick Installation of Student's plug-in

In Sublime Text 3, go to View -> Show Console

Copy this code to the console and hit enter:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4Brownies.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu", c4b_menu)
```


### Manual installation of Student's plug-in

- Download [C4BStudent.zip](https://github.com/vtphan/Code4Brownies/raw/master/C4BStudent.zip).
- Unzip the file into a directory called C4BStudent.
- Place this directory into the Sublime Text 3 Packages folder.  To open this Packages folder,
in Sublime Text, click on Preferences / Browse Packages.


