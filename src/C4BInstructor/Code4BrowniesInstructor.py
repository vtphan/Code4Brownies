# Code4Brownies - Instructor module
# Author: Vinhthuy Phan, 2015-2017
#

import sublime, sublime_plugin
import urllib.parse
import urllib.request
import os
import json
import socket
# import ntpath
import webbrowser

SERVER_ADDR = "http://localhost:4030"
c4bi_BROADCAST_PATH = "broadcast"
# c4bi_BROWNIE_PATH = "give_points"
c4bi_PEEK_PATH = "peek"
c4bi_REQUEST_ENTRY_PATH = "get_post_by_index"
c4bi_REQUEST_ENTRIES_PATH = "get_posts"
c4bi_NEW_PROBLEM_PATH = "new_problem"
c4bi_START_POLL_PATH = "start_poll"
c4bi_ANSWER_POLL_PATH = "answer_poll"
c4bi_QUIZ_QUESTION_PATH = "send_quiz_question"
c4bi_FEEDBACK_PATH = "feedback"
c4bi_SHARE_WITH_TA_PATH = "share_with_ta"
c4bi_GET_FROM_TA_PATH = "get_from_ta"
c4bi_ADD_PUBLIC_BOARD_PATH = "add_public_board"

c4bi_FEEDBACK_CODE = os.path.join(os.path.dirname(os.path.realpath(__file__)), "feedback_code.txt")

TIMEOUT = 7

POSTS_DIR = os.path.join(os.path.dirname(os.path.realpath(__file__)), "Posts")

# ------------------------------------------------------------------
def c4biRequest(url, data, headers={}):
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
class c4biViewPublicBoardCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		webbrowser.open(SERVER_ADDR + "/view_public_board?i=0")

# ------------------------------------------------------------------
class c4biAddPublicBoardCommand(sublime_plugin.TextCommand):
	def run(self, edit):
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
		}).encode('utf-8')
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_ADD_PUBLIC_BOARD_PATH)
		response = c4biRequest(url, data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4biShareWithTaCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		this_file_name = self.view.file_name()
		if this_file_name is None:
			sublime.message_dialog('Do not share an empty file.')
			return

		ext = this_file_name.rsplit('.',1)[-1]
		content = self.view.substr(sublime.Region(0, self.view.size()))
		data = urllib.parse.urlencode({
			'content': 		content,
			'ext': 			ext,
		}).encode('utf-8')
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_SHARE_WITH_TA_PATH)
		response = c4biRequest(url, data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4biGetFromTa(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_GET_FROM_TA_PATH)
		data = urllib.parse.urlencode({}).encode('utf-8')
		response = c4biRequest(url,data)
		if response is not None:
			entries = json.loads(response)
			if entries:
				for i in range(len(entries)-1, -1, -1):
					entry = entries[i]
					if not os.path.isdir(POSTS_DIR):
						os.mkdir(POSTS_DIR)
					outfile_name = 'ta_{}.{}'.format(i,entry['Ext'])
					outfile = os.path.join(POSTS_DIR, outfile_name)
					with open(outfile, 'w', encoding='utf-8') as fp:
						fp.write(entry['Content'])
					new_view = self.view.window().open_file(outfile)
			else:
				sublime.message_dialog("TA has not shared anything yet.")

# ------------------------------------------------------------------
class c4biCleanCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Remove all student-submitted files."):
			if os.path.isdir(POSTS_DIR):
				files = [ f for f in os.listdir(POSTS_DIR) if f.startswith('c4b_') ]
				for f in files:
					local_file = os.path.join(POSTS_DIR, f)
					os.remove(local_file)
					sublime.status_message("remove " + local_file)

# ------------------------------------------------------------------
class c4biTrackSubmissionsCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		webbrowser.open(SERVER_ADDR + "/track_submissions?view=instructor")

# ------------------------------------------------------------------
class c4biViewQuestionsCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		webbrowser.open(SERVER_ADDR + "/view_questions")

# ------------------------------------------------------------------
class c4biClearWhiteboardsCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		url = urllib.parse.urljoin(SERVER_ADDR, "clear_whiteboards")
		data = urllib.parse.urlencode({}).encode('utf-8')
		response = c4biRequest(url,data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4biClearQuestionsCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		url = urllib.parse.urljoin(SERVER_ADDR, "clear_questions")
		data = urllib.parse.urlencode({}).encode('utf-8')
		response = c4biRequest(url,data)
		if response is not None:
			sublime.message_dialog(response)

# ------------------------------------------------------------------
class c4biStartPollCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		if sublime.ok_cancel_dialog('Poll the content of the current file.'):
			file_name = self.view.file_name()
			content = open(file_name, 'r', encoding='utf-8').read().strip()
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_START_POLL_PATH)
			data = urllib.parse.urlencode({'description': content}).encode('utf-8')
			response = c4biRequest(url,data)
			if response == 'Empty':
				sublime.message_dialog('Poll is empty. Please redo.')
			else:
				sublime.message_dialog('Poll started.')

# ------------------------------------------------------------------
class c4biViewPollCommand(sublime_plugin.ApplicationCommand):
	def run(self):
		webbrowser.open(SERVER_ADDR + "/view_poll")

# ------------------------------------------------------------------
class c4biAnswerPoll(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().show_input_panel("Answer",
			"",
			self.send_answer,
			None,
			None)

	def send_answer(self, answer):
		answer = answer.strip()
		if len(answer) > 0:
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_ANSWER_POLL_PATH)
			data = urllib.parse.urlencode({'answer': answer}).encode('utf-8')
			response = c4biRequest(url,data)
			if response is not None:
				sublime.message_dialog(response)
		else:
			sublime.message_dialog("Answer cannot be empty.")


# ------------------------------------------------------------------
def count_hints(str):
	s = str.strip()
	if s=='':
		return 0
	break_pattern = s.split('\n', 1)[0]
	hints = s.split(break_pattern)
	hints.pop(0)
	return len(hints)

# ------------------------------------------------------------------
class c4biQuizCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		file_name = self.view.file_name()
		if file_name is None:
			return
		ext = '' if file_name is None else file_name.rsplit('.',1)[-1]
		header = ''
		content = open(file_name, 'r', encoding='utf-8').read()
		if 'ANSWER:' not in content:
			sublime.message_dialog('''Each question of a quiz must have the following format:
<problem description>
ANSWER: <one_line_answer>
			''')
			return
		Q = content.split('ANSWER:')
		answers, questions = [], [Q[0].strip()]
		for i in range(1, len(Q)):
			items = Q[i].split('\n', 1)
			answers.append(items[0].strip())
			if i < len(Q)-1:
				questions.append(items[1].strip())

		if sublime.ok_cancel_dialog('The quiz appears to have {} questions.\nDo you want to hand out this quiz?'.format(len(questions))):
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_QUIZ_QUESTION_PATH)
			for i in range(len(questions)):
				data = urllib.parse.urlencode({
					'question': questions[i],
					'answer':	answers[i],
				}).encode('utf-8')
				response = c4biRequest(url,data)
				if response is not None:
					sublime.status_message(response)

# ------------------------------------------------------------------
class c4biTestCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		url = urllib.parse.urljoin(SERVER_ADDR, 'test')
		q = [
			dict(sids='XXX', content='first 1\nfirst 2'),
			dict(content='second 1\nsecond 2', something='YYY', sids='__all__', help='XXX'),
		]
		c = json.dumps(q).encode('utf-8')
		response = c4biRequest(url, c, headers={'content-type': 'application/json; charset=utf-8'})
		if response is not None:
			sublime.status_message(response)

# ------------------------------------------------------------------
def c4bi_share_feedback(self, edit, points=-1):
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
		'passcode':		'',
		'name':			'instructor',
		'has_feedback': _has_feedback(this_file_name, content)
	}).encode('utf-8')

	url = urllib.parse.urljoin(SERVER_ADDR, c4bi_FEEDBACK_PATH)
	response = c4biRequest(url, data)
	if response is not None:
		sublime.message_dialog(response)
		self.view.window().run_command('close')

# ------------------------------------------------------------------
# Instructor shares feedback on the current file
# ------------------------------------------------------------------
class c4biShareFeedbackUngradedCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_share_feedback(self, edit)

class c4biShareFeedbackZeroCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_share_feedback(self, edit, 0)

class c4biShareFeedbackOneCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_share_feedback(self, edit, 1)

class c4biShareFeedbackTwoCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_share_feedback(self, edit, 2)

class c4biShareFeedbackThreeCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_share_feedback(self, edit, 3)

class c4biShareFeedbackFourCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_share_feedback(self, edit, 4)

class c4biShareFeedbackFiveCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_share_feedback(self, edit, 5)

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
class c4biInsertFeedbackCommand(sublime_plugin.TextCommand):
	def read_feedback_def(self):
		if not os.path.exists(c4bi_FEEDBACK_CODE):
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
			with open(c4bi_FEEDBACK_CODE, 'w') as f:
				f.write('\n'.join(feedback_code))

		with open(c4bi_FEEDBACK_CODE, 'r') as f:
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
class c4biModifyFeedbackDef(sublime_plugin.WindowCommand):
	def run(self):
		sublime.active_window().open_file(c4bi_FEEDBACK_CODE)

# ------------------------------------------------------------------
# mode: 0 (unicast, current tab)
#		1 (multicast, all tabs)
#		2 (multicast, all tabs, randomized)
# ------------------------------------------------------------------
def _multicast(self, file_names, sids, mode):
	data = []
	for file_name in file_names:
		ext = file_name.rsplit('.',1)[-1]
		header = ''
		if file_name is not None:
			lines = open(file_name, 'r', encoding='utf-8').readlines()
			if len(lines)>0 and (lines[0].startswith('#') or lines[0].startswith('//')):
				header = lines[0]

		content = ''.join(lines)

		# Skip empty tabs
		if content.strip() == '':
			print('Skipping empty file:', file_name)
			break

		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_BROADCAST_PATH)
		if file_name is not None:
			basename = os.path.basename(file_name)
			dirname = os.path.dirname(file_name)
			help_content, test_content = '', ''
			if ext in ['py', 'go', 'java', 'c', 'pl', 'rb', 'txt', 'md']:
				prefix = basename.rsplit('.', 1)[0]
				help_file = os.path.join(dirname, prefix+'_hints.'+ext)
				if os.path.exists(help_file):
					help_content = open(help_file).read()
			# if basename.startswith('c4b_'):
			# 	original_sid = basename.split('.')[0]
			# 	original_sid = original_sid.split('c4b_')[1]
			# else:
			# 	original_sid = ''

		num_of_hints = count_hints(help_content)
		if num_of_hints > 0:
			sublime.message_dialog('There are {} hints associated with this exercise.'.format(num_of_hints))

		data.append({
			'content': 		content,
			'sids':			sids,
			'ext': 			ext,
			'help_content':	help_content,
			'hints':		num_of_hints,
			# 'original_sid':	original_sid,
			'mode': 		mode,
		})

	json_data = json.dumps(data).encode('utf-8')
	response = c4biRequest(url, json_data, headers={'content-type': 'application/json; charset=utf-8'})
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

# # ------------------------------------------------------------------
# Instructor broadcasts content on group defined by current window
# ------------------------------------------------------------------
class c4biMulticastRandCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		_broadcast(self, sids='__all__', mode=2)

# ------------------------------------------------------------------
# Instructor broadcasts content on group defined by current window
# ------------------------------------------------------------------
class c4biBroadcastGroupCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		fnames = [ v.file_name() for v in sublime.active_window().views() ]
		names = [ os.path.basename(n.rsplit('.',-1)[0]) for n in fnames if n is not None ]
		# Remove c4b_ prefix from file name
		sids = [ n.split('c4b_')[-1] for n in names if n.startswith('c4b_') ]
		if sids == []:
			sublime.message_dialog("No students' files in this window.")
			return
		if sublime.ok_cancel_dialog("Share this file with {} students whose submissions arein this window?".format(len(sids))):
			_broadcast(self, sids=','.join(sids), mode=0)

# ------------------------------------------------------------------
class c4biBroadcastCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		_broadcast(self, sids='__all__', mode=0)

# ------------------------------------------------------------------
# how_many = -1 means gets all submissions.
# ------------------------------------------------------------------
def c4bi_get_submissions(self, edit, how_many):
	url = urllib.parse.urljoin(SERVER_ADDR, c4bi_REQUEST_ENTRIES_PATH)
	data = urllib.parse.urlencode({
		'how_many': 	how_many,
		'name': 		'instructor',
	}).encode('utf-8')
	response = c4biRequest(url,data)
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
			sublime.status_message("There is no new submissions.  Some might be pending.")

# ------------------------------------------------------------------
# Instructor retrieves submissions
# ------------------------------------------------------------------
class c4biGetOneSubmissionCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_get_submissions(self, edit, 1)

# ------------------------------------------------------------------
class c4biGetThreeSubmissionsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_get_submissions(self, edit, 3)

# ------------------------------------------------------------------
class c4biGetAllSubmissionsCommand(sublime_plugin.TextCommand):
	def run(self, edit):
		c4bi_get_submissions(self, edit, -1)

# ------------------------------------------------------------------
# Preview submissions and select by index
# ------------------------------------------------------------------
class c4biPeekCommand(sublime_plugin.TextCommand):
	def request_entry(self, users, edit):
		def foo(selected):
			if selected < 0:
				return
			url = urllib.parse.urljoin(SERVER_ADDR, c4bi_REQUEST_ENTRY_PATH)
			data = urllib.parse.urlencode({
				'post':		selected,
				'name':		'instructor',
			}).encode('utf-8')
			response = c4biRequest(url,data)
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
		url = urllib.parse.urljoin(SERVER_ADDR, c4bi_PEEK_PATH)
		data = urllib.parse.urlencode({'name': 'instructor'}).encode('utf-8')
		response = c4biRequest(url,data)
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
class c4biAboutCommand(sublime_plugin.WindowCommand):
	def run(self):
		try:
			version = open(os.path.join(sublime.packages_path(), "C4BInstructor", "VERSION")).read().strip()
		except:
			version = 'Unknown'
		sublime.message_dialog("Code4Brownies (v%s)\nCopyright © 2015-2018 Vinhthuy Phan" % version)

# ------------------------------------------------------------------
class c4biUpdate(sublime_plugin.WindowCommand):
	def run(self):
		if sublime.ok_cancel_dialog("Are you sure you want to update Code4Brownies to the latest version?"):
			package_path = os.path.join(sublime.packages_path(), "C4BInstructor");
			if not os.path.isdir(package_path):
				os.mkdir(package_path)
			c4b_py = os.path.join(package_path, "Code4BrowniesInstructor.py")
			c4b_menu = os.path.join(package_path, "Main.sublime-menu")
			c4b_version = os.path.join(package_path, "VERSION")
			try:
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Code4BrowniesInstructor.py", c4b_py)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/C4BInstructor/Main.sublime-menu", c4b_menu)
				urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/Code4Brownies/master/src/VERSION", c4b_version)
				version = open(c4b_version).read().strip()
				sublime.message_dialog("Code4Brownies has been updated to version %s.  Latest server is at https://github.com/vtphan/Code4Brownies" % version)
			except:
				sublime.message_dialog("A problem occurred during update.")


