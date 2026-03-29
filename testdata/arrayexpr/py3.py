_ = None
N = int(input())
A = [0] * N
s = 0
for i in range(0, N):
	if _ == None: _ = input().split()
	A[i] = int(_.pop(0))
	s = s+A[i]
	_ = None
