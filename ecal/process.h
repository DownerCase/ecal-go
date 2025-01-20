#ifndef ECAL_GO_PROCESS_SEVERITY
#define ECAL_GO_PROCESS_SEVERITY

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

enum CProcessSeverity {
  process_severity_unknown  = 0,
  process_severity_healthy  = 1,
  process_severity_warning  = 2,
  process_severity_critical = 3,
  process_severity_failed   = 4,
};

#ifdef __cplusplus
}
#endif

#endif
