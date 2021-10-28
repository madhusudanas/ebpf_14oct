package main

const bpfProgram = `
#include <uapi/linux/ptrace.h>

struct data_t {
  char arg1[256];
  char arg2[256];
  char arg3[256];
  char arg4[256];
  char arg5[256];
  char ret[256];
};

struct args_t {
  char arg1[256];
  char arg2[256];
  char arg3[256];
  char arg4[256];
  char arg5[256];
};

BPF_PERF_OUTPUT(trace);
BPF_HASH(margs, u32, struct args_t);

inline int start(struct pt_regs *ctx) {
  u32 pid = bpf_get_current_pid_tgid();
  struct args_t args = {};
  bpf_probe_read_user_str(args.arg1, sizeof(args.arg1), &ctx->ax);
  bpf_probe_read_user_str(args.arg2, sizeof(args.arg2), &ctx->bx);
  bpf_probe_read_user_str(args.arg3, sizeof(args.arg3), &ctx->cx);
  bpf_probe_read_user_str(args.arg4, sizeof(args.arg4), &ctx->di);
  bpf_probe_read_user_str(args.arg5, sizeof(args.arg5), &ctx->si);
  margs.insert(&pid, &args);
  return 0;
}

inline int end(struct pt_regs *ctx) {
  u32 pid = bpf_get_current_pid_tgid();
  struct args_t* args = margs.lookup(&pid);
  struct data_t data = {};
  if(args != NULL) {
    bpf_probe_read_user_str(data.arg1, sizeof(data.arg1), args->arg1);
    bpf_probe_read_user_str(data.arg2, sizeof(data.arg2), args->arg2);
    bpf_probe_read_user_str(data.arg3, sizeof(data.arg3), args->arg3);
    bpf_probe_read_user_str(data.arg4, sizeof(data.arg4), args->arg4);
    bpf_probe_read_user_str(data.arg5, sizeof(data.arg5), args->arg5);
    bpf_probe_read_user_str(data.ret, sizeof(data.ret), &ctx->ax);
  }
  trace.perf_submit(ctx, &data, sizeof(data));
  margs.delete(&pid);
  return 0;
}
`
