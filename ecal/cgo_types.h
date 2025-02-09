#ifndef ECAL_GO_TYPES
#define ECAL_GO_TYPES

#include <stddef.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct CDatatype {
  const char *name;
  const char *encoding;
  const void *descriptor;
  int descriptor_len;
  char _PADDING[4];
};

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
