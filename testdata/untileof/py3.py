_ = None
while True:
	try:
		if _ == None: _ = input().split()
		cmd = string(_.pop(0))
		if cmd=="PUSH"||cmd=="REPEAT":
			param = int(_.pop(0))
		if cmd=="REPEAT":
			pass
		_ = None
	except EOFError as _:
		break
