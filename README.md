Code4Brownies
Author: Vinhthuy Phan, 2015

###Installation Sublime Text plugin

In Sublime Text 3, go to View -> Show Console

Copy this code to the console and hit enter:

```
from os.path import join; package_path = join(sublime.packages_path(), "C4BStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); c4b_py = join(package_path, "Code4Brownies.py") ; c4b_menu = join(package_path, "Main.sublime-menu"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Code4Brownies.py", c4b_py); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/C4BStudent/Main.sublime-menu", c4b_menu)
```
