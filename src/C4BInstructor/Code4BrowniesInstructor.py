#
# Author: Vinhthuy Phan, 2015
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket

c4bi_BROADCAST_PATH = "broadcast"
c4bi_BROWNIE_PATH = "give_point"
c4bi_PEEK_PATH = "peek"
c4bi_POINTS_PATH = "points"
c4bi_REQUEST_ENTRY_PATH = "get_post"
c4bi_REQUEST_ENTRIES_PATH = "get_posts"
TIMEOUT = 10
ACTIVE_USERS = {}
SERVER_ADDR, PASSCODE = "", ""

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")
try:
	os.mkdir(POSTS_DIR)
except:
	pass

def c4biRequest(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError:
		sublime.message_dialog("HTTP error: possibly due to incorrect passcode.\n\nIf you don't know the passcode, restart the server and Sublime Text.")
	except urllib.error.URLError:
		sublime.message_dialog("Server not running or incorrect server address.")
	return None


class c4biBroadcastCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		content = self.view.substr(sublime.Region(0, self.view.size()))
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_BROADCAST_PATH)
		this_file_name = self.view.file_name()
		if this_file_name is not None:
			if '.' not in this_file_name:
				ext = ''
			else:
				ext = this_file_name.split('.')[-1]
			data = urllib.parse.urlencode({'passcode':PASSCODE, 'content':content, 'ext':ext}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				sublime.status_message(response)


class c4biPointsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_POINTS_PATH)
		data = urllib.parse.urlencode({'passcode':PASSCODE}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			new_view = self.view.window().new_file()
			users = [ "%s,%s" % (k,v) for k,v in json_obj.items() ]
			new_view.insert(edit, 0, "\n".join(users))


class c4biGetAllCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_REQUEST_ENTRIES_PATH)
		data = urllib.parse.urlencode({'passcode':PASSCODE}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			entries = json.loads(response)
			if entries:
				for entry in entries:
					ext = '' if entry['Ext']=='' else '.'+entry['Ext']
					userFile = os.path.join(POSTS_DIR, entry['Uid'] + ext)
					with open(userFile, 'w', encoding='utf-8') as fp:
						fp.write(entry['Body'])
					new_view = self.view.window().open_file(userFile)
					ACTIVE_USERS[new_view.id()] = entry['Uid']
			else:
				sublime.status_message("Queue is empty.")


class c4biPeekCommand(sublime_plugin.TextCommand):
	def request_entry(self, users, edit):
		def foo(selected):
			if selected < 0:
				return
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_REQUEST_ENTRY_PATH)
			data = urllib.parse.urlencode({'passcode':PASSCODE, 'post':selected}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				json_obj = json.loads(response)
				ext = '' if json_obj['Ext']=='' else '.'+json_obj['Ext']
				userFile = os.path.join(POSTS_DIR, users[selected] + ext)
				with open(userFile, 'w', encoding='utf-8') as fp:
					fp.write(json_obj['Body'])
				new_view = self.view.window().open_file(userFile)
				ACTIVE_USERS[new_view.id()] = users[selected]
		return foo

	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_PEEK_PATH)
		data = urllib.parse.urlencode({'passcode':PASSCODE}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			if json_obj is None:
				sublime.status_message("Queue is empty.")
			else:
				users = [ entry['Uid'] for entry in json_obj ]
				if users:
					self.view.show_popup_menu(users, self.request_entry(users, edit))
				else:
					sublime.status_message("Queue is empty.")


class c4biAwardPointCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		uid = ACTIVE_USERS.get(self.view.id())
		if uid is None:
			sublime.message_dialog("There is no user associated with this file.")
		else:
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_BROWNIE_PATH)
			data = urllib.parse.urlencode({'passcode':PASSCODE, 'uid':uid}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				sublime.status_message(response)

class c4biPasscode(sublime_plugin.WindowCommand):
	def run(self):
		def set_passcode(p):
			global PASSCODE
			PASSCODE = p
			sublime.message_dialog('Passcode is set.')

		if sublime.ok_cancel_dialog("Set a passcode only if the server is not running on this computer, or you don't want to use the default passcode.  In that case, first, run the server with a passcode. Then, use the same passcode here."):
			sublime.active_window().show_input_panel('Passcode','',set_passcode,None,None)

class c4biAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BInstructor", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		sublime.message_dialog("Code4Brownies (v%s)\nCopyright © 2015-2016 Vinhthuy Phan\n%s" %
			(version, SERVER_ADDR))


class c4biUpgrade(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Are you sure you want to upgrade Code4Brownies to the latest version?"):
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
				version = open(c4b_version).read().strip()
				sublime.message_dialog("Code4Brownies has been upgraded to version %s.  Latest server is at https://github.com/vtphan/Code4Brownies" % version)
			except:
				sublime.message_dialog("A problem occurred during upgrade.")

def Init():
	global PASSCODE, SERVER_ADDR
	s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	s.connect(("8.8.8.8", 80))
	SERVER_ADDR = "http://%s:4030 " % s.getsockname()[0]
	PASSCODE = s.getsockname()[0]

Init()

