#ifndef ECAL_GO_CORE_H
#define ECAL_GO_CORE_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

struct version {
  int major;
  int minor;
  int patch;
};

struct config {
  char _padding;
};

const char *GetVersionString(void);
const char *GetVersionDateString(void);
struct version GetVersion(void);

int Initialize(
    struct config *config,
    const char *unit_name,
    unsigned int components
);
int Finalize(void);
bool IsInitialized(void);
bool IsComponentInitialized(unsigned int component);
bool SetUnitName(const char *unit_name);
bool Ok(void);

#ifdef __cplusplus
}
#endif

#endif
