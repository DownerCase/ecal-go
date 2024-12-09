#ifndef ECAL_GO_CORE_H
#define ECAL_GO_CORE_H

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

struct version {
  int major;
  int minor;
  int patch;
};

struct CConfigLogging {
  bool receive_enabled;
};

struct CConfig {
  struct CConfigLogging logging;
};

const char *GetVersionString(void);
const char *GetVersionDateString(void);
struct version GetVersion(void);

int Initialize(
    struct CConfig *config,
    const char *unit_name,
    unsigned int components
);
int Finalize(void);
bool IsInitialized(void);
bool IsComponentInitialized(unsigned int component);
bool SetUnitName(const char *unit_name);
bool Ok(void);

void GetConfig(uintptr_t handle);

#ifdef __cplusplus
}
#endif

#endif
