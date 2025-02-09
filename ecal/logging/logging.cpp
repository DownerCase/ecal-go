#include "logging.h"

#include "cgo_types.hpp"

#include <ecal/log.h>
#include <ecal/types/logging.h>

extern "C" {
extern void goCopyLogging(uintptr_t handle, CLogging *);
}

namespace {
CLogMessage toCLogMessage(const eCAL::Logging::SLogMessage &log) {
  return {
      log.time,
      log.host_name.c_str(),
      log.process_name.c_str(),
      log.unit_name.c_str(),
      log.content.c_str(),
      log.process_id,
      static_cast<CLogLevel>(log.level)
  };
}
} // namespace

void Log(CLogLevel level, const char *const msg, size_t len) {
  eCAL::Logging::Log(
      static_cast<eCAL::Logging::eLogLevel>(level),
      std::string(msg, len)
  );
}

void GetLogging(uintptr_t handle) {
  eCAL::Logging::SLogging logs{};
  eCAL::Logging::GetLogging(logs);
  const auto clogs = containerTo<std::vector>(logs.log_messages, toCLogMessage);
  CLogging clogging{};
  clogging.messages     = clogs.data();
  clogging.num_messages = clogs.size();
  goCopyLogging(handle, &clogging);
}

// Ensure exposed constants are correct
static_assert(
    log_level_none == static_cast<CLogLevel>(eCAL::Logging::log_level_none)
);
static_assert(
    log_level_info == static_cast<CLogLevel>(eCAL::Logging::log_level_info)
);
static_assert(
    log_level_warning ==
    static_cast<CLogLevel>(eCAL::Logging::log_level_warning)
);
static_assert(
    log_level_error == static_cast<CLogLevel>(eCAL::Logging::log_level_error)
);
static_assert(
    log_level_fatal == static_cast<CLogLevel>(eCAL::Logging::log_level_fatal)
);
static_assert(
    log_level_debug1 == static_cast<CLogLevel>(eCAL::Logging::log_level_debug1)
);
static_assert(
    log_level_debug2 == static_cast<CLogLevel>(eCAL::Logging::log_level_debug2)
);
static_assert(
    log_level_debug3 == static_cast<CLogLevel>(eCAL::Logging::log_level_debug3)
);
static_assert(
    log_level_debug4 == static_cast<CLogLevel>(eCAL::Logging::log_level_debug4)
);
static_assert(
    log_level_all == static_cast<CLogLevel>(eCAL::Logging::log_level_all)
);
