# Code4Brownies - TA module
# Author: Vinhthuy Phan, 2015-2018
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket

c4ba_FILE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "info")
c4ba_BROADCAST_PATH = "broadcast"
c4ba_BROWNIE_PATH = "give_points"
c4ba_REQUEST_ENTRIES_PATH = "get_posts"
TIMEOUT = 7

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")

# ------------------------------------------------------------------
def c4ba_get_attr():
	try:
		with open(c4ba_FILE, 'r') as f:
			json_obj = json.loads(f.read())
	except:
		sublime.message_dialog("Please set server address and your name.")
		return None
	if 'Server' not in json_obj or len(json_obj['Server']) < 4:
		sublime.message_dialog("Please set server address.")
		return None
	if 'Passcode' not in json_obj or len(json_obj['Passcode']) < 4:
		sublime.message_dialog("Please set passcode.")
		return None
	return json_obj

# ------------------------------------------------------------------
def c4baRequest(url, data, headers={}):
	req = urllib.request.Request(url, data, headers=headers)
	try:
		with urllib.request.urlopen(req, None, TIMEOUT) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.HTTPError as err:
		sublime.message_dialog("{0}".format(err))
	except urllib.error.URLError as err:
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
# mode: 0 (unicast, current tab)
#		1 (multicast, all tabs)
#		2 (multicast, all tabs, randomized)
# ------------------------------------------------------------------
def _multicast(self, file_names, sids, mode):
	info = c4ba_get_attr()
	if info is None:
		return

	data = []
	file_names = [ f for f in file_names if f is not None ]
	for file_name in file_names:
		ext = file_name.rsplit('.',1)[-1]
		header = ''
		lines = open(file_name, 'r', encoding='utf-8').readlines()
		if len(lines)>0 and (lines[0].startswith('#') or lines[0].startswith('//')):
			header = lines[0]

		content = ''.join(lines)

		# Skip empty tabs
		if content.strip() == '':
			print('Skipping empty file:', file_name)
			break

		basename = os.path.basename(file_name)
		dirname = os.path.dirname(file_name)
		if basename.startswith('c4b_'):
			original_sid = basename.split('.')[0]
			original_sid = original_sid.split('c4b_')[1]
		else:
			original_sid = ''
		data.append({
			'content': 		content,
			'sids':			sids,
			'ext': 			ext,
			'help_content':	'',
			'hints':		0,
			'original_sid':	original_sid,
			'mode': 		mode,
			'passcode':		info['Passcode'],
		})

	url = urllib.parse.urljoin(info['Server'], c4ba_BROADCAST_PATH)
	json_data = json.dumps(data).encode('utf-8')
	response = c4baRequest(url, json_data, headers={'content-type': 'application/json; charset=utf-8'})
	if response is not None:
		sublime.status_message(response)

# ------------------------------------------------------------------
def _broadcast(self, sids='__all__', mode=0):
	if mode == 0:
		fname = self.view.file_name()
		if fname is None:
			sublime.message_dialog('Cannot broadcast unsaved content.')
			return
		_multicast(self, [fname], sids, mode)
	else:
		fnames = [ v.file_name() for v in sublime.active_window().views() ]
		fnames = [ fname for fname in fnames if fname is not None ]
		if mode == 1:
			mesg = 'Broadcast all {} tabs in this window?'.format(len(fnames))
		elif mode == 2:
			mesg = 'Broadcast (randomized) all {} tabs in this window?'.format(len(fnames))
		else:
			sublime.message_dialog('Unable to broadcast. Unknown mode:', mode)
			return
		if sublime.ok_cancel_dialog(mesg):
			_multicast(self, fnames, sids, mode)

# ------------------------------------------------------------------
# TA gives feedback on this specific file
# ------------------------------------------------------------------
class c4baGiveFeedbackCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		this_file_name = self.view.file_name()
		if this_file_name is not None:
			sid = this_file_name.rsplit('.',-1)[0]
			sid = os.path.basename(sid)
			if sid.startswith('c4b_'):
				sid = sid.split('c4b_')[-1]
				_broadcast(self, sid)
			else:
				sublime.message_dialog("No student associated to this window.")

# ------------------------------------------------------------------
# TA broadcasts content on group defined by current window
# ------------------------------------------------------------------
# class c4baBroadcastGroupCommand(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		fnames = [ v.file_name() for v in sublime.active_window().views() ]
# 		names = [ os.path.basename(n.rsplit('.',-1)[0]) for n in fnames if n is not None ]
# 		# Remove c4b_ prefix from file name
# 		sids = [ n.split('c4b_')[-1] for n in names if n.startswith('c4b_') ]
# 		if sids == []:
# 			sublime.message_dialog("No students' files in this window.")
# 			return
# 		if sublime.ok_cancel_dialog("Share this file with {} students whose submissions arein this window?".format(len(sids))):
# 			_broadcast(self, ','.join(sids))

# ------------------------------------------------------------------
# class c4baBroadcastCommand(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		_broadcast(self)

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
			sublime.status_message("Failed to give brownies.")
		else:
			sublime.status_message(response)
			self.view.window().run_command('close')

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
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/Assistant/Code4BrowniesAssistant.py", c4b_py)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/Assistant/Main.sublime-menu", c4b_menu)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
				version = open(c4b_version).read().strip()
				sublime.message_dialog("Code4Brownies has been updated to version %s.  Latest server is at https://github.com/vtphan/Code4Brownies" % version)
			except:
				sublime.message_dialog("A problem occurred during update.")


