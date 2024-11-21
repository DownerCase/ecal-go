#ifndef ECAL_GO_REGISTRATION
#define ECAL_GO_REGISTRATION

#include <stddef.h>
#include <stdint.h>

#include "types.h"

#ifdef __cplusplus
extern "C" {
#endif

struct CQualityInfo {
  struct CDatatype datatype;
  uint8_t qualityFlags;
};

size_t AddPublisherEventCallback(uintptr_t handle);
void RemPublisherEventCallback(size_t handle);

#ifdef __cplusplus
}
#endif

#endif
