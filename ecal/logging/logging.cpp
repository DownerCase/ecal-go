#include "logging.h"

#include "types.hpp"

#include <ecal/ecal_log.h>
#include <ecal/types/logging.h>

extern "C" {
extern void goCopyLogging(uintptr_t handle, CLogging *);
}

void Log(eCAL_Logging_eLogLevel level, const char *const msg, size_t len) {
  eCAL::Logging::Log(level, std::string(msg, len));
}

void GetLogging(uintptr_t handle) {
  eCAL::Logging::SLogging logs{};
  eCAL::Logging::GetLogging(logs);
  const auto clogs = toCTypes<CLogMessage>(logs.log_messages);
  CLogging clogging{};
  clogging.messages     = clogs.data();
  clogging.num_messages = clogs.size();
  goCopyLogging(handle, &clogging);
}
