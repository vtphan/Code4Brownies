import re

begin_tag = '# Begin Hint' + '\s+([1-9][0-9]*)$'
end_tag = '# End Hint' + '\s+([1-9][0-9]*)$'
empty_block_text_sub = '# ...'

'''
A block is a sequence of strings interleaved by other blocks.
'''
class Block(object):
	def __init__(self, lines, bid=0):
		self.bid = int(bid)
		self.lines = lines
		self.empty_sub = empty_block_text_sub
		self.parse()

	def parse(self):
		self.content = []
		while len(self.lines) > 0:
			new_block = self.new_block()
			if new_block is None:
				self.content.append( self.lines.pop(0) )
			else:
				self.content.append( new_block )

	# --------------------------------------------------------------------------
	# If self.lines[0] does not start a new block, return None.
	# Else, return the new block.
	# --------------------------------------------------------------------------
	def new_block(self):
		first_line = self.lines[0].strip()
		match = re.match(begin_tag, first_line)
		if match == None:
			return None
		block_id = match.groups()[0]
		first_line = self.lines.pop(0)
		new_lines = []
		while len(self.lines) > 0 and not self.block_ends(block_id):
			new_lines.append(self.lines.pop(0))
		items = re.split('(\s+)', first_line, 1)
		b = Block(new_lines,block_id)
		if len(items) == 1:
			b.empty_sub = empty_block_text_sub
		else:
			b.empty_sub = items[1] + empty_block_text_sub
		return b

	# --------------------------------------------------------------------------
	# If self.lines[0] ends the current block, pop that line and return True
	# Else, return False
	# --------------------------------------------------------------------------
	def block_ends(self, block_id):
		first_line = self.lines[0].strip()
		match = re.match(end_tag, first_line)
		if match == None:
			return False
		cur_block_id = match.groups()[0]
		if cur_block_id != block_id:
			return False
		self.lines.pop(0)
		return True

	# --------------------------------------------------------------------------
	def get_blocks(self, up_to):
		output = ''
		if up_to >= self.bid:
			for thing in self.content:
				if type(thing) == str:
					output += thing
				else:
					output += thing.get_blocks(up_to)
		else:
			output = self.empty_sub
		return output

	# --------------------------------------------------------------------------
	def __str__(self):
		return 'Block %s: %s, %s' % (self.bid, [ str(i) for i in self.content], self.empty_sub)

if __name__ == '__main__':
	lines = open('test2.py').readlines()
	b = Block(lines)
	print(b)
	print(">>Here's the block:")
	print(b.get_blocks(3))

