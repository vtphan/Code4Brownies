import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json

SERVER = "http://localhost:4030"
DEQUE_PATH = "deque"
CUR_ENTRY_PATH = "currentEntry"
BROWNIE_PATH = "brownie"

class GetCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		global CUR_USER
		url = urllib.parse.urljoin(SERVER, DEQUE_PATH)
		req = urllib.request.Request(url)
		with urllib.request.urlopen(req) as response:
			json_obj = json.loads(response.read().decode(encoding="utf-8"))
			body, N = json_obj['Body'], json_obj['N']
			self.view.replace(edit, sublime.Region(0, self.view.size()), body)
			sublime.status_message(str(N) + " entries left")

class AwardPointCommand(sublime_plugin.WindowCommand):
	def run(self):
		url = urllib.parse.urljoin(SERVER, CUR_ENTRY_PATH)
		req = urllib.request.Request(url)
		with urllib.request.urlopen(req) as response:
			json_obj = json.loads(response.read().decode(encoding="utf-8"))
		if json_obj is not None:
			if sublime.ok_cancel_dialog("Give a brownie point to "+json_obj['User']) == True:
				url = urllib.parse.urljoin(SERVER, BROWNIE_PATH)
				req = urllib.request.Request(url)
				with urllib.request.urlopen(req) as response:
					print(response.read())
		else:
			print("No user has been dequed.")
