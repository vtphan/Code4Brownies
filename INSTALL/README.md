# Installation of [Sublime Text 3](https://www.sublimetext.com/3)'s plug-in for **students**.

First, install [Sublime Text 3](https://www.sublimetext.com/3).

Then, open Sublime Text, go to View -> Show Console.  Then, copy this code to the console and hit enter:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4Brownies.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); c4b_version = os.path.join(package_path, "VERSION"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Main.sublime-menu", c4b_menu); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
```

After installation, students can share codes using the menu "ShareCode".



# Installation of [Sublime Text 3](https://www.sublimetext.com/3)'s plug-in for **instructor**

First, donwload and install [Sublime Text](https://www.sublimetext.com/3).

Then, go to View, open Console and execute the following code:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4BrowniesInstructor.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); c4b_version = os.path.join(package_path, "VERSION"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Code4BrowniesInstructor.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Main.sublime-menu", c4b_menu); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
```

Finally, donwload the server and run it on the instructor's machine.

- [Windows 64bit](https://github.com/vtphan/Code4Brownies/raw/master/INSTALL/c4b_windows_amd64)
- [Mac 64bit](https://github.com/vtphan/Code4Brownies/raw/master/INSTALL/c4b_darwin_amd64).
- Create a directory called "db" to store student records (in CSV format).



# Uninstall student's plugin

In Sublime Text, go to View, open Console.  Then, execute this code:

```
import os; import shutil; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); shutil.rmtree(package_path)
```

# Uninstall instructor's plugin

In Sublime Text, go to View, open Console.  Then, execute this code:

```
import os; import shutil; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); shutil.rmtree(package_path)
```


Copyright Vinhthuy Phan, 2015
