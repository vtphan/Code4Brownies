import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json

LC_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "LiveCodingInfo")
REMOTE_SHARE_HANDLER = "share"

def lc_set_attr(attr):
	def foo(value):
		try:
			with open(LC_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			json_obj = json.loads('{}')
		json_obj[attr] = value
		with open(LC_FILE, 'w') as f:
			f.write(json.dumps(json_obj))

	return foo


def lc_get_attr():
	try:
		with open(LC_FILE, 'r') as f:
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
	if not json_obj['Address'].startswith('http'):
		json_obj['Address'] = 'http://' + json_obj['Address']
	return json_obj


class LcSetServerAddressCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Server address: ', '', lc_set_attr('Address'), None, None)


class LcSetUsernameCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Username: ', '', lc_set_attr('Username'), None, None)


class LcShareCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = lc_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Address'], REMOTE_SHARE_HANDLER)
		content = self.view.substr(sublime.Region(0, self.view.size()))
		values = {'login':os.getlogin(), 'username':info['Username'],  'body':content}
		data = urllib.parse.urlencode(values)
		data = data.encode('ascii')
		req = urllib.request.Request(url, data)
		try:
			with urllib.request.urlopen(req) as response:
				the_page = response.read()
				print(the_page)
				sublime.message_dialog("Got it!")
		except urllib.error.URLError:
			sublime.message_dialog("URL Error: reset server address.")


