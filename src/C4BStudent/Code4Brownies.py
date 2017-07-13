# Code4Brownies - Student module
# Author: Vinhthuy Phan, 2015-2017
#
import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json

c4b_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
c4b_SUBMIT_POST_PATH = "submit_post"
c4b_MY_POINTS_PATH = "my_points"
c4b_RECEIVE_BROADCAST_PATH = "receive_broadcast"
# c4b_RECEIVE_FEEDBACK_PATH = "receive_feedback"

TIMEOUT = 10

# ------------------------------------------------------------------
def c4b_get_attr():
	try:
		with open(c4b_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		sublime.message_dialog("Please set information first.")
		return None
	if 'Server' not in json_obj or 'Name' not in json_obj:
		sublime.message_dialog("Please set information completely.")
		return None
	if not json_obj['Server'].startswith("http://"):
		sublime.message_dialog("Server must starts with http://\nReset information.")
		return None
	return json_obj

# ------------------------------------------------------------------
def c4bRequest(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.URLError as err:
		sublime.message_dialog("{0}\nPossibly server not running or incorrect server address.".format(err))
		return None

# ------------------------------------------------------------------
def _receive_broadcast(self, edit, uid=''):
	info = c4b_get_attr()
	if info is None:
		return
	data = urllib.parse.urlencode({'uid':uid}).encode('ascii')
	url = urllib.parse.urljoin(info['Server'], c4b_RECEIVE_BROADCAST_PATH)
	response = c4bRequest(url, data)
	if response is not None:
		json_obj = json.loads(response)
		content = json_obj['content']
		if len(content.strip()) > 0:
			new_view = self.view.window().new_file()
			new_view.insert(edit, 0, content)
		else:
			if uid=='':
				sublime.message_dialog("Whiteboard is empty.")
			else:
				sublime.message_dialog("Your board is empty.")

# ------------------------------------------------------------------
class c4bReceivebroadcastCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		_receive_broadcast(self, edit)

# ------------------------------------------------------------------
class c4bMyBoardCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4b_get_attr()
		if info is None:
			return
		_receive_broadcast(self, edit, info['Name'])

# ------------------------------------------------------------------
# class c4bGetFeedbackCommand(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		info = c4b_get_attr()
# 		if info is None:
# 			return
# 		data = urllib.parse.urlencode({'uid':info['Name']}).encode('ascii')
# 		url = urllib.parse.urljoin(info['Server'], c4b_RECEIVE_FEEDBACK_PATH)
# 		response = c4bRequest(url, data)
# 		if response is not None:
# 			json_obj = json.loads(response)
# 			content = json_obj['content']
# 			if len(content.strip()) > 0:
# 				new_view = self.view.window().new_file()
# 				new_view.insert(edit, 0, content)
# 			else:
# 				sublime.message_dialog("You have no feedback.")

# ------------------------------------------------------------------
class c4bShareCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		this_file_name = self.view.file_name()
		if this_file_name is not None:
			if '.' not in this_file_name:
				ext = ''
			else:
				ext = this_file_name.split('.')[-1]
			info = c4b_get_attr()
			if info is None:
				return
			url = urllib.parse.urljoin(info['Server'], c4b_SUBMIT_POST_PATH)
			content = self.view.substr(sublime.Region(0, self.view.size()))
			values = {'uid':info['Name'], 'body':content, 'ext':ext}
			data = urllib.parse.urlencode(values).encode('ascii')
			response = c4bRequest(url,data)
			if response is not None:
				sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4bShowPoints(sublime_plugin.WindowCommand):
	def run(self):
		info = c4b_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4b_MY_POINTS_PATH)
		values = {'uid':info['Name']}
		data = urllib.parse.urlencode(values).encode('ascii')
		response = c4bRequest(url,data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4bSetInfo(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4b_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'Name' not in info:
			info['Name'] = 'JohnSmith'
		if 'Server' not in info:
			info['Server'] = 'http://0.0.0.0:4030'

		with open(c4b_FILE, 'w') as f:
			f.write(json.dumps(info, indent=4))

		sublime.active_window().open_file(c4b_FILE)

# ------------------------------------------------------------------
class c4bAbout(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BStudent", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		sublime.message_dialog("Code4Brownies (v%s)\nCopyright 2015-2016 Vinhthuy Phan" % version)

# ------------------------------------------------------------------
class c4bUpgrade(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Are you sure you want to upgrade Code4Brownies to the latest version?", "Yes"):
			package_path = os.path.join(sublime.packages_path(), "C4BStudent")
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			c4b_py = os.path.join(package_path, "Code4Brownies.py")
			c4b_menu = os.path.join(package_path, "Main.sublime-menu")
			c4b_version = os.path.join(package_path, "VERSION")
			try:
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Code4Brownies.py", c4b_py)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Main.sublime-menu", c4b_menu)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
				version = open(c4b_version).read()
				sublime.message_dialog("Code4Brownies has been upgraded to version %s" % version)
			except:
				sublime.message_dialog("A problem occurred during upgrade.")