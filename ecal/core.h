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

const char *GetVersionString();
const char *GetVersionDateString();
struct version GetVersion();

int Initialize(
    struct config *config,
    const char *unit_name,
    unsigned int components
);
int Finalize();
bool IsInitialized(unsigned int component);
bool SetUnitName(const char *unit_name);
bool Ok();

#ifdef __cplusplus
}
#endif

#endif
