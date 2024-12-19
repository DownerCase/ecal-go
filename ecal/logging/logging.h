#ifndef ECAL_GO_LOGGING
#define ECAL_GO_LOGGING

#include <ecal/ecal_log_level.h>

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct CLogging {
  const struct CLogMessage *messages;
  size_t num_messages;
};

void Log(enum eCAL_Logging_eLogLevel level, const char *const msg, size_t len);

void GetLogging(uintptr_t handle);

#ifdef __cplusplus
}
#endif

#endif
