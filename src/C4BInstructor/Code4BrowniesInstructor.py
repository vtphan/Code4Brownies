# Code4Brownies - Instructor module
# Author: Vinhthuy Phan, 2015-2017
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket
import ntpath
import webbrowser

SERVER_ADDR = "http://localhost:4030"
c4bi_BROADCAST_PATH = "broadcast"
c4bi_TEST_DATA_PATH = "test_data"
c4bi_CLEAR_BOARD_PATH = "clear_board"
c4bi_BROWNIE_PATH = "give_points"
c4bi_PEEK_PATH = "peek"
c4bi_POINTS_PATH = "points"
c4bi_REQUEST_ENTRY_PATH = "get_post"
c4bi_REQUEST_ENTRIES_PATH = "get_posts"
c4bi_START_POLL_PATH = "start_poll"
c4bi_NEW_PROBLEM_PATH = "new_problem"
TIMEOUT = 10

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")
try:
	os.mkdir(POSTS_DIR)
except:
	pass

# ------------------------------------------------------------------
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

# ------------------------------------------------------------------
class c4biCleanCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Remove all student-submitted files."):
			files = [ f for f in os.listdir(POSTS_DIR) if f.startswith('c4b_') ]
			for f in files:
				local_file = os.path.join(POSTS_DIR, f)
				os.remove(local_file)
				sublime.status_message("remove " + local_file)

# ------------------------------------------------------------------
class c4biViewPollCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		webbrowser.open(SERVER_ADDR + "/poll")

# ------------------------------------------------------------------
class c4biBroadcastCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		content = self.view.substr(sublime.Region(0, self.view.size()))
		pid = content.split('\n',1)[0]
		file_name = self.view.file_name()
		if file_name is not None:
			ext = file_name.split('.')[-1] if '.' in file_name else ''
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_BROADCAST_PATH)
			data = urllib.parse.urlencode({'content':content, 'ext':ext, 'pid':pid}).encode('ascii')
			response = c4biRequest(url,data)
			if response is not None:
				sublime.status_message(response)
				local_dir = os.path.dirname(file_name)

				# For future usage: currently not supporting sending test data to server
				# self.send_test_data(local_dir)
			else:
				sublime.status_message("Unexpected response from server.")
		else:
			sublime.status_message("Must specify an existing file to broadcast.")

	def send_test_data(self, local_dir):
		def on_done(selected_idx):
			if selected_idx >= 0:
				selected_file = os.path.join(local_dir, files[selected_idx])
				try:
					with open(selected_file, 'r', encoding='utf-8') as f:
						content = f.read()
					url = urllib.parse.urljoin(SERVER_ADDR, c4bi_TEST_DATA_PATH)
					data = urllib.parse.urlencode({'content':content}).encode('ascii')
				except:
					sublime.status_message("Test data must be in ASCII format.")
				else:
					response = c4biRequest(url,data)
					if response is not None:
						sublime.status_message(response)
					else:
						sublime.status_message("Unexpected response from server.")
			else:
				sublime.status_message("No test data for this challenge.")

		files = os.listdir(local_dir)
		self.view.window().show_quick_panel(files, on_done)


# ------------------------------------------------------------------
# Instructor sends signal to clear whiteboard.
# ------------------------------------------------------------------
class c4biClearBoardCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_CLEAR_BOARD_PATH)
		data = urllib.parse.urlencode({}).encode('ascii')
		response = c4biRequest(url, data)
		if response is not None:
			sublime.status_message(response)

# ------------------------------------------------------------------
# Instructor retrieves all current and past points of all users.
# ------------------------------------------------------------------
class c4biPointsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_POINTS_PATH)
		data = urllib.parse.urlencode({}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			points, entries = {}, {}
			for k,v in json_obj.items():
				if v['Uid'] not in points:
					points[v['Uid']] = 0
					entries[v['Uid']] = 0
				points[v['Uid']] += v['Points']
				entries[v['Uid']] += 1
			new_view = self.view.window().new_file()
			users = [ "%s,%s,%s" % (k,entries[k],points[k]) for k,v in sorted(points.items()) ]
			new_view.insert(edit, 0, "\n".join(users))

# ------------------------------------------------------------------
# Instructor retrieves all posts.
# ------------------------------------------------------------------
class c4biGetAllCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_REQUEST_ENTRIES_PATH)
		data = urllib.parse.urlencode({}).encode('ascii')
		response = c4biRequest(url,data)
		if response is not None:
			entries = json.loads(response)
			# print(entries)
			if entries:
				for entry in reversed(entries):
					# print(entry)
					ext = '' if entry['Ext']=='' else '.'+entry['Ext']
					userFile = os.path.join(POSTS_DIR, entry['Sid'] + ext)
					with open(userFile, 'w', encoding='utf-8') as fp:
						fp.write(entry['Body'])
					new_view = self.view.window().open_file(userFile)
			else:
				sublime.status_message("Queue is empty.")

# ------------------------------------------------------------------
# Instructor starts poll mode
# ------------------------------------------------------------------
class c4biStartPollCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_START_POLL_PATH)
		data = urllib.parse.urlencode({}).encode('ascii')
		response = c4biRequest(url, data)
		if response == "true":
			sublime.message_dialog("A new poll has started.")
		elif response == "false":
			sublime.message_dialog("Poll is now closed.")
		else:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
# Instructor looks at new posts and is able to select one.
# ------------------------------------------------------------------
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
				# print(json_obj)
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

# ------------------------------------------------------------------
# Instructor rewards brownies.
# ------------------------------------------------------------------
class c4biAwardPoint1Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 1)

class c4biAwardPoint2Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 2)

class c4biAwardPoint3Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 3)

class c4biAwardPoint4Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 4)

class c4biAwardPoint5Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 5)

def award_points(self, edit, points):
	this_file_name = self.view.file_name()
	if this_file_name:
		sid = this_file_name.split('.')[0]
		sid = ntpath.basename(sid)
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_BROWNIE_PATH)
		data = urllib.parse.urlencode({'sid':sid, 'points':points}).encode('ascii')
		response = c4biRequest(url,data)
		if response:
			sublime.status_message(response)
			self.view.window().run_command('close')
		else:
			sublime.status_message("no uid associated with this file.")

# ------------------------------------------------------------------
class c4biAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BInstructor", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		sublime.message_dialog("Code4Brownies (v%s)\nCopyright © 2015-2016 Vinhthuy Phan" % version)

# ------------------------------------------------------------------
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


