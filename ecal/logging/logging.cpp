#include "logging.h"

#include <ecal/ecal_log.h>

void Log(eCAL_Logging_eLogLevel level, const char *const msg, size_t len) {
  eCAL::Logging::Log(level, std::string(msg, len));
}

void SetFileFilter(eCAL_Logging_Filter filter_bitset) {
  eCAL::Logging::SetFileLogFilter(filter_bitset);
}

void SetUDPFilter(eCAL_Logging_Filter filter_bitset) {
  eCAL::Logging::SetUDPLogFilter(filter_bitset);
}

void SetConsoleFilter(eCAL_Logging_Filter filter_bitset) {
  eCAL::Logging::SetConsoleLogFilter(filter_bitset);
}
