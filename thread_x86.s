#include "textflag.h"

#define SYS_gettid	186

TEXT ·ThreadId(SB),NOSPLIT,$0-4
	MOVL	$SYS_gettid, AX
	SYSCALL
	MOVL	AX, ret+0(FP)
	RET

TEXT ·Gid(SB),NOSPLIT,$0-8
	MOVQ	(TLS), CX
	MOVQ	152(CX), AX
	MOVQ	AX, ret+0(FP)
	RET


