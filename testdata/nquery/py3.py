_ = None
N = int(input())
for i in range(0, N):
	if _ == None: _ = input().split()
	q = int(_.pop(0))
	if q==1:
		if _ == None: _ = input().split()
		w = string(_.pop(0))
	elif q==2:
		if _ == None: _ = input().split()
		l = int(_.pop(0))
		h = int(_.pop(0))
	else:
		pass
	_ = None
