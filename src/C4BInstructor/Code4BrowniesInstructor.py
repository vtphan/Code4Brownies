#
# Author: Vinhthuy Phan, 2015
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket
import ntpath

SERVER_ADDR = "http://localhost:4030"
c4bi_BROADCAST_PATH = "broadcast"
c4bi_BROWNIE_PATH = "give_point"
c4bi_PEEK_PATH = "peek"
c4bi_POINTS_PATH = "points"
c4bi_REQUEST_ENTRY_PATH = "get_post"
c4bi_REQUEST_ENTRIES_PATH = "get_posts"
TIMEOUT = 10

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
	except urllib.error.HTTPError as err:
		sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
		sublime.message_dialog("{0}\nPossibly server not running or incorrect server address.".format(err))
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
			data = urllib.parse.urlencode({'content':content, 'ext':ext}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				sublime.status_message(response)


class c4biPointsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_POINTS_PATH)
		data = urllib.parse.urlencode({}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			users = {}
			for k,v in json_obj.items():
				if v['Uid'] not in users:
					users[v['Uid']] = 0
				users[v['Uid']] += v['Points']
			new_view = self.view.window().new_file()
			users = [ "%s,%s" % (k,v) for k,v in sorted(users.items()) ]
			new_view.insert(edit, 0, "\n".join(users))


class c4biGetAllCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_REQUEST_ENTRIES_PATH)
		data = urllib.parse.urlencode({}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			entries = json.loads(response)
			if entries:
				for entry in entries:
					# print(entry)
					ext = '' if entry['Ext']=='' else '.'+entry['Ext']
					userFile = os.path.join(POSTS_DIR, entry['Sid'] + ext)
					with open(userFile, 'w', encoding='utf-8') as fp:
						fp.write(entry['Body'])
					new_view = self.view.window().open_file(userFile)
			else:
				sublime.status_message("Queue is empty.")


class c4biPeekCommand(sublime_plugin.TextCommand):
	def request_entry(self, users, edit):
		def foo(selected):
			if selected < 0:
				return
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_REQUEST_ENTRY_PATH)
			data = urllib.parse.urlencode({'post':selected}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				json_obj = json.loads(response)
				ext = '' if json_obj['Ext']=='' else '.'+json_obj['Ext']
				userFile = os.path.join(POSTS_DIR, json_obj['Sid'] + ext)
				with open(userFile, 'w', encoding='utf-8') as fp:
					fp.write(json_obj['Body'])
				new_view = self.view.window().open_file(userFile)
		return foo

	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_PEEK_PATH)
		data = urllib.parse.urlencode({}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			if json_obj is None:
				sublime.status_message("Queue is empty.")
			else:
				users = [ '%s: %s, %s' % (entry['Uid'], entry['Pid'], entry['Sid']) for entry in json_obj ]
				if users:
					self.view.show_popup_menu(users, self.request_entry(users, edit))
				else:
					sublime.status_message("Queue is empty.")


class c4biAwardPointCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		this_file_name = self.view.file_name()
		if this_file_name:
			sid = this_file_name.split('.')[0]
			sid = ntpath.basename(sid)
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_BROWNIE_PATH)
			data = urllib.parse.urlencode({'sid':sid}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				# sublime.status_message(response)
				sublime.message_dialog(response)

class c4biAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BInstructor", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		sublime.message_dialog("Code4Brownies (v%s)\nCopyright © 2015-2016 Vinhthuy Phan" % version)


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


