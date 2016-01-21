Code4Brownies
Author: Vinhthuy Phan, 2015

### Manual installation

- Download [C4BStudent.zip](https://github.com/vtphan/Code4Brownies/raw/master/C4BStudent.zip).
- Unzip the file into a directory called C4BStudent.
- Place this directory into the Sublime Text 3 Packages folder.  To open this Packages folder, 
in Sublime Text, click on Preferences / Browse Packages.


### Quick Installation (may not always work)

In Sublime Text 3, go to View -> Show Console

Copy this code to the console and hit enter:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4Brownies.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu", c4b_menu)
```
