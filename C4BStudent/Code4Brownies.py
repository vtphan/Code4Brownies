#
# Live Coding (student module)
# Author: Vinhthuy Phan, 2015
#
import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json

c4b_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
c4b_REGISTER_PATH = "register"
c4b_SUBMIT_POST_PATH = "submit_post"
c4b_MY_POINTS_PATH = "my_points"
TIMEOUT = 5

def c4b_set_attr(attr):
	def foo(value):
		try:
			with open(c4b_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			json_obj = json.loads('{}')
		json_obj[attr] = value
		if 'Address' in json_obj and not json_obj['Address'].startswith('http'):
			json_obj['Address'] = 'http://' + json_obj['Address']
		with open(c4b_FILE, 'w') as f:
			f.write(json.dumps(json_obj))

	return foo


def c4b_get_attr():
	try:
		with open(c4b_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		sublime.message_dialog("Please set server address and username.")
		return None
	if 'Address' not in json_obj:
		sublime.message_dialog("Please set server address.")
		return None
	if 'Username' not in json_obj:
		sublime.message_dialog("Please set username.")
		return None
	# if not json_obj['Address'].startswith('http'):
	# 	json_obj['Address'] = 'http://' + json_obj['Address']
	return json_obj


class c4bSetServerAddressCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Server address: ', '', c4b_set_attr('Address'), None, None)


class c4bSetUsernameCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Username: ', '', c4b_set_attr('Username'), None, None)



def c4bRequest(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.URLError:
		sublime.message_dialog("Server not running or incorrect server address.")
	return None


def c4b_register(value):
	uid, name = value.split(",")
	try:
		with open(c4b_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		json_obj = json.loads('{}')
	json_obj["Username"] = uid
	with open(c4b_FILE, 'w') as f:
		f.write(json.dumps(json_obj))

	if 'Address' not in json_obj:
		sublime.message_dialog("Please set server address first.")
		return None

	url = urllib.parse.urljoin(json_obj['Address'], c4b_REGISTER_PATH)
	print("URL",url, json_obj['Address'], json_obj)
	data = urllib.parse.urlencode({'uid':uid, 'name':name}).encode('ascii')
	response = c4bRequest(url,data)
	sublime.message_dialog(response)


class c4bRegisterCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Username: ', '', c4b_register, None, None)


class c4bShareCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4b_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Address'], c4b_SUBMIT_POST_PATH)
		content = self.view.substr(sublime.Region(0, self.view.size()))
		values = {'login':os.getlogin(), 'uid':info['Username'],  'body':content}
		data = urllib.parse.urlencode(values).encode('ascii')
		response = c4bRequest(url,data)
		if response == "1":
			sublime.message_dialog("Entry submitted succesfully.")
		elif response == "0":
			sublime.message_dialog("Invalid submission by a non-existent user.")


class c4bShowPoints(sublime_plugin.WindowCommand):
	def run(self):
		info = c4b_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Address'], c4b_MY_POINTS_PATH)
		values = {'login':os.getlogin(), 'username':info['Username']}
		data = urllib.parse.urlencode(values).encode('ascii')
		response = c4bRequest(url,data)
		sublime.message_dialog(response)



