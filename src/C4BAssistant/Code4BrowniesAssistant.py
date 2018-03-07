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
c4ba_FEEDBACK_PATH = "ta_feedback"
c4ba_BROWNIE_PATH = "ta_give_points"
c4ba_REQUEST_ENTRIES_PATH = "ta_get_posts"
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
def c4ba_share_feedback(self, edit, points=-1):
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
	# lines = open(this_file_name, 'r', encoding='utf-8').readlines()
	# if len(lines)>0 and (lines[0].startswith('#') or lines[0].startswith('//')):
	# 	header = lines[0]
	# content = ''.join(lines)
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
		if points != -1:
			self.view.window().run_command('close')

# ------------------------------------------------------------------
# TA shares feedback on the current file
# ------------------------------------------------------------------
class c4baShareFeedbackUngradedCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4ba_share_feedback(self, edit)

class c4baShareFeedbackZeroCommand(sublime_plugin.TextCommand):
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


# class c4baInsertFeedbackCommand(sublime_plugin.TextCommand):
# 	#--------------------------------------------------------
# 	def on_cancel(self):
# 		self.view.hide_popup()

# 	#--------------------------------------------------------
# 	def process_feedback(self, line):
# 		self.view.hide_popup()

# 		# Retrieve feedback: <line no>,<feedback no> separated by a spaces
# 		items = line.strip().split()
# 		items = [ i.strip().split(',') for i in items ]
# 		feedback = [ (int(i[0]), self.feedback_def[int(i[1])]) for i in items ]

# 		# Insert feedback and build content
# 		this_file_name = self.view.file_name()
# 		content = ''
# 		with open(this_file_name, 'r', encoding='utf-8') as fp:
# 			lines = fp.readlines()
# 			for fb in feedback:
# 				line_no = fb[0]-1
# 				if line_no < len(lines):
# 					cur_line = lines[line_no].rstrip()
# 					if cur_line == '':
# 						cur_line = _build_feedback_line(this_file_name, fb[1])
# 					else:
# 						cur_line +=  '\t' + _build_feedback_line(this_file_name, fb[1])
# 					lines[line_no] = cur_line
# 				else:
# 					lines.append('\n' + _build_feedback_line(this_file_name, fb[1]))
# 			content = ''.join(lines)

# 		# Write back to file
# 		with open(this_file_name, 'w', encoding='utf-8') as fp:
# 			fp.write(content)

# 	#--------------------------------------------------------
# 	def read_feedback_def(self):
# 		if not os.path.exists(c4ba_FEEDBACK_CODE):
# 			feedback_code = [
# 				'GOOD JOB!!!',
# 				'Syntax',
# 				'Incorrect logic',
# 				'Base case',
# 				'Will not stop',
# 				'Return value',
# 				'Incorrect parameters',
# 				'Unreachable',
# 			]
# 			with open(c4ba_FEEDBACK_CODE, 'w') as f:
# 				f.write('\n'.join(feedback_code))

# 		self.feedback_def, instr = {}, []
# 		with open(c4ba_FEEDBACK_CODE, 'r') as f:
# 			i = 1
# 			for line in f:
# 				if line.strip():
# 					# print('{}. {}'.format(i, line.strip()))
# 					instr.append('{}. {}'.format(i,line))
# 					self.feedback_def[i] = line.strip()
# 					i += 1
# 		# sublime.status_message('Feedback code printed in console.')
# 		return '<br>'.join(instr)

# 	#--------------------------------------------------------
# 	def run(self, edit):
# 		# cursor = self.view.sel()[0].begin()
# 		# self.view.insert(edit, cursor, 'HELLO WORLD')

# 		sublime.active_window().show_input_panel("<line no>,<feedback no> (separate multiple items by a space)",
# 			"",
# 			self.process_feedback,
# 			None,
# 			self.on_cancel)

# 		# Show feedback code
# 		instr = self.read_feedback_def()
# 		if self.feedback_def:
# 			self.view.show_popup(instr)


# ------------------------------------------------------------------
# TA retrieves all posts.
# ------------------------------------------------------------------
class c4baGetAllCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		info = c4ba_get_attr()
		if info is None:
			return
		url = urllib.parse.urljoin(info['Server'], c4ba_REQUEST_ENTRIES_PATH)
		data = urllib.parse.urlencode({
			'passcode':	info['Passcode'],
			'name':		info['Name'],
		}).encode('utf-8')
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
				sublime.message_dialog("Queue is empty.")


# ------------------------------------------------------------------
# TA rewards brownies.
# ------------------------------------------------------------------
# class c4baAwardPoint0Command(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		award_points(self, edit, 0)

# class c4baAwardPoint1Command(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		award_points(self, edit, 1)

# class c4baAwardPoint2Command(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		award_points(self, edit, 2)

# class c4baAwardPoint3Command(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		award_points(self, edit, 3)

# class c4baAwardPoint4Command(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		award_points(self, edit, 4)

# class c4baAwardPoint5Command(sublime_plugin.TextCommand):
# 	def run(self, edit):
# 		award_points(self, edit, 5)

# def award_points(self, edit, points):
# 	this_file_name = self.view.file_name()
# 	if this_file_name:
# 		info = c4ba_get_attr()
# 		if info is None:
# 			return
# 		basename = os.path.basename(this_file_name)
# 		if not basename.startswith('c4b_'):
# 			sublime.status_message("This is not a student submission.")
# 			return
# 		sid = basename.rsplit('.',-1)[0]
# 		sid = sid.split('c4b_')[-1]
# 		url = urllib.parse.urljoin(info['Server'], c4ba_BROWNIE_PATH)
# 		data = urllib.parse.urlencode({
# 			'sid':		sid,
# 			'points':	points,
# 			'passcode':	info['Passcode'],
# 			'name':		info['Name'],
# 		}).encode('utf-8')
# 		response = c4baRequest(url,data)
# 		if response == 'Failed':
# 			sublime.message_dialog("Failed to give brownies.")
# 		else:
# 			sublime.message_dialog(response)
# 			self.view.window().run_command('close')

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


