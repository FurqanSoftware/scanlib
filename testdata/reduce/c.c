#include <stdio.h>

int main() {
	int N;
	scanf("%d", &N);
	int s;
	for (int i = 0; i < N; ++i) {
		int A;
		scanf("%d", &A);
		s = s+A;
	}
	
	return 0;
}
