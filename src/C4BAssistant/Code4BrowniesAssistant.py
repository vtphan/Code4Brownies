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
c4ba_FEEDBACK_CODE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "feedback_code.txt")
c4ba_PEEK_PATH = "peek"
c4ba_REQUEST_ENTRY_PATH = "get_post_by_index"
c4ba_FEEDBACK_PATH = "feedback"
c4ba_BROWNIE_PATH = "give_points"
c4ba_REQUEST_ENTRIES_PATH = "get_posts"
c4ba_SHARE_WITH_TEACHER_PATH = "ta_share_with_teacher"
c4ba_GET_FROM_TEACHER_PATH = "ta_get_from_teacher"
c4ba_ADD_PUBLIC_BOARD_PATH = "add_public_board"
TIMEOUT = 7

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")

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
	if 'Name' not in json_obj:
		if verbose:
			sublime.message_dialog("Please set your name.")
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
class c4baAddPublicBoardCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4ba_get_attr()
		if info is None:
			return
		this_file_name = self.view.file_name()
		if this_file_name is None:
			sublime.message_dialog('Do not share an empty file.')
			return
		ext = this_file_name.rsplit('.',1)[-1]
		beg, end = self.view.sel()[0].begin(), self.view.sel()[0].end()
		content = self.view.substr(sublime.Region(beg,end))
		if len(content) <= 20:
			sublime.message_dialog('Select a larger region to share.')
			return
		data = urllib.parse.urlencode({
			'content': 		content,
			'ext': 			ext,
			'passcode':	info['Passcode'],
			'name':		info['Name'],
		}).encode('utf-8')
		url = urllib.parse.urljoin(info['Server'], c4ba_ADD_PUBLIC_BOARD_PATH)
		response = c4baRequest(url, data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4baViewPublicBoardCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		info = c4ba_get_attr()
		if info is None:
			return
		webbrowser.open(info['Server'] + "/view_public_board?i=0")

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
class c4baModifyFeedbackDef(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().open_file(c4ba_FEEDBACK_CODE)

# ------------------------------------------------------------------
class c4baShareWithTeacherCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4ba_get_attr()
		if info is None:
			return

		this_file_name = self.view.file_name()
		if this_file_name is None:
			sublime.message_dialog('Do not share an empty file.')
			return

		ext = this_file_name.rsplit('.',1)[-1]
		content = self.view.substr(sublime.Region(0, self.view.size()))
		data = urllib.parse.urlencode({
			'content': 		content,
			'ext': 			ext,
			'passcode':	info['Passcode'],
			'name':		info['Name'],
		}).encode('utf-8')
		url = urllib.parse.urljoin(info['Server'], c4ba_SHARE_WITH_TEACHER_PATH)
		response = c4baRequest(url, data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4baGetFromTeacherCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4ba_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4ba_GET_FROM_TEACHER_PATH)
		data = urllib.parse.urlencode({
			'passcode':	info['Passcode'],
			'name':		info['Name'],
		}).encode('utf-8')
		response = c4baRequest(url,data)
		if response is not None:
			entries = json.loads(response)
			if entries:
				for i in range(len(entries)-1, -1, -1):
					entry = entries[i]
					if not os.path.isdir(POSTS_DIR):
						os.mkdir(POSTS_DIR)
					outfile_name = 'teacher_{}.{}'.format(i,entry['Ext'])
					outfile = os.path.join(POSTS_DIR, outfile_name)
					with open(outfile, 'w', encoding='utf-8') as fp:
						fp.write(entry['Content'])
					new_view = self.view.window().open_file(outfile)
			else:
				sublime.message_dialog("Teacher has not shared anything yet.")

# ------------------------------------------------------------------
def c4ba_share_feedback(self, edit, points):
	# Get info
	info = c4ba_get_attr()
	if info is None:
		return

	# Detect empty buffer
	this_file_name = self.view.file_name()
	if this_file_name is None:
		sublime.message_dialog('Empty file is not a student submission.')
		return

	# Determine sid
	basename = os.path.basename(this_file_name)
	if not basename.startswith('c4b_'):
		sublime.message_dialog('This does not look like a student submission.')
		return

	# Get content and ext
	ext = this_file_name.rsplit('.',1)[-1]
	content = self.view.substr(sublime.Region(0, self.view.size()))
	header = content.split('\n',1)[0]
	if not header.startswith('#') or not header.startswith('//'):
		header = ''

	# Determine sid
	basename = os.path.basename(this_file_name)
	if not basename.startswith('c4b_'):
		sublime.message_dialog('This does not look like a student submission.')
		return
	sid = basename.split('.')[0]
	sid = sid.split('c4b_')[1]

	# Prepare and send feedback
	data = urllib.parse.urlencode({
		'content': 		content,
		'sid':			sid,
		'ext': 			ext,
		'points':		points,
		'passcode':		info['Passcode'],
		'name':			info['Name'],
		'has_feedback': _has_feedback(this_file_name, content)
	}).encode('utf-8')
	url = urllib.parse.urljoin(info['Server'], c4ba_FEEDBACK_PATH)
	response = c4baRequest(url, data)
	if response is not None:
		sublime.message_dialog(response)
		self.view.window().run_command('close')

# ------------------------------------------------------------------
# TA shares feedback on the current file
# ------------------------------------------------------------------
class c4baShareFeedbackUngradedCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_share_feedback(self, edit, 0)

class c4baShareFeedbackOneCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_share_feedback(self, edit, 1)

class c4baShareFeedbackTwoCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_share_feedback(self, edit, 2)

class c4baShareFeedbackThreeCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_share_feedback(self, edit, 3)

class c4baShareFeedbackFourCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_share_feedback(self, edit, 4)

class c4baShareFeedbackFiveCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_share_feedback(self, edit, 5)

# ------------------------------------------------------------------
def _build_feedback_line(file_name, fb):
	ext = file_name.rsplit('.', 1)[-1]
	if ext in ['py', 'pl', 'rb', 'txt', 'md']:
		prefix = '##> '
	else:
		prefix = '///> '
	return prefix + fb

# ------------------------------------------------------------------
def _has_feedback(file_name, content):
	ext = file_name.rsplit('.', 1)[-1]
	if ext in ['py', 'pl', 'rb', 'txt', 'md']:
		prefix = '##> '
	else:
		prefix = '///> '
	if prefix in content:
		return 1
	return 0

# ------------------------------------------------------------------
# TA inserts feedback on the current file
# ------------------------------------------------------------------
class c4baInsertFeedbackCommand(sublime_plugin.TextCommand):
	def read_feedback_def(self):
		if not os.path.exists(c4ba_FEEDBACK_CODE):
			feedback_code = [
				'GOOD JOB!!!',
				'Syntax',
				'Incorrect logic',
				'Base case',
				'Will not stop',
				'Return value',
				'Incorrect parameters',
				'Unreachable',
			]
			with open(c4ba_FEEDBACK_CODE, 'w') as f:
				f.write('\n'.join(feedback_code))

		with open(c4ba_FEEDBACK_CODE, 'r') as f:
			fb = f.readlines()
			fb = [ l.strip() for l in fb if l.strip() ]
		return fb

	#--------------------------------------------------------
	def run(self, edit):
		def on_done(i):
			if i < len(items):
				this_file_name = self.view.file_name()
				fb = _build_feedback_line(this_file_name, items[i])
				selection = self.view.sel()[0]
				line = self.view.substr(self.view.line(selection))
				if line != '':
					fb = '\t' + fb
				cursor = selection.begin()
				self.view.insert(edit, cursor, fb)

		items = self.read_feedback_def()
		self.view.show_popup_menu(items, on_done)

# ------------------------------------------------------------------
# how_many = -1 means gets all submissions.
# ------------------------------------------------------------------
def c4ba_get_submissions(self, edit, how_many):
	info = c4ba_get_attr()
	if info is None:
		return
	url = urllib.parse.urljoin(info['Server'], c4ba_REQUEST_ENTRIES_PATH)
	data = urllib.parse.urlencode({
		'how_many': 	how_many,
		'passcode':		info['Passcode'],
		'name':			info['Name'],
	}).encode('utf-8')
	response = c4baRequest(url,data)
	if response is not None:
		entries = json.loads(response)
		if entries:
			for entry in reversed(entries):
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
			sublime.message_dialog("There is no new submission.  Some might be pending.")

# ------------------------------------------------------------------
# Instructor retrieves submissions
# ------------------------------------------------------------------
class c4baGetOneSubmissionCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_get_submissions(self, edit, 1)

# ------------------------------------------------------------------
class c4baGetThreeSubmissionsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_get_submissions(self, edit, 3)

# ------------------------------------------------------------------
class c4baGetAllSubmissionsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_get_submissions(self, edit, -1)


# ------------------------------------------------------------------
# Preview submissions and select by index
# ------------------------------------------------------------------
class c4baPeekCommand(sublime_plugin.TextCommand):
	def request_entry(self, users, edit):
		def foo(selected):
			info = c4ba_get_attr()
			if info is None:
				return
			if selected < 0:
				return
			url = urllib.parse.urljoin(info['Server'], c4ba_REQUEST_ENTRY_PATH)
			data = urllib.parse.urlencode({
				'post':		selected,
				'passcode':		info['Passcode'],
				'name':			info['Name'],
			}).encode('utf-8')
			response = c4baRequest(url,data)
			if response is not None:
				json_obj = json.loads(response)
				# print(json_obj)
				ext = '' if json_obj['Ext']=='' else '.'+json_obj['Ext']
				if not os.path.isdir(POSTS_DIR):
					os.mkdir(POSTS_DIR)
				# Prefix c4b_ to file name
				userFile_name = 'c4b_' + json_obj['Sid'] + ext
				userFile = os.path.join(POSTS_DIR, userFile_name)
				with open(userFile, 'w', encoding='utf-8') as fp:
					fp.write(json_obj['Body'])
				new_view = self.view.window().open_file(userFile)
		return foo

	def run(self, edit):
		info = c4ba_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4ba_PEEK_PATH)
		data = urllib.parse.urlencode({
			'passcode':		info['Passcode'],
			'name':			info['Name'],
		}).encode('utf-8')
		response = c4baRequest(url,data)
		if response is not None:
			json_obj = json.loads(response)
			if json_obj is None:
				sublime.status_message("Queue is empty.")
			else:
				users = [ '%s: %s' % (entry['Uid'], entry['Sid']) for entry in json_obj ]
				if users:
					self.view.show_popup_menu(users, self.request_entry(users, edit))
				else:
					sublime.status_message("Queue is empty.")

# ------------------------------------------------------------------
class c4baTrackSubmissionsCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		info = c4ba_get_attr()
		if info is None:
			return
		webbrowser.open(info['Server'] + "/track_submissions?view=ta")

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
class c4baSetName(sublime_plugin.WindowCommand):
	def run(self):
		try:
			with open(c4ba_FILE, 'r') as f:
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
		try:
			with open(c4ba_FILE, 'r') as f:
				info = json.loads(f.read())
		except:
			info = dict()
		info['Name'] = name
		with open(c4ba_FILE, 'w') as f:
			f.write(json.dumps(info, indent=4))

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


