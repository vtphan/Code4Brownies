# Simulation of students submitting / sharing files
import argparse
import random
from urllib.parse import urljoin, urlencode
import urllib.request
import time

SERVER = 'http://localhost:4030'
SUBMIT_URL = urljoin(SERVER, 'submit_post')

def rand_text(line_count=1):
	alphabets = '     abcdefghijklmnopqrstuvwxyz0123456789'
	lines = []
	for i in range(line_count):
		line_width = random.randint(5,80)
		lines.append(''.join([alphabets[random.randint(0,len(alphabets)-1)] for j in range(line_width)]))
	return '\n'.join(lines)

def make_request(url, data):
	req = urllib.request.Request(url, data)
	try:
		with urllib.request.urlopen(req, None, 10) as response:
			return response.read().decode(encoding="utf-8")
	except urllib.error.URLError as err:
		sublime.message_dialog("{0}\nPossibly server not running or incorrect server address.".format(err))
		return None

parser = argparse.ArgumentParser(description='Simulate student submissions')
parser.add_argument('uid', type=str)
parser.add_argument('submit_intervals', type=int, nargs='+')
args = parser.parse_args()

print("Making", len(args.submit_intervals)+1, "requests to", SUBMIT_URL)
for i in range(len(args.submit_intervals)+1):
	values = {'uid':args.uid, 'body':'Submission #%s from %s.'%(i,args.uid), 'ext':'txt'}
	data = urllib.parse.urlencode(values).encode('ascii')
	if i==0:
		print("Request ", i)
	else:
		print("Request", i, "duration", args.submit_intervals[i-1])
		time.sleep(args.submit_intervals[i-1])
	response = make_request(SUBMIT_URL, data)
	print("")
