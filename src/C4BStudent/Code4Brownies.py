# Code4Brownies - Student module
# Author: Vinhthuy Phan, 2015-2017
#
import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import threading
import time
import random
import datetime

c4b_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
c4b_CHECKIN_PATH = "checkin"
c4b_SHARE_PATH = "share"
c4b_MY_POINTS_PATH = "my_points"
c4b_RECEIVE_BROADCAST_PATH = "receive_broadcast"
c4b_CHECK_BROADCAST_PATH = "check_broadcast"
c4b_DEFAULT_FOLDER = os.path.join(os.path.expanduser('~'), 'C4B')
TIMEOUT = 7
CHECKED_IN = ''

c4b_WHITEBOARD_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Whiteboard")

Hints = {}
CUR_BID = None
# ------------------------------------------------------------------
def c4b_get_attr():
	try:
		with open(c4b_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		sublime.message_dialog("Please set server address and your name.")
		return None
	if 'Name' not in json_obj or len(json_obj['Name']) < 2:
		sublime.message_dialog("Please set your name.")
		return None
	if 'Server' not in json_obj or len(json_obj['Server']) < 4:
		sublime.message_dialog("Please set server address.")
		return None
	if 'Folder' not in json_obj:
		sublime.message_dialog("Please set course folder.")
		return None
	return json_obj

# ------------------------------------------------------------------
def c4bRequest(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError as err:
		sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
		sublime.message_dialog("{0}\nCannot connect to server.".format(err))
	return None

# ------------------------------------------------------------------
def get_hints(str):
	s = str.strip()
	if s=='':
		return []
	break_pattern = s.split('\n', 1)[0]
	hints = s.split(break_pattern)
	hints.pop(0)
	return hints

# ------------------------------------------------------------------
class c4bCheckinCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		global CHECKED_IN
		today = datetime.datetime.now().strftime('%d-%m-%Y')
		if CHECKED_IN == today:
			sublime.message_dialog("You are already checked in today.")
			return
		info = c4b_get_attr()
		if info is None:
			return
		data = urllib.parse.urlencode({'uid':info['Name']}).encode('utf-8')
		url = urllib.parse.urljoin(info['Server'], c4b_CHECKIN_PATH)
		response = c4bRequest(url, data)
		if response == 'Ok':
			sublime.message_dialog('Hi {}. You are now checked in.'.format(info['Name']))
			CHECKED_IN = today
		else:
			sublime.message_dialog('Fail to check in.')


# ------------------------------------------------------------------
class c4bMyBoardCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4b_get_attr()
		if info is None:
			return
		data = urllib.parse.urlencode({'uid':info['Name']}).encode('utf-8')
		url = urllib.parse.urljoin(info['Server'], c4b_RECEIVE_BROADCAST_PATH)
		response = c4bRequest(url, data)
		if response != None:
			json_obj = json.loads(response)
			if json_obj == []:
				sublime.message_dialog("Whiteboard is empty.")
			else:
				for board in json_obj:
					content = board['Content']
					ext = board['Ext']
					bid = board['Bid']
					Hints[bid] = [0, get_hints(board['HelpContent'])]
					if len(content.strip()) > 0:
						wb = os.path.join(info['Folder'], bid)
						wb += '.'+ext if ext!='' else ''
						with open(wb, 'w', encoding='utf-8') as f:
							f.write(content)
						new_view = sublime.active_window().open_file(wb)

# ------------------------------------------------------------------
class c4bHintCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		global Hints
		info = c4b_get_attr()
		if info is None:
			return
		bid = None
		file_name = self.view.file_name()
		if file_name is not None:
			basename = os.path.basename(file_name)
			prefix = basename.rsplit('.', 1)[0]
			if prefix in Hints:
				bid = prefix
		if bid in Hints:
			if Hints[bid][1] == []:
				sublime.message_dialog("There is no hint associated to this exercise.")
				return
			i = Hints[bid][0]
			if i >= len(Hints[bid][1]):
				sublime.message_dialog("No more hint.")
			else:
				help_content = Hints[bid][1][i]
				Hints[bid][0] = i+1
				# cwd = os.path.dirname(file_name)
				hint_file = os.path.join(info['Folder'], bid) + '.' + str(Hints[bid][0])
				with open(hint_file, 'w', encoding='utf-8') as f:
					f.write(help_content)
				new_view = sublime.active_window().open_file(hint_file)
		else:
			sublime.message_dialog("No hints associated with this file.")

# ------------------------------------------------------------------
class c4bShareCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4b_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4b_SHARE_PATH)

		# Guesstimate extension
		this_file_name = self.view.file_name()
		if this_file_name is None:
			return
		fname = this_file_name.rsplit('/',1)[-1]
		ext = 'py' if fname is None else fname.rsplit('.',1)[-1]
		bid = fname.split('.')[0]
		if not bid.startswith('wb_'):
			bid = ""
		header = ''
		if this_file_name is not None:
			lines = open(this_file_name, 'r', encoding='utf-8').readlines()
			if len(lines)>0 and (lines[0].startswith('#') or lines[0].startswith('//')):
				header = lines[0]

		# Determine content
		content = ''.join([ self.view.substr(s) for s in self.view.sel() ])
		if len(content) < 10:  # probably selected by mistake
			content = self.view.substr(sublime.Region(0, self.view.size()))
		else:
			content = header + '\n' + content

		# Now send
		hints_used = -1 if bid not in Hints else Hints[bid][0]
		values = {
			'uid':			info['Name'],
			'body':			content,
			'ext':			ext,
			'mode':			'code',
			'bid':			bid,
			'hints_used':	hints_used,
		}
		data = urllib.parse.urlencode(values).encode('utf-8')
		response = c4bRequest(url,data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4bAsk(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel(
			"Type your question. Press Enter.",
			"",
			self.send_question,
			None,
			None
		)

	def send_question(self, question):
		question = question.strip()
		if len(question) > 0:
			info = c4b_get_attr()
			if info is None:
				return
			url = urllib.parse.urljoin(info['Server'], c4b_SHARE_PATH)
			values = {
				'uid':info['Name'],
				'body':question,
				'ext':'',
				'mode': 'ask',
				'hints_used': -1,
			}
			data = urllib.parse.urlencode(values).encode('utf-8')
			response = c4bRequest(url,data)
			if response is not None:
				sublime.message_dialog(response)
		else:
			sublime.message_dialog("Question cannot be empty.")

# ------------------------------------------------------------------
class c4bVote(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel("ENTER to Vote or ESC to Cancel.",
			"",
			self.vote,
			None,
			None)

	def vote(self, answer):
		answer = answer.strip()
		if len(answer) > 0:
			info = c4b_get_attr()
			if info is None:
				return
			url = urllib.parse.urljoin(info['Server'], c4b_SHARE_PATH)
			values = {
				'uid':info['Name'],
				'body':answer,
				'ext':'',
				'mode': 'poll',
				'hints_used': -1,
			}
			data = urllib.parse.urlencode(values).encode('utf-8')
			response = c4bRequest(url,data)
			if response is not None:
				sublime.message_dialog(response)
		else:
			sublime.message_dialog("Answer cannot be empty.")

# ------------------------------------------------------------------
class c4bShowPoints(sublime_plugin.WindowCommand):
	def run(self):
		info = c4b_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4b_MY_POINTS_PATH)
		values = {'uid':info['Name']}
		# data = urllib.parse.urlencode(values).encode('ascii')
		data = urllib.parse.urlencode(values).encode('utf-8')
		response = c4bRequest(url,data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4bSetInfo(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4b_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()

		if 'Name' not in info:
			info['Name'] = ''
		if 'Server' not in info:
			info['Server'] = ''

		with open(c4b_FILE, 'w') as f:
			f.write(json.dumps(info, indent=4))

		sublime.active_window().open_file(c4b_FILE)

# ------------------------------------------------------------------
class c4bSetServer(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4b_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Server' not in info:
			info['Server'] = ''
		sublime.active_window().show_input_panel("Set server address.  Press Enter:",
			info['Server'],
			self.set,
			None,
			None)

	def set(self, addr):
		addr = addr.strip()
		if len(addr) > 0:
			try:
				with open(c4b_FILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			if not addr.startswith('http://'):
				addr = 'http://' + addr
			info['Server'] = addr
			with open(c4b_FILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
		else:
			sublime.message_dialog("Server address is empty.")

# ------------------------------------------------------------------
class c4bSetName(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4b_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Name' not in info:
			info['Name'] = ''
		sublime.active_window().show_input_panel("Set your name.  Press Enter:",
			info['Name'],
			self.set,
			None,
			None)

	def set(self, name):
		name = name.strip()
		if len(name) > 0:
			try:
				with open(c4b_FILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Name'] = name
			with open(c4b_FILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
		else:
			sublime.message_dialog("Server address cannot be empty.")

# ------------------------------------------------------------------
class c4bSetFolder(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4b_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Folder' not in info:
			info['Folder'] = c4b_DEFAULT_FOLDER
		sublime.active_window().show_input_panel("Set your course folder name.  Press Enter:",
			info['Folder'],
			self.set,
			None,
			None)

	def set(self, name):
		name = name.strip()
		if len(name) > 0:
			try:
				with open(c4b_FILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Folder'] = name
			if not os.path.exists(name):
				try:
					os.mkdir(name)
				except:
					sublime.message_dialog('Error creating directory {}, defaulting to {}.'.format(name, c4b_DEFAULT_FOLDER))
					info['Folder'] = c4b_DEFAULT_FOLDER
					os.mkdir(c4b_DEFAULT_FOLDER)

			with open(c4b_FILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
		else:
			sublime.message_dialog("Folder is empty.")

# ------------------------------------------------------------------
class c4bAbout(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BStudent", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		sublime.message_dialog("Code4Brownies (v%s)\nCopyright 2015-2017 Vinhthuy Phan" % version)

# ------------------------------------------------------------------
class c4bUpdate(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Are you sure you want to update Code4Brownies to the latest version?", "Yes"):
			package_path = os.path.join(sublime.packages_path(), "C4BStudent")
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			c4b_py = os.path.join(package_path, "Code4Brownies.py")
			c4b_menu = os.path.join(package_path, "Main.sublime-menu")
			c4b_version = os.path.join(package_path, "VERSION")
			try:
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Code4Brownies.py", c4b_py)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BStudent/Main.sublime-menu", c4b_menu)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
				version = open(c4b_version).read()
				sublime.message_dialog("Code4Brownies has been updated to version %s" % version)
			except:
				sublime.message_dialog("A problem occurred during update.")