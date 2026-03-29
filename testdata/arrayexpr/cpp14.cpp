#include <iostream>

using namespace std;

int main() {
	int N;
	cin >> N;
	int A[N];
	int s;
	for (int i = 0; i < N; ++i) {
		cin >> A[i];
		s = s+A[i];
	}
	
	return 0;
}
