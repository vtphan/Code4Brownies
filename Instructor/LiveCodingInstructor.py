import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json

LCI_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "LiveCodingInstructorInfo")
LCI_BROWNIE_PATH = "brownie"
LCI_ENTRIES_PATH = "entries"
LCI_REQUEST_ENTRY_PATH = "request_entry"
CURRENT_USER = None

def lci_set_attr(attr):
	def foo(value):
		try:
			with open(LCI_FILE, 'r') as f:
				json_obj = json.loads(f.read())
		except:
			json_obj = json.loads('{}')
		json_obj[attr] = value
		with open(LCI_FILE, 'w') as f:
			f.write(json.dumps(json_obj))

	return foo


def lci_get_attr():
	try:
		with open(LCI_FILE, 'r') as f:
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


class LciSetServerAddressCommand(sublime_plugin.WindowCommand):
	def run(self):
		info = lci_get_attr()
		sublime.active_window().show_input_panel('Server address: ', info['Address'], lci_set_attr('Address'), None, None)


class LciSetPasscodeCommand(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel('Passcode: ', '', lci_set_attr('Passcode'), None, None)


def LciRequest(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError:
		sublime.message_dialog("HTTP error: possibly due to incorrect passcode.")
	except urllib.error.URLError:
		sublime.message_dialog("URL error: reset server address.")


class LciGetCommand(sublime_plugin.TextCommand):
	def request_entry(self, info, users, edit):
		def foo(selected):
			global CURRENT_USER, CURRENT_USER_NAME
			if selected < 0:
				return
			url = urllib.parse.urljoin(info['Address'], LCI_REQUEST_ENTRY_PATH)
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'entry':selected}).encode('ascii')
			response = LciRequest(url,data)
			json_obj = json.loads(response)
			self.view.replace(edit, sublime.Region(0, self.view.size()), json_obj['Body'])
			CURRENT_USER = users[selected]
		return foo

	def run(self, edit):
		info = lci_get_attr()
		if info is None:
			return
		
		url = urllib.parse.urljoin(info['Address'], LCI_ENTRIES_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('ascii')
		response = LciRequest(url,data)
		json_obj = json.loads(response)
		users = [ entry['User'] for entry in json_obj ]
		if users:
			self.view.show_popup_menu(users, self.request_entry(info, users, edit))
		else:
			sublime.message_dialog("No entry.")


class LciAwardPointCommand(sublime_plugin.WindowCommand):
	def run(self):
		global CURRENT_USER
		if CURRENT_USER is None:
			return
		info = lci_get_attr()
		if info is None:
			return
		if sublime.ok_cancel_dialog("Award 1 point to "+CURRENT_USER+"?"):
			url = urllib.parse.urljoin(info['Address'], LCI_BROWNIE_PATH)
			data = urllib.parse.urlencode({'passcode':info['Passcode'], 'user':CURRENT_USER}).encode('ascii')
			response = LciRequest(url,data)
			print(response)

