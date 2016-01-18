#
# Author: Vinhthuy Phan, 2015
#
import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket

FILE_EXTENSION = ".py"

c4bi_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
c4bi_BROWNIE_PATH = "give_point"
c4bi_ENTRIES_PATH = "posts"
c4bi_POINTS_PATH = "points"
c4bi_REQUEST_ENTRY_PATH = "get_post"
c4bi_CHECK_PATH = "check_post"
TIMEOUT = 10
TRACKING = False
TRACKING_INTERVAL = 30000
ACTIVE_USERS = {}

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")
try:
	os.mkdir(POSTS_DIR)
except:
	pass


def c4bi_check_for_posts():
	print("Tracking", TRACKING)
	if TRACKING:
		info = c4bi_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4bi_CHECK_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = c4biRequest(url,data)
		print(response)
		if response == 'yes':
			os.system('afplay /System/Library/Sounds/Glass.aiff')
	sublime.set_timeout_async(c4bi_check_for_posts, TRACKING_INTERVAL)

sublime.set_timeout_async(c4bi_check_for_posts, TRACKING_INTERVAL)

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

class c4biTrackCode(sublime_plugin.WindowCommand):
	def run(self):
		global TRACKING
		TRACKING = True

class c4biUntrackCode(sublime_plugin.WindowCommand):
	def run(self):
		global TRACKING
		TRACKING = False

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


class c4biGetCommand(sublime_plugin.TextCommand):
	def request_entry(self, info, users, edit):
		def foo(selected):
			if selected < 0:
				return
			url = urllib.parse.urljoin(info['Server'], c4bi_REQUEST_ENTRY_PATH)
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'post':selected}).encode('ascii')
			response = c4biRequest(url,data)
			json_obj = json.loads(response)
			userFile = os.path.join(POSTS_DIR, users[selected] + FILE_EXTENSION)
			with open(userFile, 'w') as fp:
				fp.write(json_obj['Body'])
			new_view = self.view.window().open_file(userFile)
			ACTIVE_USERS[new_view.id()] = users[selected]
		return foo

	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		url = urllib.parse.urljoin(info['Server'], c4bi_ENTRIES_PATH)
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
		elif sublime.ok_cancel_dialog("Award 1 point to "+uid+"?"):
			url = urllib.parse.urljoin(info['Server'], c4bi_BROWNIE_PATH)
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'uid':uid}).encode('ascii')
			response = c4biRequest(url,data)


class c4biAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		addr = socket.gethostbyname(socket.gethostname()) + ":4030"
		sublime.message_dialog("Code4Brownies\nServer address: %s\n\nCopyright © 2015 Vinhthuy Phan." %
			addr)


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