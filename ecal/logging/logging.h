#ifndef ECAL_GO_LOGGING
#define ECAL_GO_LOGGING

#include <ecal/ecal_log_level.h>

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

void Log(enum eCAL_Logging_eLogLevel level, const char *const msg, size_t len);

void SetFileFilter(eCAL_Logging_Filter filter_bitset);
void SetUDPFilter(eCAL_Logging_Filter filter_bitset);
void SetConsoleFilter(eCAL_Logging_Filter filter_bitset);

#ifdef __cplusplus
}
#endif

#endif
