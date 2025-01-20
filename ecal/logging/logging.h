#ifndef ECAL_GO_LOGGING
#define ECAL_GO_LOGGING

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

enum CLogLevel {
  log_level_none    = 0,
  log_level_info    = 1,
  log_level_warning = 2,
  log_level_error   = 4,
  log_level_fatal   = 8,
  log_level_debug1  = 16,
  log_level_debug2  = 32,
  log_level_debug3  = 64,
  log_level_debug4  = 128,
  log_level_all     = 255,
};

struct CLogMessage {
  int64_t time;
  const char *host_name;
  const char *process_name;
  const char *unit_name;
  const char *content;
  int32_t pid;
  enum CLogLevel level;
};

struct CLogging {
  const struct CLogMessage *messages;
  size_t num_messages;
};

void Log(enum CLogLevel level, const char *msg, size_t len);

void GetLogging(uintptr_t handle);

#ifdef __cplusplus
}
#endif

#endif
