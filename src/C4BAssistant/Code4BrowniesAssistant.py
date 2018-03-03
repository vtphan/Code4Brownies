# Code4Brownies - TA module
# Author: Vinhthuy Phan, 2015-2018
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket
import webbrowser

c4ba_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
c4ba_FEEDBACK_PATH = "ta_feedback"
c4ba_BROWNIE_PATH = "ta_give_points"
c4ba_REQUEST_ENTRIES_PATH = "ta_get_posts"
TIMEOUT = 7

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")

# ------------------------------------------------------------------

# def check_new_submissions():
# 	delay = 5000
# 	info = c4ba_get_attr(verbose=False)
# 	if info is None:
# 		return
# 	url = urllib.parse.urljoin(info['Server'], c4ba_QUEUE_LENGTH_PATH)
# 	data = urllib.parse.urlencode({}).encode('utf-8')
# 	response = c4baRequest(url, data, verbose=False)
# 	if response is not None:
# 		count = len(response)
# 		if count > 0:
# 			print('There are {} submissions in the queue.'.format(count))
# 			sublime.status_message('There are {} submissions in the queue.'.format(count))
# 	else:
# 		print('Error checking for new submissions. Response is None')
# 	sublime.set_timeout_async(check_new_submissions, delay)

# sublime.set_timeout_async(check_new_submissions, 5000)

# ------------------------------------------------------------------
def c4ba_get_attr(verbose=True):
	try:
		with open(c4ba_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		if verbose:
			sublime.message_dialog("Please set server address and your name.")
		return None
	if 'Server' not in json_obj or len(json_obj['Server']) < 4:
		if verbose:
			sublime.message_dialog("Please set server address.")
		return None
	if 'Passcode' not in json_obj or len(json_obj['Passcode']) < 4:
		if verbose:
			sublime.message_dialog("Please set passcode.")
		return None
	return json_obj

# ------------------------------------------------------------------
def c4baRequest(url, data, headers={}, verbose=True):
	req = urllib.request.Request(url, data, headers=headers)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError as err:
		if verbose:
			sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
		if verbose:
			sublime.message_dialog("{0}\nCannot connect to server.".format(err))
	return None

# ------------------------------------------------------------------
class c4baCleanCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Remove all student-submitted files."):
			if os.path.isdir(POSTS_DIR):
				files = [ f for f in os.listdir(POSTS_DIR) if f.startswith('c4b_') ]
				for f in files:
					local_file = os.path.join(POSTS_DIR, f)
					os.remove(local_file)
					sublime.status_message("remove " + local_file)

# ------------------------------------------------------------------
# TA gives feedback on this specific file
# ------------------------------------------------------------------
class c4baGiveFeedbackCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4ba_get_attr()
		if info is None:
			return

		this_file_name = self.view.file_name()
		if this_file_name is None:
			sublime.message_dialog('No student is associated with this file.')
			return

		# Get content and ext
		ext = this_file_name.rsplit('.',1)[-1]
		header = ''
		lines = open(this_file_name, 'r', encoding='utf-8').readlines()
		if len(lines)>0 and (lines[0].startswith('#') or lines[0].startswith('//')):
			header = lines[0]
		content = ''.join(lines)

		# Determine sid
		basename = os.path.basename(this_file_name)
		if not basename.startswith('c4b_'):
			sublime.message_dialog('No student is associated with this file.')
			return
		sid = basename.split('.')[0]
		sid = sid.split('c4b_')[1]

		data = urllib.parse.urlencode({
			'content': 		content,
			'sid':			sid,
			'ext': 			ext,
			'passcode':		info['Passcode'],
		}).encode('utf-8')

		url = urllib.parse.urljoin(info['Server'], c4ba_FEEDBACK_PATH)
		response = c4baRequest(url, data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
# TA retrieves all posts.
# ------------------------------------------------------------------
class c4baGetAllCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4ba_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4ba_REQUEST_ENTRIES_PATH)
		data = urllib.parse.urlencode({'passcode':info['Passcode']}).encode('utf-8')
		response = c4baRequest(url,data)
		if response is not None:
			entries = json.loads(response)
			# print(entries)
			if entries:
				for entry in reversed(entries):
					# print(entry)
					ext = '' if entry['Ext']=='' else '.'+entry['Ext']
					if not os.path.isdir(POSTS_DIR):
						os.mkdir(POSTS_DIR)
					# Prefix c4b_ to file name
					userFile_name = 'c4b_' + entry['Sid'] + ext
					userFile = os.path.join(POSTS_DIR, userFile_name)
					with open(userFile, 'w', encoding='utf-8') as fp:
						fp.write(entry['Body'])
					new_view = self.view.window().open_file(userFile)
			else:
				sublime.status_message("Queue is empty.")


# ------------------------------------------------------------------
# TA rewards brownies.
# ------------------------------------------------------------------
class c4baAwardPoint0Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 0)

class c4baAwardPoint1Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 1)

class c4baAwardPoint2Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 2)

class c4baAwardPoint3Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 3)

class c4baAwardPoint4Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 4)

class c4baAwardPoint5Command(sublime_plugin.TextCommand):
	def run(self, edit):
		award_points(self, edit, 5)

def award_points(self, edit, points):
	this_file_name = self.view.file_name()
	if this_file_name:
		info = c4ba_get_attr()
		if info is None:
			return
		basename = os.path.basename(this_file_name)
		if not basename.startswith('c4b_'):
			sublime.status_message("This is not a student submission.")
			return
		sid = basename.rsplit('.',-1)[0]
		sid = sid.split('c4b_')[-1]
		url = urllib.parse.urljoin(info['Server'], c4ba_BROWNIE_PATH)
		data = urllib.parse.urlencode({'sid':sid, 'points':points, 'passcode':info['Passcode']}).encode('utf-8')
		response = c4baRequest(url,data)
		if response == 'Failed':
			sublime.message_dialog("Failed to give brownies.")
		else:
			sublime.message_dialog(response)
			self.view.window().run_command('close')

# ------------------------------------------------------------------
class c4baTrackSubmissionsCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		info = c4ba_get_attr()
		if info is None:
			return
		webbrowser.open(info['Server'] + "/track_submissions")

# ------------------------------------------------------------------
class c4baSetServer(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4ba_FILE, 'r') as f:
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
				with open(c4ba_FILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			if not addr.startswith('http://'):
				addr = 'http://' + addr
			info['Server'] = addr
			with open(c4ba_FILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
		else:
			sublime.message_dialog("Server address is empty.")

# ------------------------------------------------------------------
class c4baSetPasscode(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4ba_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		if 'Passcode' not in info:
			info['Passcode'] = ''
		sublime.active_window().show_input_panel("Set passcode.  Press Enter:",
			info['Passcode'],
			self.set,
			None,
			None)

	def set(self, name):
		name = name.strip()
		if len(name) > 4:
			try:
				with open(c4ba_FILE, 'r') as f:
					info = json.loads(f.read())
			except:
				info = dict()
			info['Passcode'] = name
			with open(c4ba_FILE, 'w') as f:
				f.write(json.dumps(info, indent=4))
		else:
			sublime.message_dialog("Passcode is too short.")

# ------------------------------------------------------------------
class c4baAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BAssistant", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		sublime.message_dialog("Code4Brownies (v%s)\nCopyright © 2015-2018 Vinhthuy Phan" % version)

# ------------------------------------------------------------------
class c4baUpdate(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Are you sure you want to update Code4Brownies to the latest version?"):
			package_path = os.path.join(sublime.packages_path(), "C4BAssistant");
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			c4b_py = os.path.join(package_path, "Code4BrowniesAssistant.py")
			c4b_menu = os.path.join(package_path, "Main.sublime-menu")
			c4b_version = os.path.join(package_path, "VERSION")
			try:
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BAssistant/Code4BrowniesAssistant.py", c4b_py)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BAssistant/Main.sublime-menu", c4b_menu)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
				version = open(c4b_version).read().strip()
				sublime.message_dialog("Code4Brownies has been updated to version %s." % version)
			except:
				sublime.message_dialog("A problem occurred during update.")


