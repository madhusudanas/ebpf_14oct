package main

const bpfProgram = `
#include <uapi/linux/ptrace.h>

struct data_t {
  u64 arg1;
  u64 arg2;
  u64 arg3;
  u64 arg4;
  u64 arg5;
  u64 ret;
};

struct args_t {
  u64 arg1;
  u64 arg2;
  u64 arg3;
  u64 arg4;
  u64 arg5;
};

BPF_PERF_OUTPUT(trace);
BPF_HASH(margs, u32, struct args_t);

inline int start(struct pt_regs *ctx) {
  u32 pid = bpf_get_current_pid_tgid();
  struct args_t args = {};
  args.arg1 = ctx->ax;
  args.arg2 = ctx->bx;
  args.arg3 = ctx->cx;
  args.arg4 = ctx->di;
  args.arg5 = ctx->si;
  margs.insert(&pid, &args);
  return 0;
}

inline int end(struct pt_regs *ctx) {
  u32 pid = bpf_get_current_pid_tgid();
  struct args_t* args = margs.lookup(&pid);
  struct data_t data = {};
  if(args != NULL) {
    data.arg1 = args->arg1;
    data.arg2 = args->arg2;
    data.arg3 = args->arg3;
    data.arg4 = args->arg4;
    data.arg5 = args->arg5;
    data.ret = ctx->ax;
  }
  trace.perf_submit(ctx, &data, sizeof(data));
  margs.delete(&pid);
  return 0;
}
`
