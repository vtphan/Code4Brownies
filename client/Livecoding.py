import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json

INFO_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "LiveCodingInfo")
REMOTE_SHARE_HANDLER = "share"

class SetServerAddressCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Server address: ', '', self.set_server_address, None, None)

	def set_server_address(self, addr):
		if not addr.startswith("http"):
			addr = "http://" + addr
		try:
			with open(INFO_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			json_obj = json.loads('{}')
			json_obj['Username'] = 'Anonymous'

		json_obj['Address'] = addr
		with open(INFO_FILE, 'w') as f:
			f.write(json.dumps(json_obj))


class SetUsernameCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Username: ', '', self.set_username, None, None)

	def set_username(self, u):
		try:
			with open(INFO_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			json_obj = json.loads('{}')
		
		json_obj['Username'] = u
		with open(INFO_FILE, 'w') as f:
			f.write(json.dumps(json_obj))


class ShareCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		try:
			with open(INFO_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			sublime.message_dialog("Please set server address and username.")
			return
		if 'Address' not in json_obj:
			sublime.message_dialog("Please set server address.")
			return

		print(json_obj, json_obj['Address'], json_obj['Username'])
		
		url = urllib.parse.urljoin(json_obj['Address'], REMOTE_SHARE_HANDLER)
		content = self.view.substr(sublime.Region(0, self.view.size()))
		values = {'login':os.getlogin(), 'username':json_obj['Username'],  'body':content}
		data = urllib.parse.urlencode(values)
		data = data.encode('ascii')
		req = urllib.request.Request(url, data)
		try:
			with urllib.request.urlopen(req) as response:
				the_page = response.read()
				print(the_page)
				sublime.message_dialog("Got it!")
		except:
			sublime.message_dialog("Unable to share data. Make sure server address is correct.")


