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
c4bi_REGISTERED_USERS_PATH = "registered_users"
c4bi_APPROVE_PATH = "approve"
TIMEOUT = 5
ACTIVE_USERS = {}

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")
try:
	os.mkdir(POSTS_DIR)
except:
	pass

def c4bi_set_attr(attr):
	def foo(value):
		try:
			with open(c4bi_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			json_obj = json.loads('{}')
		json_obj[attr] = value
		if 'Address' in json_obj and not json_obj['Address'].startswith('http'):
			json_obj['Address'] = 'http://' + json_obj['Address']
		with open(c4bi_FILE, 'w') as f:
			f.write(json.dumps(json_obj))
	return foo


def c4bi_get_attr():
	try:
		with open(c4bi_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		sublime.message_dialog("Please set server address and passcode.")
		return None
	if 'Address' not in json_obj:
		sublime.message_dialog("Please set server address.")
		return None
	if 'Passcode' not in json_obj:
		sublime.message_dialog("Please set passcode.")
		return None
	if not json_obj['Address'].startswith('http'):
		json_obj['Address'] = 'http://' + json_obj['Address']
	return json_obj


class c4biSetServerAddressCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Server address: ', '', c4bi_set_attr('Address'), None, None)


class c4biSetPasscodeCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Passcode: ', '', c4bi_set_attr('Passcode'), None, None)


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


class c4biRegisteredUsersCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		url = urllib.parse.urljoin(info['Address'], c4bi_REGISTERED_USERS_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = c4biRequest(url,data)
		json_obj = json.loads(response)
		new_view = self.view.window().new_file()
		users = [ "%s,%s" % (k,v) for k,v in json_obj.items() ]
		new_view.insert(edit, 0, "\n".join(users))


class c4biPointsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		url = urllib.parse.urljoin(info['Address'], c4bi_POINTS_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = c4biRequest(url,data)
		json_obj = json.loads(response)
		new_view = self.view.window().new_file()
		users = [ "%s,%s" % (k,v) for k,v in json_obj.items() ]
		new_view.insert(edit, 0, "\n".join(users))


class c4biApproveRegistrationCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4bi_get_attr()
		if info is None:
			return

		url = urllib.parse.urljoin(info['Address'], c4bi_APPROVE_PATH)
		text = self.view.substr(sublime.Region(0, self.view.size()))
		approved = {}
		for line in text.splitlines():
			if line != '':
				if line.count(",") != 1:
					sublime.message_dialog("Aborted due to an invalid entry: " + line)
					return
				uid, name = line.split(",")
				approved[uid] = name
		approved = json.dumps(approved)		
		data = urllib.parse.urlencode({'passcode':info['Passcode'], 'approved':approved}).encode('ascii')
		response = c4biRequest(url,data)
		sublime.message_dialog("Done.")



class c4biGetCommand(sublime_plugin.TextCommand):
	def request_entry(self, info, users, edit):
		def foo(selected):
			if selected < 0:
				return
			url = urllib.parse.urljoin(info['Address'], c4bi_REQUEST_ENTRY_PATH)
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

		url = urllib.parse.urljoin(info['Address'], c4bi_ENTRIES_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = c4biRequest(url,data)
		if response is None:
			return
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
			url = urllib.parse.urljoin(info['Address'], c4bi_BROWNIE_PATH)
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'uid':uid}).encode('ascii')
			response = c4biRequest(url,data)
			print(response)


class c4biAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		addr = socket.gethostbyname(socket.gethostname()) + ":4030"
		sublime.message_dialog("Code4Brownies\nServer address: %s\n\nCopyright © 2015 Vinhthuy Phan." %
			addr)
