#
# Author: Vinhthuy Phan, 2015
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket

c4bi_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
c4bi_BROADCAST_PATH = "broadcast"
c4bi_BROWNIE_PATH = "give_point"
c4bi_PEEK_PATH = "peek"
c4bi_POINTS_PATH = "points"
c4bi_REQUEST_ENTRY_PATH = "get_post"
c4bi_REQUEST_ENTRIES_PATH = "get_posts"
TIMEOUT = 10
ACTIVE_USERS = {}

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")
try:
	os.mkdir(POSTS_DIR)
except:
	pass
def c4bi_get_attr():
	try:
		with open(c4bi_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		sublime.message_dialog("Please set information first.")
		return None
	if 'Server' not in json_obj or 'Passcode' not in json_obj:
		sublime.message_dialog("Please set information completely.")
		return None
	if not json_obj['Server'].startswith("http://"):
		sublime.message_dialog("Server must starts with http://\nReset information.")
		return None
	return json_obj


def c4biRequest(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError:
		sublime.message_dialog("HTTP error: possibly due to incorrect passcode.")
	except urllib.error.URLError:
		sublime.message_dialog("Server not running or incorrect server address.")
	return None


class c4biBroadcastCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		content = self.view.substr(sublime.Region(0, self.view.size()))
		url = urllib.parse.urljoin(info['Server'], c4bi_BROADCAST_PATH)
		this_file_name = self.view.file_name()
		if this_file_name is not None:
			if '.' not in this_file_name:
				ext = ''
			else:
				ext = this_file_name.split('.')[-1]
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'content':content, 'ext':ext}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				sublime.status_message(response)


class c4biPointsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		url = urllib.parse.urljoin(info['Server'], c4bi_POINTS_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			new_view = self.view.window().new_file()
			users = [ "%s,%s" % (k,v) for k,v in json_obj.items() ]
			new_view.insert(edit, 0, "\n".join(users))


class c4biGetAllCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		url = urllib.parse.urljoin(info['Server'], c4bi_REQUEST_ENTRIES_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			entries = json.loads(response)
			if entries:
				for entry in entries:
					ext = '' if entry['Ext']=='' else '.'+entry['Ext']
					userFile = os.path.join(POSTS_DIR, entry['Uid'] + ext)
					with open(userFile, 'w') as fp:
						fp.write(entry['Body'])
					new_view = self.view.window().open_file(userFile)
					ACTIVE_USERS[new_view.id()] = entry['Uid']
			else:
				sublime.status_message("Queue is empty.")


class c4biGetCommand(sublime_plugin.TextCommand):
	def request_entry(self, info, users, edit):
		def foo(selected):
			if selected < 0:
				return
			url = urllib.parse.urljoin(info['Server'], c4bi_REQUEST_ENTRY_PATH)
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'post':selected}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				json_obj = json.loads(response)
				ext = '' if json_obj['Ext']=='' else '.'+json_obj['Ext']
				userFile = os.path.join(POSTS_DIR, users[selected] + ext)
				with open(userFile, 'w') as fp:
					fp.write(json_obj['Body'])
				new_view = self.view.window().open_file(userFile)
				ACTIVE_USERS[new_view.id()] = users[selected]
		return foo

	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		url = urllib.parse.urljoin(info['Server'], c4bi_PEEK_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			if json_obj is None:
				sublime.status_message("Queue is empty.")
			else:
				users = [ entry['Uid'] for entry in json_obj ]
				if users:
					self.view.show_popup_menu(users, self.request_entry(info, users, edit))
				else:
					sublime.status_message("Queue is empty.")


class c4biAwardPointCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return
		uid = ACTIVE_USERS.get(self.view.id())
		if uid is None:
			sublime.message_dialog("There is no user associated with this file.")
		else:
			url = urllib.parse.urljoin(info['Server'], c4bi_BROWNIE_PATH)
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'uid':uid}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				sublime.status_message(response)


class c4biAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BInstructor", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		addr = socket.gethostbyname(socket.gethostname()) + ":4030"
		sublime.message_dialog("Code4Brownies (v%s)\nServer address: %s\n\nCopyright © 2015-2016 Vinhthuy Phan." %
			(version,addr))


class c4biSetInfo(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4bi_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'Passcode' not in info:
			info['Passcode'] = 'password'
		if 'Server' not in info:
			info['Server'] = 'http://0.0.0.0:4030'

		with open(c4bi_FILE, 'w') as f:
			f.write(json.dumps(info, indent=4))

		sublime.active_window().open_file(c4bi_FILE)


class c4biUpgrade(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Are you sure you want to upgrade Code4Brownies to the latest version?", "Yes"):
			package_path = os.path.join(sublime.packages_path(), "C4BInstructor");
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			c4b_py = os.path.join(package_path, "Code4BrowniesInstructor.py")
			c4b_menu = os.path.join(package_path, "Main.sublime-menu")
			c4b_version = os.path.join(package_path, "VERSION")
			try:
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Code4BrowniesInstructor.py", c4b_py)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Main.sublime-menu", c4b_menu)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
				version = open(c4b_version).read()
				sublime.message_dialog("Code4Brownies has been upgraded to version %s" % version)
			except:
				sublime.message_dialog("A problem occurred during upgrade.")

