_ = None
T = int(input())
for i in range(0, T):
	n, q = map(int, input().split())
	A = map(int, input().split())
	for j in range(0, q):
		if _ == None: _ = input().split()
		c = int(_.pop(0))
		if c==1:
			x = int(_.pop(0))
			y = int(_.pop(0))
		elif c==2:
			idx = int(_.pop(0))
		_ = None
