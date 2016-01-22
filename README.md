This package includes a Sublime-Text 3 plugin, which allows students share their codes with an instructor
in real time.

## Update of the Student's plug-in

Assuming that you have properly installed the plug-in.  The quick way is to Show Console (in View)
paste and run this code:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4Brownies.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu", c4b_menu)
```

If this gives an error (e.g. "os not found"), then save [Code4Brownies.py](https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py) and [Main.sublime-menu](https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu) into your Packages/C4BStudent directory.  They will replace the old files.  Caution: when you save, the browser may outsmart itself by adding a ".txt" extension to Main.sublime-menu. Do not let the browser do this. Otherwise, you must remove the ".txt" extension.


# Installation of the Student's plug-in

To install the Student's plugin, try the quick installation method first.  If that does not work, try the manual installation method.

### Quick Installation of Student's plug-in

In Sublime Text 3, go to View -> Show Console

Copy this code to the console and hit enter:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4Brownies.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu", c4b_menu)
```


### Manual installation of Student's plug-in

- Download [C4BStudent.zip](https://github.com/vtphan/Code4Brownies/raw/master/downloads/C4BStudent.zip).
- Unzip the file into a directory called C4BStudent.
- Place this directory into the Sublime Text 3 Packages folder.  To open this Packages folder,
in Sublime Text, click on Preferences / Browse Packages.


# Installation of Instructor's server and plug-in

- Download the server: [Windows 64bit](https://github.com/vtphan/Code4Brownies/raw/master/downloads/c4b_windows_amd64) or [Mac 64bit](https://github.com/vtphan/Code4Brownies/raw/master/downloads/c4b_darwin_amd64).

- Quick install of Sublime Text 3 plug in: (a) open Console and execute the following code:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4BrowniesInstructor.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BInstructor/Code4BrowniesInstructor.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BInstructor/Main.sublime-menu", c4b_menu)
```

- If the quick install method does not work, you can install the plug in manually:

    - Download [C4BStudentInstructor.zip](https://github.com/vtphan/Code4Brownies/raw/master/downloads/C4BInstructor.zip).
    - Unzip the file into a directory called C4BInstructor.
    - Place this directory into the Sublime Text 3 Packages folder.  To open this Packages folder,
in Sublime Text, click on Preferences / Browse Packages.

Copyright Vinhthuy Phan, 2015
