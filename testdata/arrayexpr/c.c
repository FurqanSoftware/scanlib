#include <stdio.h>

int main() {
	int N;
	scanf("%d", &N);
	int A[N];
	int s;
	for (int i = 0; i < N; ++i) {
		scanf("%d", &A[i]);
		s = s+A[i];
	}
	
	return 0;
}
