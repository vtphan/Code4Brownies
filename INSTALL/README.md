# Overview

First, install [Sublime Text 3](https://www.sublimetext.com/3).  Then, students can install the Student's plug in.  After installation, there is a dedicated menu for students on the ST3 menu bar.  That's it.

For instructors, install the server and Instructor's plug-in.  After installation, there is a dedicated menu for instructors on the ST3 menu bar.


# Installation/Update for Students

In Sublime Text 3, go to View -> Show Console.  Then, copy this code to the console and hit enter:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4Brownies.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Main.sublime-menu", c4b_menu)
```

To uninstall the plugin, in ST3 console, copy, paste and execute this code:

```
import os; import shutil; package_path = os.path.join(sublime.packages_path(), "C4BStudent"); shutil.rmtree(package_path)
```


# Installation/Update for Instructor

##### Donwload the server and run it on the instructor's machine.

- [Windows 64bit](https://github.com/vtphan/Code4Brownies/raw/master/INSTALL/c4b_windows_amd64)
- [Mac 64bit](https://github.com/vtphan/Code4Brownies/raw/master/INSTALL/c4b_darwin_amd64).
- Create a directory called "db" to store student records (in CSV format).


#####  Install/update of Instructor's plug in

Open Console and execute the following code:

```
import os; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = os.path.join(package_path, "Code4BrowniesInstructor.py") ; c4b_menu = os.path.join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Code4BrowniesInstructor.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Main.sublime-menu", c4b_menu)
```

To uninstall the plugin, in ST3 console, copy, paste and execute this code:

```
import os; import shutil; package_path = os.path.join(sublime.packages_path(), "C4BInstructor"); shutil.rmtree(package_path)
```

# Manual Installation/Update of Student's plug-in

If the Installation/Update method above does not work, do this:

- Download [C4BStudent.zip](https://github.com/vtphan/Code4Brownies/raw/master/INSTALL/C4BStudent.zip).
- Unzip the file into a directory called C4BStudent.
- Place this directory into the Sublime Text 3 Packages folder.  To open this Packages folder,
in Sublime Text, click on Preferences / Browse Packages.



# Manual Installation/Update of Instructor's plug-in

If the Installation/Update method above does not work, do this:

 - Download [C4BStudentInstructor.zip](https://github.com/vtphan/Code4Brownies/raw/master/INSTALL/C4BInstructor.zip).
 - Unzip the file into a directory called C4BInstructor.
 - Place this directory into the Sublime Text 3 Packages folder.  To open this Packages folder,
in Sublime Text, click on Preferences / Browse Packages.


Copyright Vinhthuy Phan, 2015
