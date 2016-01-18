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
c4b_SUBMIT_POST_PATH = "submit_post"
c4b_MY_POINTS_PATH = "my_points"
TIMEOUT = 10


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
	return json_obj


def c4bRequest(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.URLError:
		return "Server not running or incorrect server address."


class c4bShareCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4b_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4b_SUBMIT_POST_PATH)
		content = self.view.substr(sublime.Region(0, self.view.size()))
		values = {'login':os.getlogin(), 'uid':info['Name'],  'body':content}
		data = urllib.parse.urlencode(values).encode('ascii')
		response = c4bRequest(url,data)
		sublime.message_dialog(response)


class c4bShowPoints(sublime_plugin.WindowCommand):
	def run(self):
		info = c4b_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4b_MY_POINTS_PATH)
		values = {'login':os.getlogin(), 'uid':info['Name']}
		data = urllib.parse.urlencode(values).encode('ascii')
		response = c4bRequest(url,data)
		sublime.message_dialog(response)


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

