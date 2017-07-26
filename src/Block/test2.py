
def foo(L):
	# Begin Hint 1
	# What do you do when L is empty?
	if L == []:
		pass
	# End Hint 1

	# Begin Hint 2
	# What do you do when L is not empty?
	a = L[0]

	# Begin Hint 3
	# What do you do with the first item of L
	b = L[-1]
	# End Hint 3

	return 0
	# End Hint 2


if __name__ == '__main__':
	foo([1,2,3,4])